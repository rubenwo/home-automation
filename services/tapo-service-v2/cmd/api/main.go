package main

import (
	"github.com/rubenwo/home-automation/services/tapo-service/internal/service"
	"log"
)

func main() {
	cfg, err := service.LoadConfigFromPath("./config.json")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("tapo-service is online!")
	if err := service.Run(cfg); err != nil {
		log.Fatal(err)
	}
	log.Println("tapo-service is offline!")
}
