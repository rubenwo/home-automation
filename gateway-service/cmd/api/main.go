package main

import (
	"github.com/rs/cors"
	"github.com/rubenwo/home-automation/gateway-service/pkg/ingress"
	"log"
	"net/http"
)

func main() {
	cfg, err := ingress.ParseConfig("./ingress.yaml")
	if err != nil {
		log.Fatal(err)
	}

	router, err := ingress.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	handler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}).Handler(router)

	if err := http.ListenAndServe(":80", handler); err != nil {
		log.Fatal(err)
	}
}
