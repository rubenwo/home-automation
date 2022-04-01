package api

import (
	"github.com/go-chi/chi"
	"github.com/rubenwo/home-automation/services/qr-service/internal/api/config"
	"github.com/rubenwo/home-automation/services/qr-service/internal/pkg/qr"
	"image/jpeg"
	"log"
	"net/http"
	"os"
)

type api struct {
	cfg       config.Config
	generator qr.Generator
}

func New(cfg config.Config) (*api, error) {
	logoFile, err := os.Open("./logo.png")
	if err != nil {
		return nil, err
	}
	defer logoFile.Close()

	logo, err := jpeg.Decode(logoFile)
	if err != nil {
		return nil, err
	}

	generator := qr.NewGenerator(logo, 0)
	return &api{cfg: cfg, generator: generator}, nil
}

func (a *api) Run() error {
	router := chi.NewRouter()

	router.Post("/generate", a.generateQRCode)
	router.Get("/render/{uuid}.{ext}", a.renderQRCode)

	if err := http.ListenAndServe(a.cfg.Listen, router); err != nil {
		return err
	}
	return nil
}

func (a *api) generateQRCode(w http.ResponseWriter, r *http.Request) {

}

func castExt(ext string) qr.OutputFormat {
	switch qr.OutputFormat(ext) {
	case qr.SVG:
		return qr.SVG
	case qr.JPG:
		return qr.JPG
	case qr.PNG:
		return qr.PNG
	}
	return qr.PNG
}

func (a *api) renderQRCode(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "uuid")
	log.Printf("received QR render request for id: %s\n", id)
	ext := chi.URLParam(r, "ext")

	b, err := a.generator.GenerateQRCode("some content here", 2000, 1, castExt(ext), qr.WithAlphaBlending)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Internal Server Error"))
	}

	var contentType string
	switch castExt(ext) {
	case qr.SVG:
		contentType = "image/svg+xml"
	case qr.JPG:
		contentType = "image/jpg"
	case qr.PNG:
		contentType = "image/png"
	}

	w.Header().Set("Content-Type", contentType)
	_, _ = b.WriteTo(w)
}
