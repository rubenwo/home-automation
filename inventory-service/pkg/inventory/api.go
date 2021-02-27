package inventory

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-pg/pg/v10"
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
		return nil, fmt.Errorf("cfg can't be nil")
	}
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	db := pg.Connect(&pg.Options{
		Addr:     cfg.DatabaseAddr,
		User:     cfg.DatabaseUser,
		Password: cfg.DatabasePassword,
		Database: cfg.DatabaseName,
	})

	if err := db.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
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
