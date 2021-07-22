package main

import (
	"github.com/rubenwo/home-automation/services/video-streaming-hub-service/internal/app"
	"log"
	"time"
)

func main() {
	if err := app.Run(
		app.Config{
			//DatabaseAddr: "192.168.2.135:5432",
			DatabaseAddr: "postgres.default.svc.cluster.local:5432",
			//DatabaseAddr:     "localhost:5432",
			DatabaseUser:     "user",
			DatabasePassword: "password",
			DatabaseName:     "camera_database",
			ApiAddr:          ":80",
			StreamInterval:   time.Millisecond * 30,
		}); err != nil {
		log.Fatal(err)
	}
}
