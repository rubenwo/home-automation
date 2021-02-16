package main

import (
	"fmt"
	"github.com/jasonlvhit/gocron"
	"github.com/rubenwo/home-automation/registry-service/pkg/registry"
	"log"
	"net/http"
)

func main() {
	go func() {
		for <-gocron.Start() {
		}
	}()

	scheduler := &registry.Scheduler{}
	if err := scheduler.CreateSchedule(registry.Schedule{
		Every:  1,
		At:     "15:38",
		Script: fmt.Sprint("turn_on>http://tapo.default.svc.cluster.local/tapo/lights/dfa20eb0-d9e4-40d0-a6ae-19eb5060a1fc?command=off&brightness=100\n"),
	}); err != nil {
		log.Fatal(err)
	}
	scheduler.CreateSchedule(registry.Schedule{
		Every:  1,
		At:     "15:39",
		Script: fmt.Sprint("turn_on>http://tapo.default.svc.cluster.local/tapo/lights/dfa20eb0-d9e4-40d0-a6ae-19eb5060a1fc?command=om&brightness=100\n"),
	})

	router, err := registry.New(&registry.Config{DatabaseBackend: "boltdb"})
	if err != nil {
		log.Fatal(err)
	}
	if err := http.ListenAndServe(":80", router); err != nil {
		log.Fatal(err)
	}
}
