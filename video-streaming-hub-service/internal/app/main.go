package app

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-pg/pg/v10"
	"github.com/gorilla/websocket"
	"github.com/rubenwo/home-automation/video-streaming-hub-service/pkg/rtsp"
	"log"
	"net/http"
	"strconv"
	"time"
)

type api struct {
	db          *pg.DB
	rtspClients []*rtsp.Client
}

func Run(cfg Config) error {
	if err := cfg.Validate(); err != nil {
	}

	db := pg.Connect(&pg.Options{
		Addr:     cfg.DatabaseAddr,
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
		db: db,
	}

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

	router.Get("/cameras", a.getCameras)
	router.Post("/cameras", a.addCamera)

	router.Delete("/camera/{id}", a.deleteCamera)
	router.Put("/camera/{id}", a.updateCamera)
	router.Get("/camera/{id}", a.getCamera)

	router.Get("/stream/{id}", a.streamVideoWS)

	if err := http.ListenAndServe(cfg.ApiAddr, router); err != nil {
		return fmt.Errorf("http.ListenAndServe error: %w", err)
	}

	return nil
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

func (a *api) addCamera(w http.ResponseWriter, r *http.Request) {
	var camera Camera
	if err := json.NewDecoder(r.Body).Decode(&camera); err != nil {
		http.Error(w, fmt.Sprintf("unable to decode camera: %s", err.Error()), http.StatusBadRequest)
		return
	}

	result, err := a.db.Model(&camera).Insert()
	if err != nil {
		http.Error(w, fmt.Sprintf("couldn't insert model into database: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	fmt.Println(result)
	// return all cameras now
	a.getCameras(w, r)
}

func (a *api) getCameras(w http.ResponseWriter, r *http.Request) {
	var cameras []Camera
	if err := a.db.Model(&cameras).Select(); err != nil {
		http.Error(w, fmt.Sprintf("couldn't load item from database: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	if cameras == nil {
		cameras = []Camera{}
	}

	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(&cameras); err != nil {
		log.Printf("error sending healthz: %s\n", err.Error())
	}
}

func (a *api) getCamera(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("couldn't convert id: %s", err.Error()), http.StatusBadRequest)
		return
	}

	var camera Camera
	if err := a.db.Model(&camera).Where("camera.Id = ?", id).Select(); err != nil {
		http.Error(w, fmt.Sprintf("couldn't load item from database: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(&camera); err != nil {
		log.Printf("error sending healthz: %s\n", err.Error())
	}
}

func (a *api) updateCamera(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("couldn't convert id: %s", err.Error()), http.StatusBadRequest)
		return
	}
	var camera Camera
	if err := json.NewDecoder(r.Body).Decode(&camera); err != nil {
		http.Error(w, fmt.Sprintf("couldn't decode body: %s", err.Error()), http.StatusBadRequest)
		return
	}
	if camera.Id != id {
		http.Error(w, "the id in the request path does not equal the id in the request body", http.StatusBadRequest)
		return
	}

	result, err := a.db.Model(&camera).Where("camera.id = ?", id).Update()
	if err != nil {
		http.Error(w, fmt.Sprintf("an error occured when updating item with id: %d, error: %s", id, err.Error()), http.StatusInternalServerError)
		return
	}
	fmt.Println(result)
	a.getCameras(w, r)
}

func (a *api) deleteCamera(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("couldn't convert id: %s", err.Error()), http.StatusBadRequest)
		return
	}

	result, err := a.db.Model(&Camera{Id: id}).Where("camera.id = ?", id).Delete()
	if err != nil {
		http.Error(w, fmt.Sprintf("an error occured when deleteing item with id: %d, error: %s", id, err.Error()), http.StatusInternalServerError)
		return
	}
	fmt.Println(result)
	a.getCameras(w, r)
}

func (a *api) streamVideoWS(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("couldn't convert id: %s", err.Error()), http.StatusBadRequest)
		return
	}

	var client *rtsp.Client
	found := false
	for _, c := range a.rtspClients {
		if c.Id == id {
			found = true
			client = c
			break
		}
	}
	if !found {
		var camera Camera
		if err := a.db.Model(&camera).Where("camera.id = ?", id).Select(); err != nil {
			http.Error(w, fmt.Sprintf("couldn't find camera with id: %d", id), http.StatusNotFound)
			return
		}
		client, err = rtsp.NewClient(rtsp.Config{Host: camera.Host, Id: camera.Id})
		if err != nil {
			http.Error(w, fmt.Sprintf("couldn't create rtsp client: %s", err.Error()), http.StatusNotFound)
			return
		}
	}
	if client == nil {
		http.Error(w, fmt.Sprintf("camera is nil"), http.StatusInternalServerError)
		return
	}

	// Upgrade connection
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	// Read messages from socket
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			conn.Close()
			return
		}
		log.Printf("msg: %s", string(msg))
	}
}
