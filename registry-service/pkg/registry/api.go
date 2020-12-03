package registry

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log"
	"net/http"
	"time"
)

type api struct {
	router  *chi.Mux
	devices []DeviceInfo
}

func New(cfg *Config) (*api, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config can't be nil")
	}
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("error validating config: %w", err)
	}
	a := &api{router: chi.NewRouter()}
	// A good base middleware stack
	a.router.Use(middleware.RequestID)
	a.router.Use(middleware.RealIP)
	a.router.Use(middleware.Logger)
	a.router.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	a.router.Use(middleware.Timeout(60 * time.Second))

	a.router.Get("/devices", a.getDevices)
	a.router.Post("/devices", a.postDevice)

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

func (a *api) postDevice(w http.ResponseWriter, r *http.Request) {
	var dev DeviceInfo
	if err := json.NewDecoder(r.Body).Decode(&dev); err != nil {
		log.Println("error decoding the body:", err)
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	a.devices = append(a.devices, dev)
	a.getDevices(w, r)
}

func (a *api) getDevices(w http.ResponseWriter, r *http.Request) {
	var devices struct {
		Devices []DeviceInfo `json:"devices"`
	}
	devices.Devices = a.devices

	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(&devices); err != nil {
		log.Printf("error sending healthz: %s\n", err.Error())
	}
}
