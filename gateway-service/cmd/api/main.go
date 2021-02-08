package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/rubenwo/home-automation/gateway-service/pkg/auth"
	"github.com/rubenwo/home-automation/gateway-service/pkg/ingress"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	jwtKey := os.Getenv("JWT_KEY")

	if jwtKey == "" {
		log.Fatal("jwt key can't be empty")
	}

	adminEnabled := false
	if os.Getenv("ENABLE_ADMIN") == "true" {
		adminEnabled = true
	}
	fmt.Println(adminEnabled)
	cfg, err := ingress.ParseConfig("./ingress.yaml")
	if err != nil {
		log.Fatal(err)
	}

	authenticator := auth.NewDefaultClient([]byte(jwtKey), time.Hour*1, adminEnabled)

	mfw := []mux.MiddlewareFunc{
		ingress.LoggingMiddleware,
		authenticator.AuthorizationMiddleware,
	}

	router, err := ingress.New(cfg, authenticator, mfw...)
	if err != nil {
		log.Fatal(err)
	}

	handler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:8080", "http://localhost", "http://192.168.2.135"},
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}).Handler(router)

	if err := http.ListenAndServe(":80", handler); err != nil {
		log.Fatal(err)
	}
}
