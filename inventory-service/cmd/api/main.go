package main

import (
	"fmt"
	"github.com/rubenwo/home-automation/inventory-service/pkg/inventory"
	"log"
	"net/http"
)

func main() {
	// TODO: get the database credentials from environment
	router, err := inventory.New(&inventory.Config{
		DatabaseAddr:     "postgres.default.svc.cluster.local:5432",
		//DatabaseAddr:     "192.168.2.135:5432",
		DatabaseUser:     "user",
		DatabasePassword: "password",
		DatabaseName:     "inventory_database",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("inventory service is now online")
	if err := http.ListenAndServe(":80", router); err != nil {
		log.Fatal(err)
	}
}
