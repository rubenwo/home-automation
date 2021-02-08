package ingress

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/rubenwo/home-automation/gateway-service/pkg/auth"
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

//New create a new ingress router. The config specifies with back-ends the gateway has. This cannot be nil
//The authenticator is used for the /auth(/login,/logout,/register,/refresh) endpoints.
//The last arguments are apiMiddleware. This is a slice of mux.MiddlewareFuncs that are applied to the /api/v1 endpoints.
func New(cfg *Config, authenticator auth.Authenticator, apiMiddleware ...mux.MiddlewareFunc) (*Ingress, error) {
	if cfg == nil {
		return nil, fmt.Errorf("cfg cannot be nil")
	}

	router := mux.NewRouter()
	mqttClient := mqtt.New()

	authRouter := router.PathPrefix(authPrefix).Subrouter()
	authRouter.HandleFunc("/login", authenticator.Login).Methods("POST")
	authRouter.HandleFunc("/logout", authenticator.Logout).Methods("GET")
	authRouter.HandleFunc("/register", authenticator.Register).Methods("POST")
	authRouter.HandleFunc("/refresh", authenticator.RefreshToken).Methods("GET")

	apiRouter := router.PathPrefix(apiPrefix).Subrouter()
	apiRouter.Use(apiMiddleware...)

	switch cfg.ApiVersion {
	case "v1":
		log.Println("Using v1 ingress spec")
		for index, spec := range cfg.Spec {
			switch strings.ToUpper(spec.Protocol) {
			case "HTTP":
				target := strings.ToLower(spec.Protocol) + "://" + spec.Host
				u, err := url.Parse(target)
				if err != nil {
					return nil, fmt.Errorf("invalid spec at index: %d, parsing error: %w", index, err)
				}
				apiRouter.HandleFunc(spec.Path, func(writer http.ResponseWriter, request *http.Request) {
					request.URL.Path = strings.TrimPrefix(request.URL.Path, apiPrefix)
					serveReverseProxy(u, writer, request)
				}).Methods(spec.Methods...)

			case "MQTT":
				fmt.Println("doing things with mqtt")
				fmt.Println(spec.Methods)
				for _, method := range spec.Methods {
					switch strings.ToUpper(method) {
					case "POST":
						if err := mqttClient.Register(spec.Path, spec.Host, 10); err != nil {
							log.Println(err)
						}
						apiRouter.HandleFunc(spec.Path, func(writer http.ResponseWriter, request *http.Request) {
							request.URL.Path = strings.TrimPrefix(request.URL.Path, apiPrefix)
							fmt.Println("MQTT")
							mqttClient.BrokerMQTTRequest(writer, request)
						}).Methods("POST")
					case "GET":
						fmt.Println("TODO: implement websocket to mqtt listener")
					default:
						return nil, fmt.Errorf("%s is not a supported method for: %s", method, spec.Protocol)
					}
				}
			default:
				return nil, fmt.Errorf("%s is not a supported spec.Protocol", spec.Protocol)
			}
		}

	default:
		return nil, fmt.Errorf("%s is not a supported ApiVersion", cfg.ApiVersion)
	}

	target := "http://web.default.svc.cluster.local"
	u, err := url.Parse(target)
	if err != nil {
		return nil, fmt.Errorf("invalid catch-all target, parsing error: %w", err)
	}
	router.PathPrefix("/").HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		serveReverseProxy(u, writer, request)
	}).Methods("GET")

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
