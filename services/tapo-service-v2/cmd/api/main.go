package main

import (
	"github.com/rubenwo/home-automation/services/tapo-service/internal/service"
	"log"
)

func main() {

	log.Println("tapo-service is online!")
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
	log.Println("tapo-service is offline!")
}
