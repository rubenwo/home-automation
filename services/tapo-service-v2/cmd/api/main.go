package main

import (
	"fmt"
	"github.com/rubenwo/home-automation/services/tapo-service/internal/service"
	"github.com/rubenwo/home-automation/services/tapo-service/pkg/p100"
	"log"
)

func main() {
	p, err := p100.New("", "", "")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(p.DeviceInfo())
	if err := p.SetState(false, 0); err != nil {
		log.Fatal(err)
	}
	if err := p.SetState(true, 100); err != nil {
		log.Fatal(err)
	}

	log.Println("tapo-service is online!")
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
	log.Println("tapo-service is offline!")
}
