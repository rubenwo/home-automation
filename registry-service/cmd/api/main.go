package main

import (
	"github.com/rubenwo/home-automation/registry-service/pkg/registry"
	"log"
	"net/http"
)

func main() {
	router, err := registry.New(&registry.Config{
		DatabaseBackend:   "redis",
		PgDatabaseBackend: "postgres.default.svc.cluster.local:5432",
		//PgDatabaseBackend:     "192.168.2.135:5432",
		PgDatabaseUser:     "user",
		PgDatabasePassword: "password",
		PgDatabaseName:     "registry_database",
	})
	if err != nil {
		log.Fatal(err)
	}
	if err := http.ListenAndServe(":80", router); err != nil {
		log.Fatal(err)
	}
}
