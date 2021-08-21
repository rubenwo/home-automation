package main

import (
	"github.com/rubenwo/home-automation/services/led-strip-service/internal/service"
	"log"
	"net/http"
)

func main() {
	cfg, err := service.LoadConfigFromPath("./config.json")
	if err != nil {
		log.Fatal(err)
	}
	router, err := service.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("led-strip-service is online!")
	if err := http.ListenAndServe(":80", router); err != nil {
		log.Fatal(err)
	}
	log.Println("led-strip-service is offline!")
}
