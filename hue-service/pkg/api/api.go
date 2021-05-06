package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/amimof/huego"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"github.com/rubenwo/home-automation/hue-service/pkg/api/config"
	"github.com/rubenwo/home-automation/hue-service/pkg/hue"
	"log"
	"net/http"
	"strconv"
	"time"
)

type api struct {
	registryUrl string
	controller  *hue.Controller
}

func New(cfg *config.Config) (*api, error) {
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

	controller := hue.New(db)

	return &api{
		registryUrl: cfg.RegistryUrl,
		controller:  controller,
	}, nil
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

	router.Post("/hue/bridge/register", a.registerBridge)
	router.Get("/hue/bridge/discover", a.discoverBridge)

	router.Get("/hue/lights/{bridgeId}", a.lights)
	router.Put("/hue/lights/{bridgeId}/{lightId}", a.commandLight)

	return http.ListenAndServe(":80", router)
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

func (a *api) registerBridge(w http.ResponseWriter, r *http.Request) {}
func (a *api) discoverBridge(w http.ResponseWriter, r *http.Request) {
	bridge, _ := huego.Discover()
	user, _ := bridge.CreateUser(fmt.Sprintf("home-automation-hue-service-%s", uuid.New().String())) // Link button needs to be pressed
	bridge = bridge.Login(user)
	_ = a.controller.AddBridge(bridge)

	w.WriteHeader(http.StatusCreated)
}

func (a *api) lights(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "bridgeId")
	bridge, err := a.controller.Bridge(id)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("bridge with id: %s not found", id), http.StatusNotFound)
		return
	}
	lights, err := bridge.GetLights()
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("error retrieving lights from bridge: %s: %s", id, err.Error()), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&lights); err != nil {
		log.Println(err)
	}
}
func (a *api) commandLight(w http.ResponseWriter, r *http.Request) {
	bridgeId := chi.URLParam(r, "bridgeId")
	lightId, err := strconv.Atoi(chi.URLParam(r, "lightId"))
	if err != nil {
		log.Println(err)
		http.Error(w, "lightId needs to be off type int", http.StatusBadRequest)
		return
	}
	bridge, err := a.controller.Bridge(bridgeId)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("bridge with id: %s not found", bridgeId), http.StatusNotFound)
		return
	}

	var state huego.State
	if err := json.NewDecoder(r.Body).Decode(&state); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	light, err := bridge.GetLight(lightId)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("light with id: %d not found", lightId), http.StatusNotFound)
		return
	}

	if err := light.SetState(state); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
