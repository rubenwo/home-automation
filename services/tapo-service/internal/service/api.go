package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"github.com/rubenwo/home-automation/services/tapo-service/internal/service/models"
	"github.com/rubenwo/home-automation/services/tapo-service/pkg/p100"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type api struct {
	router   *chi.Mux
	db       *pg.DB
	registry map[string]*p100.Client

	registryBaseUrl string
}

func Run(cfg *Config) error {

	db := pg.Connect(&pg.Options{
		Addr:     cfg.DatabaseBackend,
		User:     cfg.DatabaseUser,
		Password: cfg.DatabasePassword,
		Database: cfg.DatabaseName,
	})

	if err := db.Ping(context.Background()); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	if err := createSchema(db); err != nil {
		return fmt.Errorf("couldn't create schema: %w", err)
	}

	a := &api{
		router:          chi.NewRouter(),
		db:              db,
		registry:        make(map[string]*p100.Client),
		registryBaseUrl: cfg.RegistryBaseUrl,
	}

	var devices []models.DeviceConnectionInfo
	if err := a.db.Model(&devices).Select(); err != nil {
		return fmt.Errorf("couldn't load item from database: %s", err.Error())
	}

	for _, device := range devices {
		p, err := p100.New(device.IpAddress, device.Email, device.Password)
		if err != nil {
			return err
		}
		a.registry[device.DeviceId] = p
	}

	// A good base middleware stack
	a.router.Use(middleware.RequestID)
	a.router.Use(middleware.RealIP)
	a.router.Use(middleware.Logger)
	a.router.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	a.router.Use(middleware.Timeout(60 * time.Second))

	a.router.Get("/healthz", a.healthz)

	a.router.Get("/tapo/devices", a.getDevices)
	a.router.Get("/tapo/devices/{device_id}", a.getDevice)
	a.router.Delete("/tapo/devices/{device_id}", a.deleteDevice)
	a.router.Post("/tapo/devices/register", a.addDevice)
	a.router.Put("/tapo/devices/{device_id}", a.updateDevice)

	a.router.Get("/tapo/lights/{device_id}", a.commandDevice)

	if err := http.ListenAndServe(":80", a.router); err != nil {
		return err
	}
	return nil
}

func (a *api) healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(&models.HealthzModel{
		IsHealthy:    true,
		ErrorMessage: "",
	}); err != nil {
		log.Printf("error sending healthz: %s\n", err.Error())
	}
}

func (a *api) getDevices(w http.ResponseWriter, r *http.Request) {
	var devices []models.TapoDevice

	var connectionInfos []models.DeviceConnectionInfo
	if err := a.db.Model(&connectionInfos).Select(); err != nil {
		http.Error(w, fmt.Sprintf("couldn't retrieve models from database: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	for _, connectionInfo := range connectionInfos {
		if _, exists := a.registry[connectionInfo.DeviceId]; !exists {
			p, err := p100.New(connectionInfo.IpAddress, connectionInfo.Email, connectionInfo.Password)
			if err != nil {
				http.Error(w, fmt.Sprintf("couldn't connect to the tapo device"), http.StatusBadRequest)
				return
			}
			a.registry[connectionInfo.DeviceId] = p
		}
	}

	// fan out
	results := make(chan models.TapoDevice, len(connectionInfos))
	for _, connectionInfo := range connectionInfos {
		go func(ci models.DeviceConnectionInfo, result chan models.TapoDevice) {
			device := a.registry[ci.DeviceId]
			deviceInfo, err := device.DeviceInfo()
			if err != nil {
				log.Println(err)
				http.Error(w, fmt.Sprintf("couldn't retrieve tapo device info"), http.StatusBadRequest)
				return
			}
			results <- models.TapoDevice{
				DeviceId:   ci.DeviceId,
				DeviceName: device.Name(),
				DeviceType: ci.DeviceType,
				DeviceInfo: deviceInfo,
			}
		}(connectionInfo, results)
	}

	// fan in
	for i := 0; i < len(connectionInfos); i++ {
		devices = append(devices, <-results)
	}
	close(results) // We can close here since all routines should have returned by now

	var resp struct {
		Devices []models.TapoDevice `json:"devices"`
	}
	resp.Devices = devices

	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		log.Printf("error sending healthz: %s\n", err.Error())
	}
}

func (a *api) getDevice(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "device_id")

	var connectionInfo models.DeviceConnectionInfo
	if err := a.db.Model(&connectionInfo).Where("device_id = ?", id).Select(); err != nil {
		http.Error(w, fmt.Sprintf("couldn't retrieve model from database: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	p100Device, exists := a.registry[id]
	if !exists {
		p, err := p100.New(connectionInfo.IpAddress, connectionInfo.Email, connectionInfo.Password)
		if err != nil {
			http.Error(w, fmt.Sprintf("couldn't connect to the tapo device"), http.StatusBadRequest)
			return
		}
		a.registry[connectionInfo.DeviceId] = p
		p100Device = p
	}
	deviceInfo, err := p100Device.DeviceInfo()
	if err != nil {
		http.Error(w, fmt.Sprintf("couldn't retrieve tapo device info"), http.StatusBadRequest)
		return
	}
	var resp struct {
		Device models.TapoDevice `json:"device"`
	}
	resp.Device = models.TapoDevice{
		DeviceId:   connectionInfo.DeviceId,
		DeviceName: p100Device.Name(),
		DeviceType: connectionInfo.DeviceType,
		DeviceInfo: deviceInfo,
	}

	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		log.Printf("error sending healthz: %s\n", err.Error())
	}
}

func (a *api) deleteDevice(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "device_id")

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%s", a.registryBaseUrl, id), nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("error creating deletion request: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("error sending deletion request: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	if err := resp.Body.Close(); err != nil {
		log.Println(err)
	}

	result, err := a.db.Model(&models.DeviceConnectionInfo{DeviceId: id}).Where("device_connection_info.device_id = ?", id).Delete()
	if err != nil {
		http.Error(w, fmt.Sprintf("an error occured when deleteing item with id: %s, error: %s", id, err.Error()), http.StatusInternalServerError)
		return
	}
	fmt.Println(result)
	a.getDevices(w, r)
}

func (a *api) addDevice(w http.ResponseWriter, r *http.Request) {
	var device models.DeviceConnectionInfo
	if err := json.NewDecoder(r.Body).Decode(&device); err != nil {
		http.Error(w, fmt.Sprintf("couldn't decode body: %s", err.Error()), http.StatusBadRequest)
		return
	}

	p, err := p100.New(device.IpAddress, device.Email, device.Password)
	if err != nil {
		http.Error(w, fmt.Sprintf("couldn't connect to the tapo device"), http.StatusBadRequest)
		return
	}

	device.DeviceId = uuid.New().String()
	a.registry[device.DeviceId] = p

	result, err := a.db.Model(&device).Insert()
	if err != nil {
		http.Error(w, fmt.Sprintf("couldn't insert model into database: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	log.Println(result)

	var registryMsg struct {
		Id       string `json:"id"`
		Name     string `json:"name"`
		Category string `json:"category"`
		Product  struct {
			Company string `json:"company"`
			Type    string `json:"type"`
		} `json:"product"`
	}

	registryMsg.Id = device.DeviceId
	registryMsg.Name = p.Name()

	switch strings.ToLower(device.DeviceType) {
	case "p100":
		registryMsg.Category = "plug"
	case "l510e":
		registryMsg.Category = "light"
	default:
		registryMsg.Category = "unknown"
	}

	registryMsg.Product.Company = "tp-link"
	registryMsg.Product.Type = device.DeviceType

	data, err := json.Marshal(&registryMsg)
	if err != nil {
		http.Error(w, fmt.Sprintf("couldn't marshal registry message"), http.StatusInternalServerError)
		return
	}

	resp, err := http.Post(a.registryBaseUrl, "application/json", bytes.NewBuffer(data))
	if err != nil {
		http.Error(w, fmt.Sprintf("couldn't post new tapo device to the registry"), http.StatusInternalServerError)
		return
	}
	if err := resp.Body.Close(); err != nil {
		log.Println(err)
	}
	a.getDevices(w, r)
}

func (a *api) updateDevice(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "device_id")
	var connectionInfo models.DeviceConnectionInfo
	if err := json.NewDecoder(r.Body).Decode(&connectionInfo); err != nil {
		http.Error(w, fmt.Sprintf("couldn't decode body: %s", err.Error()), http.StatusBadRequest)
		return
	}
	if connectionInfo.DeviceId != id {
		http.Error(w, fmt.Sprintf("id in url does not correspond to the provided body"), http.StatusBadRequest)
		return
	}

	result, err := a.db.Model(&connectionInfo).Where("device_connection_info.device_id = ?", id).Update(&connectionInfo)
	if err != nil {
		http.Error(w, fmt.Sprintf("an error occured when deleting device info with id: %s, error: %s", id, err.Error()), http.StatusInternalServerError)
		return
	}
	log.Println(result)

	if _, exists := a.registry[id]; exists {
		p, err := p100.New(connectionInfo.IpAddress, connectionInfo.Email, connectionInfo.Password)
		if err != nil {
			http.Error(w, fmt.Sprintf("couldn't connect to the tapo device"), http.StatusInternalServerError)
			return
		}
		a.registry[connectionInfo.DeviceId] = p
	}

	w.WriteHeader(http.StatusOK)
}

func (a *api) commandDevice(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "device_id")
	query := r.URL.Query()
	command := query.Get("command")
	brightness, err := strconv.Atoi(query.Get("brightness"))
	if err != nil {
		http.Error(w, fmt.Sprintf("couldn't convert brightness to integer: %s", err.Error()), http.StatusBadRequest)
		return
	}
	var connectionInfo models.DeviceConnectionInfo
	if err := a.db.Model(&connectionInfo).Where("device_id = ?", id).Select(); err != nil {
		http.Error(w, fmt.Sprintf("couldn't retrieve model from database: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	device, exists := a.registry[id]
	if !exists {
		p, err := p100.New(connectionInfo.IpAddress, connectionInfo.Email, connectionInfo.Password)
		if err != nil {
			http.Error(w, fmt.Sprintf("couldn't connect to the tapo device"), http.StatusInternalServerError)
			return
		}
		a.registry[connectionInfo.DeviceId] = p
		device = p
	}

	if err := device.SetState(command == "on", brightness); err != nil {
		http.Error(w, fmt.Sprintf("couldn't command to the tapo device: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
