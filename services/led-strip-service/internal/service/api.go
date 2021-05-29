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
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type api struct {
	router *chi.Mux
	db     *pg.DB
}

func New(cfg *Config) (*api, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config can't be nil")
	}
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("error validating config: %w", err)
	}

	db := pg.Connect(&pg.Options{
		Addr:     cfg.DatabaseBackend,
		User:     cfg.DatabaseUser,
		Password: cfg.DatabasePassword,
		Database: cfg.DatabaseName,
	})

	if err := db.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	if err := createSchema(db); err != nil {
		return nil, fmt.Errorf("couldn't create schema: %w", err)
	}

	a := &api{
		router: chi.NewRouter(),
		db:     db,
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
	a.router.Get("/leds/devices", a.getDevices)
	a.router.Post("/leds/devices/register", a.registerDevice)
	a.router.Post("/leds/devices/{id}/command", a.commandDevice)
	return a, nil
}
func (a *api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}
func (a *api) healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(&HealthzModel{
		IsHealthy:    true,
		ErrorMessage: "",
	}); err != nil {
		log.Printf("error sending healthz: %s\n", err.Error())
	}
}

func (a *api) getDevices(w http.ResponseWriter, r *http.Request) {
	var deviceInfos []LedDeviceModel
	if err := a.db.Model(&deviceInfos).Select(); err != nil {
		http.Error(w, fmt.Sprintf("couldn't load item from database: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	if deviceInfos == nil {
		deviceInfos = []LedDeviceModel{}
	}

	var resp struct {
		Devices []LedDeviceModel `json:"devices"`
	}
	resp.Devices = deviceInfos

	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		log.Printf("error sending devices: %s\n", err.Error())
	}
}

func (a *api) registerDevice(w http.ResponseWriter, r *http.Request) {
	var d RegisterLedDeviceModel
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	c := &http.Client{}
	resp, err := c.Get(fmt.Sprintf("http://%s/healthz", d.IPAddress))
	if err != nil {
		w.Header().Set("content-type", "application/json")
		if err := json.NewEncoder(w).Encode(&JsonError{
			Code:         http.StatusInternalServerError,
			ErrorMessage: err.Error(),
		}); err != nil {
			log.Printf("error sending json_error: %s\n", err.Error())
		}
		return
	}
	var ledHealthz HealthzModel
	if err := json.NewDecoder(resp.Body).Decode(&ledHealthz); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	fmt.Println(ledHealthz)
	resp.Body.Close()

	resp, err = c.Get(fmt.Sprintf("http://%s/info", d.IPAddress))
	if err != nil {
		w.Header().Set("content-type", "application/json")
		if err := json.NewEncoder(w).Encode(&JsonError{
			Code:         http.StatusInternalServerError,
			ErrorMessage: err.Error(),
		}); err != nil {
			log.Printf("error sending json_error: %s\n", err.Error())
		}
		return
	}
	var ledControllerInfo LedControllerInfo
	if err := json.NewDecoder(resp.Body).Decode(&ledControllerInfo); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	defer resp.Body.Close()
	id := uuid.New().String()

	device := LedDeviceModel{
		ID:             id,
		Name:           ledControllerInfo.DeviceName,
		NumLeds:        150,
		SupportedModes: ledControllerInfo.DeviceInfo.SupportedModes,
		CurrentMode:    ledControllerInfo.DeviceInfo.CurrentMode,
		IPAddress:      d.IPAddress,
	}

	result, err := a.db.Model(&device).Insert()
	if err != nil {
		http.Error(w, fmt.Sprintf("couldn't insert model into database: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	fmt.Println(result)
	var data struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Category string `json:"category"`
		Product  struct {
			Company string `json:"company"`
			Type    string `json:"type"`
		} `json:"product"`
	}
	data.ID = id
	data.Name = ledControllerInfo.DeviceName
	data.Category = "led-strip"
	data.Product.Company = "esp32"
	data.Product.Type = ledControllerInfo.DeviceType
	jsonData, err := json.Marshal(&data)
	if err != nil {
		if err := json.NewEncoder(w).Encode(&JsonError{
			Code:         http.StatusBadRequest,
			ErrorMessage: "id can't be empty",
		}); err != nil {
			log.Printf("error sending json_error: %s\n", err.Error())
		}
		return
	}

	req, err := http.NewRequest("POST", "http://registry.default.svc.cluster.local/devices", bytes.NewBuffer(jsonData))
	if err != nil {
		if err := json.NewEncoder(w).Encode(&JsonError{
			Code:         http.StatusBadRequest,
			ErrorMessage: "id can't be empty",
		}); err != nil {
			log.Printf("error sending json_error: %s\n", err.Error())
		}
		return
	}
	resp, err = c.Do(req)
	if err != nil {
		if err := json.NewEncoder(w).Encode(&JsonError{
			Code:         http.StatusBadRequest,
			ErrorMessage: "id can't be empty",
		}); err != nil {
			log.Printf("error sending json_error: %s\n", err.Error())
		}
		return
	}
	raw, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Println(string(raw), err)
	a.getDevices(w, r)
}

func (a *api) commandDevice(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		if err := json.NewEncoder(w).Encode(&JsonError{
			Code:         http.StatusBadRequest,
			ErrorMessage: "id can't be empty",
		}); err != nil {
			log.Printf("error sending json_error: %s\n", err.Error())
		}
		return
	}

	var device LedDeviceModel
	if err := a.db.Model(&device).Where("led_device_model.id = ?", id).Select(); err != nil {
		http.Error(w, fmt.Sprintf("couldn't load item from database: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	fmt.Println("sending command")
	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Header().Set("content-type", "application/json")
		if err := json.NewEncoder(w).Encode(&JsonError{
			Code:         http.StatusInternalServerError,
			ErrorMessage: err.Error(),
		}); err != nil {
			log.Printf("error sending json_error: %s\n", err.Error())
		}
		return
	}
	fmt.Println(string(raw))
	proxyReq, err := http.NewRequest("POST", fmt.Sprintf("http://%s/led", device.IPAddress), bytes.NewBuffer(raw))
	if err != nil {
		w.Header().Set("content-type", "application/json")
		if err := json.NewEncoder(w).Encode(&JsonError{
			Code:         http.StatusInternalServerError,
			ErrorMessage: err.Error(),
		}); err != nil {
			log.Printf("error sending json_error: %s\n", err.Error())
		}
		return
	}

	proxyReq.Header.Set("Host", r.Host)
	proxyReq.Header.Set("X-Forwarded-For", r.RemoteAddr)

	client := &http.Client{}
	proxyRes, err := client.Do(proxyReq)
	if err != nil {
		log.Println(err)
		w.Header().Set("content-type", "application/json")
		if err := json.NewEncoder(w).Encode(&JsonError{
			Code:         http.StatusInternalServerError,
			ErrorMessage: err.Error(),
		}); err != nil {
			log.Printf("error sending json_error: %s\n", err.Error())
		}
		return
	}
	var c CommandResponseModel
	if err := json.NewDecoder(proxyRes.Body).Decode(&c); err != nil {
		log.Println(err)
		w.Header().Set("content-type", "application/json")
		if err := json.NewEncoder(w).Encode(&JsonError{
			Code:         http.StatusInternalServerError,
			ErrorMessage: err.Error(),
		}); err != nil {
			log.Printf("error sending json_error: %s\n", err.Error())
		}
		return
	}
	fmt.Println(c)
	proxyRes.Body.Close()
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(&c); err != nil {
		log.Printf("error sending response: %v", err)
	}
}
