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
			for _, route := range spec.Routes {
				switch strings.ToUpper(route.Protocol) {
				case "HTTP":
					target := strings.ToLower(route.Protocol) + "://" + spec.Host
					u, err := url.Parse(target)
					if err != nil {
						return nil, fmt.Errorf("invalid spec at index: %d, parsing error: %w", index, err)
					}
					apiRouter.HandleFunc(route.Path, func(writer http.ResponseWriter, request *http.Request) {
						request.URL.Path = strings.TrimPrefix(request.URL.Path, apiPrefix)
						serveReverseProxy(u, writer, request)
					}).Methods(route.Methods...)

				case "MQTT":
					fmt.Println("doing things with mqtt")
					fmt.Println(route.Methods)
					for _, method := range route.Methods {
						if err := mqttClient.Register(route.Path, spec.Host, 10); err != nil {
							log.Fatal(err)
						}
						switch strings.ToUpper(method) {
						case "POST":

							apiRouter.HandleFunc(route.Path, func(writer http.ResponseWriter, request *http.Request) {
								request.URL.Path = strings.TrimPrefix(request.URL.Path, apiPrefix)
								fmt.Println("MQTT POST")
								mqttClient.BrokerMQTTRequest(writer, request)
							}).Methods("POST")
						case "GET":
							apiRouter.HandleFunc(route.Path, func(writer http.ResponseWriter, request *http.Request) {
								request.URL.Path = strings.TrimPrefix(request.URL.Path, apiPrefix)
								fmt.Println("MQTT GET")
								mqttClient.SocketMQTTRequest(writer, request)
							}).Methods("GET")
						default:
							return nil, fmt.Errorf("%s is not a supported method for: %s", method, route.Protocol)
						}
					}
				default:
					return nil, fmt.Errorf("%s is not a supported spec.Protocol", route.Protocol)
				}
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

	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		tpl, err1 := route.GetPathTemplate()
		met, err2 := route.GetMethods()
		fmt.Println(tpl, err1, met, err2)
		return nil
	})

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
