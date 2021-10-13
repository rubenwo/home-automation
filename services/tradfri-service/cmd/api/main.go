package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/rubenwo/home-automation/services/tradfri-service/internal/app"
	"github.com/rubenwo/home-automation/services/tradfri-service/internal/dao"
	"github.com/rubenwo/home-automation/services/tradfri-service/internal/interfaces/api"
	"github.com/rubenwo/home-automation/services/tradfri-service/internal/services/registrysyncer"
	"github.com/rubenwo/home-automation/services/tradfri-service/internal/usecases"
	"github.com/rubenwo/home-automation/services/tradfri-service/pkg/database"
	"log"
	"net/http"
	"time"
)

func migrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://database/migrations",
		"postgres", driver)
	fmt.Println(m.Version())
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			return nil
		}
		return err
	}
	return nil
}

func fixtures(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://database/fixtures",
		"postgres", driver)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			return nil
		}
		return err
	}
	return nil
}

func main() {
	db, err := database.NewPostgresDB(database.Config{
		Host:         "192.168.178.46:5432",
		User:         "user",
		Password:     "password",
		Database:     "tradfri_database",
		Options:      []string{"sslmode=disable"},
		MaxOpenConns: 100,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := migrations(db); err != nil {
		log.Fatalf("migrations: %v", err)
	}
	if err := fixtures(db); err != nil {
		log.Fatalf("fixtures: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var dataAccessOject = &app.DataAccessObject{
		TradfriDao: dao.NewTradfriDB(db),
	}

	registrySyncerService := registrysyncer.NewService(db)
	registrySyncerService.Run(ctx)

	services := &app.Services{
		RegistrySyncerService: registrySyncerService,
	}

	var (
		usecasesTradfri = usecases.NewTradfriUsecases(dataAccessOject, services)
	)

	go func() {
		devices, err := usecasesTradfri.FetchAllDevices()
		if err != nil {
			log.Fatalf("error fetching tradfri devices: %s\n", err.Error())
		}

		for _, device := range devices {
			if err := services.RegistrySyncerService.PublishDevice(device); err != nil {
				log.Fatalf("error publishing tradfri device to registry: %s\n", err.Error())
			}
		}
		const retry = 10

		timer := time.NewTimer(time.Minute)
		errorCount := 0
		for {
			select {
			case <-timer.C:
				newDevices, err := usecasesTradfri.FetchAllDevices()
				if err != nil {
					log.Printf("error publishing tradfri device to registry: %s\n", err.Error())
					if errorCount > retry {
						log.Fatalf("max retries for fetching devices reached: %s\n", err.Error())
					}
					errorCount++
					timer.Reset(time.Minute)
					continue
				}

				for _, newDevice := range newDevices {
					add := true
					for _, device := range devices {
						if device == newDevice {
							add = false
							break
						}
					}
					if add {
						if err := services.RegistrySyncerService.PublishDevice(newDevice); err != nil {
							log.Printf("error publishing tradfri device to registry: %s\n", err.Error())
						}
					}
				}

				errorCount = 0
				timer.Reset(time.Minute)
			}
		}

	}()

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

	api.RegisterHandler(usecasesTradfri, router)

	log.Println("tradfri-service is online")
	if err := http.ListenAndServe(":80", router); err != nil {
		log.Fatal(err)
	}
	log.Println("tradfri-service is offline")
}
