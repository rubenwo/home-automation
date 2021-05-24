package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-pg/pg/v10"
	"github.com/rubenwo/home-automation/services/tapo-service/internal/service/models"
	"github.com/rubenwo/home-automation/services/tapo-service/pkg/p100"
	"log"
	"net/http"
	"time"
)

type api struct {
	router   *chi.Mux
	db       *pg.DB
	registry map[string]*p100.Client
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
		router:   chi.NewRouter(),
		db:       db,
		registry: make(map[string]*p100.Client),
	}

	var devices []models.DatabaseDevice
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
	var devices []models.Device

	for id, client := range a.registry {
		device := models.Device{
			DeviceId:   id,
			DeviceName: client.Name(),
			DeviceType: "",
			DeviceInfo: nil,
		}
		devices = append(devices, device)
	}

	var resp struct {
		Devices []models.Device `json:"devices"`
	}
	resp.Devices = devices

	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		log.Printf("error sending healthz: %s\n", err.Error())
	}
}

func (a *api) getDevice(w http.ResponseWriter, r *http.Request) {

}

func (a *api) deleteDevice(w http.ResponseWriter, r *http.Request) {

}

func (a *api) addDevice(w http.ResponseWriter, r *http.Request) {

}

func (a *api) updateDevice(w http.ResponseWriter, r *http.Request) {

}

func (a *api) commandDevice(w http.ResponseWriter, r *http.Request) {

}
