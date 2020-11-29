package ingress

import (
	"fmt"
	"log"
	"net/http"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(fmt.Sprintf("Gateway => %s request from: %s using %s, to: %s", r.Method, r.RemoteAddr, r.Proto, r.RequestURI))
		next.ServeHTTP(w, r)
	})
}
