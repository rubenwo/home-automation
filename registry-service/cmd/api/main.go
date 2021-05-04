package main

import (
	"github.com/rubenwo/home-automation/registry-service/pkg/registry"
	"log"
	"net/http"
)

func main() {
	router, err := registry.New(&registry.Config{
		DatabaseBackend: "postgres.default.svc.cluster.local:5432",
		//DatabaseBackend:  "192.168.2.135:5432",
		//DatabaseBackend:  "localhost:5432",
		DatabaseUser:     "user",
		DatabasePassword: "password",
		DatabaseName:     "registry_database",
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("registry-service is online!")
	if err := http.ListenAndServe(":80", router); err != nil {
		log.Fatal(err)
	}
}
