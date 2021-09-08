package service

import (
	"context"
	"encoding/json"
	"firebase.google.com/go/v4/messaging"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-pg/pg/v10"
	"github.com/rubenwo/home-automation/servics/notification-service/pkg/notification"
	"net/http"
)

type api struct {
	router *chi.Mux
	db     *pg.DB

	notificationManager *notification.Manager
}

func New(cfg *Config, msgClient *messaging.Client) (*api, error) {
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

	notificationManager := notification.NewManager(cfg.MqttHost, 10, msgClient, db)

	a := &api{
		router:              chi.NewRouter(),
		db:                  db,
		notificationManager: notificationManager,
	}
	a.notificationManager.SendNotification(notification.Notification{
		Title: "Notification service online",
		Body:  "The notification service is now online",
	})

	a.router.Post("/notifications/subscribe", a.subscribeToNotifications)

	return a, nil
}
func (a *api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}

func (a *api) subscribeToNotifications(w http.ResponseWriter, r *http.Request) {
	var msg NotificationSubscriber
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&JsonError{
			Code:         http.StatusBadRequest,
			ErrorMessage: err.Error(),
		})
	}

	result, err := a.db.Model(&msg).Insert()
	if err != nil {
		http.Error(w, fmt.Sprintf("couldn't insert model into database: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	fmt.Println(result)
}
