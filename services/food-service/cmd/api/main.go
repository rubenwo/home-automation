package main

import (
	"github.com/rubenwo/home-automation/services/food-service/pkg/cmd/foodservice"
	"log"
)

func main() {

	log.Println("food-service is now running!")
	if err := foodservice.Run(&foodservice.Config{
		DatabaseAddr: "postgres.default.svc.cluster.local:5432",
		//DatabaseAddr:     "localhost:5432",
		DatabaseUser:     "user",
		DatabasePassword: "password",
		DatabaseName:     "food_database",
		ApiAddr:          ":80",
	}); err != nil {
		log.Fatal(err)
	}
}
