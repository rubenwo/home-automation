package api

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log"
	"net/http"
	"strconv"
	"time"
)

type api struct {
}

func New() *api {
	return &api{}
}

func (a *api) Run() error {
	router := chi.NewRouter()
	// A good base middleware stack
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	router.Use(middleware.Timeout(60 * time.Second))

	router.Get("/healthz", a.healthz)
	router.Get("/hue/information/{device_id}", a.deviceInfo)

	return http.ListenAndServe(":80", router)
}

func (a *api) healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	// TODO: when the api is not health/ready, send a 503 with a message
	if err := json.NewEncoder(w).Encode(&HealthzModel{
		IsHealthy:    true,
		ErrorMessage: "",
	}); err != nil {
		log.Printf("error sending healthz: %s\n", err.Error())
	}
}

func (a *api) deviceInfo(w http.ResponseWriter, r *http.Request) {
	deviceID, err := strconv.Atoi(chi.URLParam(r, "device_id"))
	if err != nil {
		if err := json.NewEncoder(w).Encode(&JsonError{
			StatusCode: http.StatusNotFound,
			Message:    err.Error(),
		}); err != nil {
			log.Printf("error sending deviceInfo jsonError: %s\n", err.Error())
		}
	}
	fmt.Printf("returning device information for device: %d\n", deviceID)
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(&DeviceInfo{
		DeviceID: deviceID,
		Status:   "On",
	}); err != nil {
		log.Printf("error sending deviceInfo: %s\n", err.Error())
	}
}
