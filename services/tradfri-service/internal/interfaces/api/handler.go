package api

import (
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
)

func RegisterHandler(router chi.Router) {
	handler := Handler{}

	router.Get("/hello", handler.hello)
}

type Handler struct {
}

func (*Handler) hello(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r)
	_, _ = w.Write([]byte("ok"))
}
