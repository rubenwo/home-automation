package ingress

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/rubenwo/home-automation/gateway-service/pkg/mqtt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

const (
	apiPrefix  = "/api/v1"
	authPrefix = "/auth"
	docsPrefix = "/docs/v1"
)

type Ingress struct {
	router     *mux.Router
	mqttClient *mqtt.Client
}

func New(cfg *Config) (*Ingress, error) {
	if cfg == nil {
		return nil, fmt.Errorf("cfg cannot be nil")
	}
	router := mux.NewRouter()
	mqttClient := mqtt.New()

	apiRouter := router.PathPrefix(apiPrefix).Subrouter()
	apiRouter.Use(LoggingMiddleware)
	switch cfg.ApiVersion {
	case "v1":
		log.Println("Using v1 ingress spec")
		for _, spec := range cfg.Spec {
			switch strings.ToUpper(spec.Protocol) {
			case "HTTP":
				for index, spec := range cfg.Spec {
					target := strings.ToLower(spec.Protocol) + "://" + spec.Host
					u, err := url.Parse(target)
					if err != nil {
						return nil, fmt.Errorf("invalid spec at index: %d, parsing error: %w", index, err)
					}
					apiRouter.HandleFunc(spec.Path, func(writer http.ResponseWriter, request *http.Request) {
						request.URL.Path = strings.TrimPrefix(request.URL.Path, apiPrefix)
						serveReverseProxy(u, writer, request)
					}).Methods(spec.Methods...)
				}
			case "MQTT":
				if err := mqttClient.Register(spec.Path, spec.Host, 10); err != nil {
					log.Println(err)
				}
				apiRouter.HandleFunc(spec.Path, func(writer http.ResponseWriter, request *http.Request) {
					request.URL.Path = strings.TrimPrefix(request.URL.Path, apiPrefix)
					mqttClient.BrokerMQTTRequest(writer, request)
				}).Methods("POST")
			default:
				return nil, fmt.Errorf("%s is not a supported spec.Protocol", spec.Protocol)
			}
		}

	default:
		return nil, fmt.Errorf("%s is not a supported ApiVersion", cfg.ApiVersion)
	}

	return &Ingress{router: router, mqttClient: mqttClient}, nil
}

func (i *Ingress) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	i.router.ServeHTTP(w, r)
}

func (i *Ingress) Use(mfw ...mux.MiddlewareFunc) {
	i.router.Use(mfw...)
}

func (i *Ingress) Reload() {
	panic("Reload() not implemented")
}

func serveReverseProxy(u *url.URL, w http.ResponseWriter, r *http.Request) {
	proxy := httputil.NewSingleHostReverseProxy(u)

	r.URL.Host = u.Host
	r.URL.Scheme = u.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Host = u.Host

	r.Header.Set("X-Request-Id", uuid.New().String()) // Inject a request-id for tracing

	proxy.ServeHTTP(w, r)
}
