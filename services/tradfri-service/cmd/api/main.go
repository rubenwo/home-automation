package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rubenwo/home-automation/services/tradfri-service/internal/interfaces/api"
	"log"
	"net/http"
	"time"
)

func main() {

	//var dao struct{}
	//var services struct{}
	//var usecases struct{}

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

	api.RegisterHandler(router)

	log.Println("tradfri-service is online")
	if err := http.ListenAndServe(":80", router); err != nil {
		log.Fatal(err)
	}
	log.Println("tradfri-service is offline")
}
