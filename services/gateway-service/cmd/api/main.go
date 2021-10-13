package main

import (
	"crypto/tls"
	"flag"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/rubenwo/home-automation/services/gateway-service/pkg/auth"
	"github.com/rubenwo/home-automation/services/gateway-service/pkg/ingress"
	"golang.org/x/crypto/acme/autocert"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	ingressPathPtr := flag.String("file", "/root/ingress.yaml", "specifies the filepath for the ingress config")
	flag.Parse()

	jwtKey := os.Getenv("JWT_KEY")

	if jwtKey == "" {
		log.Fatal("jwt key can't be empty")
	}

	adminEnabled := os.Getenv("ENABLE_ADMIN") == "true"
	isProductionEnv := os.Getenv("PRODUCTION") == "true"

	cfg, err := ingress.ParseConfig(*ingressPathPtr)
	if err != nil {
		log.Fatal(err)
	}

	isAllowedAnonymous := make(map[string]bool)

	for _, spec := range cfg.Spec {
		for _, route := range spec.Routes {
			isAllowedAnonymous["/api/"+cfg.ApiVersion+route.Path] = route.AllowAnonymous
		}
	}

	authenticator := auth.NewDefaultClient([]byte(jwtKey), time.Hour*1, time.Hour*24*7, adminEnabled, isAllowedAnonymous)

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

	if isProductionEnv {
		domainsEnv := os.Getenv("DOMAINS")

		if domainsEnv == "" {
			log.Fatal("In production environment, but no domains specified")
		}
		domains := strings.Split(domainsEnv, ",")
		if len(domains) == 0 {
			log.Fatal("In production environment, but no domains specified")
		}

		// create the autocert.Manager with domains and path to the cache
		certManager := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(domains...),
			Cache:      autocert.DirCache("/certs"),
		}

		// Create the http server
		server := &http.Server{
			Addr:    ":443",
			Handler: handler,
			TLSConfig: &tls.Config{
				// Causes servers to use Go's default ciphersuite preferences,
				// which are tuned to avoid attacks. Does nothing on clients.
				PreferServerCipherSuites: true,
				// Only use curves which have assembly implementations
				CurvePreferences: []tls.CurveID{
					tls.CurveP256,
					tls.X25519, // Go 1.8 only
				},
				MinVersion: tls.VersionTLS12,
				CipherSuites: []uint16{
					tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305, // Go 1.8 only
					tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,   // Go 1.8 only
					tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				},
				GetCertificate: certManager.GetCertificate,
			},
		}
		log.Printf("Serving http/https for domains: %+v\n", domains)
		go func() {
			// serve HTTP, which will redirect automatically to HTTPS
			h := certManager.HTTPHandler(nil)
			log.Fatal(http.ListenAndServe(":80", h))
		}()

		// serve HTTPS!
		if err := server.ListenAndServeTLS("", ""); err != nil {
			log.Fatal(err)
		}
		return
	}

	// If certificate exists in non-production environment, host on 443
	if _, err := os.Stat("/certs/fullchain.pem"); err == nil {
		// path/to/whatever exists
		log.Println("gateway-service is listening on port '443'")
		if err := http.ListenAndServeTLS(":443", "/certs/fullchain.pem", "/certs/privkey.pem", handler); err != nil {
			log.Fatal(err)
		}
		return
	}

	// else host on port 80
	log.Println("gateway-service is listening on port '80'")
	if err := http.ListenAndServe(":80", handler); err != nil {
		log.Fatal(err)
	}

}
