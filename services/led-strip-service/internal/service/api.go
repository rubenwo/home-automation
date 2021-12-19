package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-pg/pg/v10"
	ledstrip "github.com/rubenwo/home-automation/services/led-strip-service/pkg/led-strip"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type api struct {
	router         *chi.Mux
	db             *pg.DB
	ledStripClient *ledstrip.Client
}

func New(cfg *Config) (*api, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config can't be nil")
	}
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("error validating config: %w", err)
	}

	db := pg.Connect(&pg.Options{
		Addr:     cfg.DatabaseBackend,
		User:     cfg.DatabaseUser,
		Password: cfg.DatabasePassword,
		Database: cfg.DatabaseName,
	})

	if err := db.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	if err := createSchema(db); err != nil {
		return nil, fmt.Errorf("couldn't create schema: %w", err)
	}

	a := &api{
		router: chi.NewRouter(),
		db:     db,
	}

	var deviceInfos []LedDeviceModel
	if err := a.db.Model(&deviceInfos).Select(); err != nil {
		log.Println(err)
	}

	knownDevices := make(map[string]string)
	if deviceInfos != nil {
		for _, info := range deviceInfos {
			knownDevices[info.ID] = info.Name
		}
	}
	ledStripClient, err := ledstrip.NewClient(cfg.MqttHost, 10, knownDevices)
	if err != nil {
		return nil, err
	}
	a.ledStripClient = ledStripClient

	a.ledStripClient.SetOnLedStripOnlineCallback(func(id, name string) {
		log.Printf("a new led strip came online with id: %s\n", id)

		var deviceInfo LedDeviceModel
		if err := a.db.Model(&deviceInfo).Where("led_device_model.id = ?", id).Select(); err != nil {
			result, err := a.db.Model(&LedDeviceModel{ID: id, Name: name}).Insert()
			if err != nil {
				log.Printf("error storing new led strip with id: %s\n", id)
				return
			}
			fmt.Println(result)
		} else {
			result, err := a.db.Model(&LedDeviceModel{ID: id, Name: name}).Where("led_device_model.id = ?", id).Update(&LedDeviceModel{ID: id, Name: name})
			if err != nil {
				log.Printf("error storing existing led strip with id: %s\n", id)
				return
			}
			fmt.Println(result)
		}
		if err := pushToRegistry(id, name); err != nil {
			log.Printf("error pushing new device to registry: %s\n", err.Error())
		}
	})

	// A good base middleware stack
	a.router.Use(middleware.RequestID)
	a.router.Use(middleware.RealIP)
	a.router.Use(middleware.Logger)
	a.router.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	a.router.Use(middleware.Timeout(60 * time.Second))

	a.router.Get("/healthz", a.healthz)
	a.router.Get("/leds/devices", a.getDevices)
	a.router.Get("/leds/devices/{id}", a.getDevice)
	a.router.Post("/leds/devices/{id}/command/{mode}", a.commandDevice)
	a.router.Post("/leds/devices/command/{mode}", a.commandDevices)

	return a, nil
}

func pushToRegistry(id, name string) error {
	c := &http.Client{}

	req, err := http.NewRequest("DELETE", "http://registry.default.svc.cluster.local/devices/"+id, nil)
	if err != nil {
		return err
	}
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err := resp.Body.Close(); err != nil {
		return err
	}

	var data struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Category string `json:"category"`
		Product  struct {
			Company string `json:"company"`
			Type    string `json:"type"`
		} `json:"product"`
	}
	data.ID = id
	data.Name = name
	data.Category = "led-strip"
	data.Product.Company = "esp32"
	data.Product.Type = "RGB_LED_STRIP"
	jsonData, err := json.Marshal(&data)
	if err != nil {
		return err
	}

	req, err = http.NewRequest("POST", "http://registry.default.svc.cluster.local/devices", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	resp, err = c.Do(req)
	if err != nil {
		return err
	}
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err := resp.Body.Close(); err != nil {
		return err
	}
	return nil
}

func (a *api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}

func (a *api) healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(&HealthzModel{
		IsHealthy:    true,
		ErrorMessage: "",
	}); err != nil {
		log.Printf("error sending healthz: %s\n", err.Error())
	}
}

func (a *api) getDevice(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "no id provided", http.StatusBadRequest)
		return
	}

	var deviceInfo LedDeviceModel
	if err := a.db.Model(&deviceInfo).Where("led_device_model.id = ?", id).Select(); err != nil {
		http.Error(w, fmt.Sprintf("couldn't load item from database: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	var resp struct {
		Device LedDeviceModel `json:"device"`
	}
	resp.Device = deviceInfo

	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		log.Printf("error sending devices: %s\n", err.Error())
	}
}

func (a *api) getDevices(w http.ResponseWriter, r *http.Request) {
	var deviceInfos []LedDeviceModel
	if err := a.db.Model(&deviceInfos).Select(); err != nil {
		http.Error(w, fmt.Sprintf("couldn't load item from database: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	if deviceInfos == nil {
		deviceInfos = []LedDeviceModel{}
	}

	var resp struct {
		Devices []LedDeviceModel `json:"devices"`
	}
	resp.Devices = deviceInfos

	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		log.Printf("error sending devices: %s\n", err.Error())
	}
}

func (a *api) commandDevice(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "no id provided", http.StatusBadRequest)
		return
	}

	mode := chi.URLParam(r, "mode")
	if mode == "" {
		http.Error(w, "no mode provided", http.StatusBadRequest)
		return
	}

	if mode == "solid" {
		var msg struct {
			Mode  string `json:"mode"`
			Red   int    `json:"red"`
			Green int    `json:"green"`
			Blue  int    `json:"blue"`
		}

		if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
			http.Error(w, fmt.Sprintf("couldn't decode the body: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		if err := a.ledStripClient.SetSolidColorById(ledstrip.Color{R: msg.Red, G: msg.Green, B: msg.Blue}, id); err != nil {
			http.Error(w, fmt.Sprintf("couldn't set the solid color for led strip with id: %s, %s", id, err.Error()), http.StatusInternalServerError)
			return
		}
	} else if mode == "colorcycle" {
		if err := a.ledStripClient.SetAnimationColorCycleById(id); err != nil {
			http.Error(w, fmt.Sprintf("couldn't set the animation color for led strip with id: %s, %s", id, err.Error()), http.StatusInternalServerError)
			return
		}

	} else if mode == "christmas" {
		if err := a.ledStripClient.SetAnimationChristmasById(id); err != nil {
			http.Error(w, fmt.Sprintf("couldn't set the animation color for led strip with id: %s, %s", id, err.Error()), http.StatusInternalServerError)
			return
		}
	} else if mode == "breathing" {
		var msg struct {
			Red   int `json:"red"`
			Green int `json:"green"`
			Blue  int `json:"blue"`
		}

		if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
			http.Error(w, fmt.Sprintf("couldn't decode the body: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		if err := a.ledStripClient.SetAnimationBreathingById(ledstrip.Color{R: msg.Red, G: msg.Green, B: msg.Blue}, id); err != nil {
			http.Error(w, fmt.Sprintf("couldn't set the breathing color for led strip with id: %s, %s", id, err.Error()), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "invalid mode", http.StatusBadRequest)
		return
	}
}

func (a *api) commandDevices(w http.ResponseWriter, r *http.Request) {
	mode := chi.URLParam(r, "mode")
	if mode == "" {
		http.Error(w, "no mode provided", http.StatusBadRequest)
		return
	}

	if mode == "solid" {
		var msg struct {
			Mode  string `json:"mode"`
			Red   int    `json:"red"`
			Green int    `json:"green"`
			Blue  int    `json:"blue"`
		}

		if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
			http.Error(w, fmt.Sprintf("couldn't decode the body: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		if err := a.ledStripClient.SetSolidColor(ledstrip.Color{R: msg.Red, G: msg.Green, B: msg.Blue}); err != nil {
			http.Error(w, fmt.Sprintf("couldn't set the solid color for led strips: %s", err.Error()), http.StatusInternalServerError)
			return
		}
	} else if mode == "colorcycle" {
		if err := a.ledStripClient.SetAnimationColorCycle(); err != nil {
			http.Error(w, fmt.Sprintf("couldn't set the animation color for led strips %s", err.Error()), http.StatusInternalServerError)
			return
		}
	} else if mode == "breathing" {
		var msg struct {
			Red   int `json:"red"`
			Green int `json:"green"`
			Blue  int `json:"blue"`
		}

		if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
			http.Error(w, fmt.Sprintf("couldn't decode the body: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		if err := a.ledStripClient.SetAnimationBreathing(ledstrip.Color{R: msg.Red, G: msg.Green, B: msg.Blue}); err != nil {
			http.Error(w, fmt.Sprintf("couldn't set the breathing color for led strips %s", err.Error()), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "invalid mode", http.StatusBadRequest)
		return
	}
}
