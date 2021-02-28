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
	"strconv"
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

	a.router.Get("/inventory", a.getInventory)
	a.router.Post("/inventory", a.addItemToInventory)
	a.router.Delete("/inventory/{id}", a.deleteItemFromInventory)
	a.router.Put("/inventory/{id}", a.updateItemFromInventory)

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

func (a *api) getInventory(w http.ResponseWriter, r *http.Request) {
	var items []Item
	if err := a.db.Model(&items).Select(); err != nil {
		http.Error(w, fmt.Sprintf("couldn't load item from database: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	if items == nil {
		items = []Item{}
	}

	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(&items); err != nil {
		log.Printf("error sending healthz: %s\n", err.Error())
	}
}

func (a *api) addItemToInventory(w http.ResponseWriter, r *http.Request) {
	var item Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, fmt.Sprintf("couldn't decode body: %s", err.Error()), http.StatusBadRequest)
		return
	}

	result, err := a.db.Model(&item).Insert()
	if err != nil {
		http.Error(w, fmt.Sprintf("couldn't insert model into database: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	fmt.Println(result)
	// return the entire inventory now
	a.getInventory(w, r)
}

func (a *api) deleteItemFromInventory(w http.ResponseWriter, r *http.Request) {
	rawId := chi.URLParam(r, "id")
	if rawId == "" {
		http.Error(w, "no id was provided in the request", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(rawId)
	if err != nil {
		http.Error(w, "provided id was not a number, thus couldn't be parsed", http.StatusBadRequest)
		return
	}

	result, err := a.db.Model(&Item{Id: int64(id)}).Where("item.id = ?", id).Delete()
	if err != nil {
		http.Error(w, fmt.Sprintf("an error occured when deleteing item with id: %d, error: %s", id, err.Error()), http.StatusInternalServerError)
		return
	}
	fmt.Println(result)
	a.getInventory(w, r)
}

func (a *api) updateItemFromInventory(w http.ResponseWriter, r *http.Request) {
	rawId := chi.URLParam(r, "id")
	if rawId == "" {
		http.Error(w, "no id was provided in the request", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(rawId)
	if err != nil {
		http.Error(w, "provided id was not a number, thus couldn't be parsed", http.StatusBadRequest)
		return
	}
	var item Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, fmt.Sprintf("couldn't decode body: %s", err.Error()), http.StatusBadRequest)
		return
	}
	if item.Id != int64(id) {
		http.Error(w, "the id in the request path does not equal the id in the request body", http.StatusBadRequest)
		return
	}

	result, err := a.db.Model(&item).Where("item.id = ?", id).Update()
	if err != nil {
		http.Error(w, fmt.Sprintf("an error occured when updating item with id: %d, error: %s", id, err.Error()), http.StatusInternalServerError)
		return
	}
	fmt.Println(result)
	a.getInventory(w, r)
}
