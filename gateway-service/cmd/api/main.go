package main

import (
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

	if err := http.ListenAndServe(":80", router); err != nil {
		log.Fatal(err)
	}
}
