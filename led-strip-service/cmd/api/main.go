package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/rubenwo/home-automation/led-strip-service/pkg/database"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type api struct {
	db database.Database
}

func main() {
	router := chi.NewRouter()
	db, err := database.Factory("redis")
	if err != nil {
		log.Fatalf("couldn't create database: %v", err)
	}

	a := &api{db: db}

	// A good base middleware stack
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	router.Use(middleware.Timeout(60 * time.Second))

	router.Get("/healthz", a.healthz)
	router.Get("/leds/devices", a.getDevices)
	router.Post("/leds/devices/register", a.registerDevice)
	router.Post("/leds/devices/{id}/command", a.commandDevice)

	log.Println("led-strip-service is online")
	if err := http.ListenAndServe(":80", router); err != nil {
		log.Println(err)
	}
	log.Println("led-strip-service is offline")
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

func (a *api) getDevices(w http.ResponseWriter, r *http.Request) {
	v, err := a.db.Get("led-strip-devices")
	if err != nil {
		w.Header().Set("content-type", "application/json")
		if err := json.NewEncoder(w).Encode(&JsonError{
			Code:         http.StatusInternalServerError,
			ErrorMessage: err.Error(),
		}); err != nil {
			log.Printf("error sending json_error: %s\n", err.Error())
		}
		return
	}
	var resp struct {
		Devices []LedDeviceModel `json:"devices"`
	}
	if v == nil {
		resp.Devices = make([]LedDeviceModel, 0)
		w.Header().Set("content-type", "application/json")
		if err := json.NewEncoder(w).Encode(&resp); err != nil {
			log.Printf("error sending getDevices: %s\n", err.Error())
		}
		return
	}

	vStr := v.(string)
	var devices []LedDeviceModel
	if err := json.Unmarshal([]byte(vStr), &devices); err != nil {
		log.Println(err)
		w.Header().Set("content-type", "application/json")
		if err := json.NewEncoder(w).Encode(&JsonError{
			Code:         http.StatusInternalServerError,
			ErrorMessage: err.Error(),
		}); err != nil {
			log.Printf("error sending json_error: %s\n", err.Error())
		}
		return
	}

	resp.Devices = devices
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		log.Printf("error sending getDevices: %s\n", err.Error())
	}
}

func (a *api) registerDevice(w http.ResponseWriter, r *http.Request) {
	var d RegisterLedDeviceModel
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	c := &http.Client{}
	resp, err := c.Get(fmt.Sprintf("http://%s/healthz", d.IPAddress))
	if err != nil {
		w.Header().Set("content-type", "application/json")
		if err := json.NewEncoder(w).Encode(&JsonError{
			Code:         http.StatusInternalServerError,
			ErrorMessage: err.Error(),
		}); err != nil {
			log.Printf("error sending json_error: %s\n", err.Error())
		}
		return
	}
	var ledHealthz HealthzModel
	if err := json.NewDecoder(resp.Body).Decode(&ledHealthz); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	fmt.Println(ledHealthz)
	resp.Body.Close()

	v, err := a.db.Get("led-strip-devices")
	fmt.Println(v, err)
	if err != nil && err != redis.Nil {
		log.Println(err)
		w.Header().Set("content-type", "application/json")
		if err := json.NewEncoder(w).Encode(&JsonError{
			Code:         http.StatusInternalServerError,
			ErrorMessage: err.Error(),
		}); err != nil {
			log.Printf("error sending json_error: %s\n", err.Error())
		}
		return
	}

	resp, err = c.Get(fmt.Sprintf("http://%s/info", d.IPAddress))
	if err != nil {
		w.Header().Set("content-type", "application/json")
		if err := json.NewEncoder(w).Encode(&JsonError{
			Code:         http.StatusInternalServerError,
			ErrorMessage: err.Error(),
		}); err != nil {
			log.Printf("error sending json_error: %s\n", err.Error())
		}
		return
	}
	var ledControllerInfo LedControllerInfo
	if err := json.NewDecoder(resp.Body).Decode(&ledControllerInfo); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	defer resp.Body.Close()

	// Add device to database in this service
	var devices []LedDeviceModel
	id := uuid.New().String()
	if v == nil {
		devices = []LedDeviceModel{
			{
				ID:             id,
				Name:           ledControllerInfo.DeviceName,
				NumLeds:        150,
				SupportedModes: []string{"SINGLE_COLOR_RGB", "SINGLE_COLOR_HSV", "GRADIENT_RGB", "GRADIENT_HSV", "ANIMATION_RGB", "ANIMATION_HSV"},
				CurrentMode:    ledControllerInfo.DeviceInfo.Mode,
				IPAddress:      d.IPAddress,
			},
		}
	} else {
		vStr := v.(string)
		if err := json.Unmarshal([]byte(vStr), &devices); err != nil {
			log.Println(err)
			w.Header().Set("content-type", "application/json")
			if err := json.NewEncoder(w).Encode(&JsonError{
				Code:         http.StatusInternalServerError,
				ErrorMessage: err.Error(),
			}); err != nil {
				log.Printf("error sending json_error: %s\n", err.Error())
			}
			return
		}
		devices = append(devices, LedDeviceModel{
			ID:             id,
			Name:           ledControllerInfo.DeviceName,
			NumLeds:        150,
			SupportedModes: ledControllerInfo.DeviceInfo.SupportedModes,
			CurrentMode:    ledControllerInfo.DeviceInfo.CurrentMode,
			IPAddress:      d.IPAddress,
		})
	}

	jsonData, err := json.Marshal(&devices)
	if err != nil {
		log.Println(jsonData)
	}
	if err := a.db.Set(fmt.Sprintf("led-strip-devices"), jsonData); err != nil {
		log.Println(err)
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
	data.Name = ledControllerInfo.DeviceName
	data.Category = "led-strip"
	data.Product.Company = "esp32"
	data.Product.Type = ledControllerInfo.DeviceType
	jsonData, err = json.Marshal(&data)
	if err != nil {
		if err := json.NewEncoder(w).Encode(&JsonError{
			Code:         http.StatusBadRequest,
			ErrorMessage: "id can't be empty",
		}); err != nil {
			log.Printf("error sending json_error: %s\n", err.Error())
		}
		return
	}

	req, err := http.NewRequest("POST", "http://registry.default.svc.cluster.local/devices", bytes.NewBuffer(jsonData))
	if err != nil {
		if err := json.NewEncoder(w).Encode(&JsonError{
			Code:         http.StatusBadRequest,
			ErrorMessage: "id can't be empty",
		}); err != nil {
			log.Printf("error sending json_error: %s\n", err.Error())
		}
		return
	}
	resp, err = c.Do(req)
	if err != nil {
		if err := json.NewEncoder(w).Encode(&JsonError{
			Code:         http.StatusBadRequest,
			ErrorMessage: "id can't be empty",
		}); err != nil {
			log.Printf("error sending json_error: %s\n", err.Error())
		}
		return
	}
	raw, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Println(string(raw), err)
	var re struct {
		Devices []LedDeviceModel `json:"devices"`
	}
	re.Devices = devices
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(&re); err != nil {
		log.Printf("error sending getDevices: %s\n", err.Error())
	}
}

func (a *api) commandDevice(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		if err := json.NewEncoder(w).Encode(&JsonError{
			Code:         http.StatusBadRequest,
			ErrorMessage: "id can't be empty",
		}); err != nil {
			log.Printf("error sending json_error: %s\n", err.Error())
		}
		return
	}
	v, err := a.db.Get("led-strip-devices")
	if err != nil {
		w.Header().Set("content-type", "application/json")
		if err := json.NewEncoder(w).Encode(&JsonError{
			Code:         http.StatusInternalServerError,
			ErrorMessage: err.Error(),
		}); err != nil {
			log.Printf("error sending json_error: %s\n", err.Error())
		}
		return
	}
	fmt.Println(v)

	if v == nil {
		http.NotFound(w, r)
		return
	}

	vStr := v.(string)
	var devices []LedDeviceModel
	if err := json.Unmarshal([]byte(vStr), &devices); err != nil {
		log.Println(err)
		w.Header().Set("content-type", "application/json")
		if err := json.NewEncoder(w).Encode(&JsonError{
			Code:         http.StatusInternalServerError,
			ErrorMessage: err.Error(),
		}); err != nil {
			log.Printf("error sending json_error: %s\n", err.Error())
		}
		return
	}
	fmt.Println(devices)
	for _, dev := range devices {
		if dev.ID == id {
			fmt.Println("sending command")
			raw, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.Header().Set("content-type", "application/json")
				if err := json.NewEncoder(w).Encode(&JsonError{
					Code:         http.StatusInternalServerError,
					ErrorMessage: err.Error(),
				}); err != nil {
					log.Printf("error sending json_error: %s\n", err.Error())
				}
				return
			}
			fmt.Println(string(raw))
			proxyReq, err := http.NewRequest("POST", fmt.Sprintf("http://%s/led", dev.IPAddress), bytes.NewBuffer(raw))
			if err != nil {
				w.Header().Set("content-type", "application/json")
				if err := json.NewEncoder(w).Encode(&JsonError{
					Code:         http.StatusInternalServerError,
					ErrorMessage: err.Error(),
				}); err != nil {
					log.Printf("error sending json_error: %s\n", err.Error())
				}
				return
			}

			proxyReq.Header.Set("Host", r.Host)
			proxyReq.Header.Set("X-Forwarded-For", r.RemoteAddr)

			client := &http.Client{}
			proxyRes, err := client.Do(proxyReq)
			if err != nil {
				log.Println(err)
				w.Header().Set("content-type", "application/json")
				if err := json.NewEncoder(w).Encode(&JsonError{
					Code:         http.StatusInternalServerError,
					ErrorMessage: err.Error(),
				}); err != nil {
					log.Printf("error sending json_error: %s\n", err.Error())
				}
				return
			}
			var c CommandResponseModel
			if err := json.NewDecoder(proxyRes.Body).Decode(&c); err != nil {
				log.Println(err)
				w.Header().Set("content-type", "application/json")
				if err := json.NewEncoder(w).Encode(&JsonError{
					Code:         http.StatusInternalServerError,
					ErrorMessage: err.Error(),
				}); err != nil {
					log.Printf("error sending json_error: %s\n", err.Error())
				}
				return
			}
			fmt.Println(c)
			proxyRes.Body.Close()
			w.Header().Set("content-type", "application/json")
			if err := json.NewEncoder(w).Encode(&c); err != nil {
				log.Printf("error sending response: %v", err)
			}
			return
		}
	}
	http.NotFound(w, r)
}
