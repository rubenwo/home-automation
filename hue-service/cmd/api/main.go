package main

import (
	"github.com/rubenwo/home-automation/hue-service/pkg/api"
	"log"
)

func main() {
	a := api.New()
	if err := a.Run(); err != nil {
		log.Fatal(err)
	}
}
