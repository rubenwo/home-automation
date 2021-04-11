package main

import (
	"github.com/rubenwo/home-automation/actions-service/internal/actions"
	"log"
)

func main() {
	if err := actions.Run(&actions.Config{
		Addr: ":80",
		IntentProcessors: map[string]actions.IntentProcessor{
			"Turn on":       &actions.ProcessTurnOnRequest{},
			"Turn off":      &actions.ProcessTurnOffRequest{},
			"Change colour": &actions.ProcessChangeColourRequest{},
		},
	}); err != nil {
		log.Fatal(err)
	}
}
