package actions

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/rubenwo/home-automation/services/actions-service/internal/actions/intentprocessor"
	"github.com/rubenwo/home-automation/services/actions-service/internal/actions/models"
	"log"
	"net/http"
)

type api struct {
	intentProcessors map[string]intentprocessor.IntentProcessor
}

func Run(cfg *Config) error {
	if cfg == nil {
		return ConfigNilError
	}
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("%v: %w", ConfigValidationError, err)
	}

	a := &api{intentProcessors: cfg.IntentProcessors}

	router := chi.NewRouter()

	router.Get("/healthz", a.healthz)
	router.Post("/webhook/request", a.onWebhookRequest)

	if err := http.ListenAndServe(cfg.Addr, router); err != nil {
		return err
	}

	return nil
}

func (a *api) healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(&models.HealthzModel{
		IsHealthy:    true,
		ErrorMessage: "",
	}); err != nil {
		log.Printf("error sending healthz: %s\n", err.Error())
	}
}

func (a *api) onWebhookRequest(w http.ResponseWriter, r *http.Request) {
	var req models.WebhookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}
	fmt.Println(req)
	fmt.Println(req.QueryResult.Intent.DisplayName)
	processor, ok := a.intentProcessors[req.QueryResult.Intent.DisplayName]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("processor for request not found")
		return
	}

	msg, err := processor.ProcessIntent(req.QueryResult.Parameters)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := models.TextResponse{FulfillmentMessages: []models.Message{
		models.Message{Text: models.Text{Text: []string{msg}}},
	}}
	if err := json.NewEncoder(w).Encode(&response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
}
