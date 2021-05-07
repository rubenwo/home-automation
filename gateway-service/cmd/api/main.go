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

	globalMiddlewares := []mux.MiddlewareFunc{
		ingress.LoggingMiddleware,
	}

	apiMiddlewares := []mux.MiddlewareFunc{
		authenticator.AuthorizationMiddleware,
	}

	router, err := ingress.New(cfg, authenticator, globalMiddlewares, apiMiddlewares)
	if err != nil {
		log.Fatal(err)
	}

	handler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:8080", "http://localhost", "https://homeautomation.rubenwoldhuis.nl"},
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

	// If certificate exists, host on 443
	if _, err := os.Stat("/certs/fullchain.pem"); err == nil {
		// path/to/whatever exists
		log.Println("gateway-service is listening on port '443'")
		if err := http.ListenAndServeTLS(":443", "/certs/fullchain.pem", "/certs/privkey.pem", handler); err != nil {
			log.Fatal(err)
		}
		return
	}
	log.Println("gateway-service is listening on port '80'")
	if err := http.ListenAndServe(":80", handler); err != nil {
		log.Fatal(err)
	}
}
