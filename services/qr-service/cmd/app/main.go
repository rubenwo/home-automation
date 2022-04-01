package main

import (
	"github.com/rubenwo/home-automation/services/qr-service/internal/api"
	"github.com/rubenwo/home-automation/services/qr-service/internal/api/config"
	"log"
)

func main() {
	apiConfig, err := config.LoadConfigFromPath("./config.json")
	if err != nil {
		log.Fatal(err)
	}
	a, err := api.New(*apiConfig)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("qr-service is up and running on %s\n", apiConfig.Listen)
	if err := a.Run(); err != nil {
		log.Fatal(err)
	}
}
