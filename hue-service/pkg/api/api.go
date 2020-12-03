package api

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log"
	"net/http"
	"time"
)

type api struct {
	registryUrl string
}

func New() *api {
	return &api{
		registryUrl: "http://registry.default.svc.cluster.local/devices",
	}
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
	router.Get("/hue/devices", a.getHueDevices)

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

func (a *api) getHueDevices(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusFound)
}
