package main

import (
	"github.com/rubenwo/home-automation/hue-service/pkg/api"
	"github.com/rubenwo/home-automation/hue-service/pkg/api/config"
	"log"
)

func main() {
	apiConfig, err := config.LoadConfigFromPath("./config.json")
	if err != nil {
		log.Fatal(err)
	}
	a, err := api.New(apiConfig)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("hue-service is up and running!")
	if err := a.Run(); err != nil {
		log.Fatal(err)
	}
}
