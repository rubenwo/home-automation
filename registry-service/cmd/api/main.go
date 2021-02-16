package main

import (
	"github.com/rubenwo/home-automation/registry-service/pkg/registry"
	"log"
	"net/http"
)

func main() {
	router, err := registry.New(&registry.Config{DatabaseBackend: "redis"})
	if err != nil {
		log.Fatal(err)
	}
	if err := http.ListenAndServe(":80", router); err != nil {
		log.Fatal(err)
	}
}
