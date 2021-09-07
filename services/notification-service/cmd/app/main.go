package main

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"fmt"
	"github.com/rubenwo/home-automation/servics/notification-service/internal/service"
	"google.golang.org/api/option"
	"log"
	"net/http"
)

func main() {
	cfg, err := service.LoadConfigFromPath("./config.json")
	if err != nil {
		log.Fatal(err)
	}

	opt := option.WithCredentialsFile("./service-account-file.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
	msgClient, err := app.Messaging(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	msg, err := msgClient.Send(context.Background(), &messaging.Message{
		Token: "blah",
		Notification: &messaging.Notification{
			Title: "Hello World",
			Body:  "Bla",
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(msg)

	router, err := service.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("notification-service is online!")
	if err := http.ListenAndServe(":80", router); err != nil {
		log.Fatal(err)
	}
	log.Println("notification-service is offline!")
}
