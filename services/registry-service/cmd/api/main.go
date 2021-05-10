package main

import (
	"github.com/rubenwo/home-automation/services/registry-service/pkg/registry"
	"log"
	"net/http"
)

func main() {
	cfg, err := registry.LoadConfigFromPath("./config.json")
	if err != nil {
		log.Fatal(err)
	}
	router, err := registry.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("registry-service is online!")
	if err := http.ListenAndServe(":80", router); err != nil {
		log.Fatal(err)
	}
}
