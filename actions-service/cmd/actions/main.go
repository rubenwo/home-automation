package main

import (
	"github.com/rubenwo/home-automation/actions-service/internal/actions"
	"github.com/rubenwo/home-automation/actions-service/internal/actions/intentprocessor"
	"log"
)

func main() {
	if err := actions.Run(&actions.Config{
		Addr: ":80",
		IntentProcessors: map[string]intentprocessor.IntentProcessor{
			"Turn on":       &intentprocessor.ProcessTurnOnRequest{},
			"Turn off":      &intentprocessor.ProcessTurnOffRequest{},
			"Change colour": &intentprocessor.ProcessChangeColourRequest{},
		},
	}); err != nil {
		log.Fatal(err)
	}
}