package registry

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/google/uuid"
	"github.com/rubenwo/home-automation/libraries/go/pkg/database"
	"log"
	"net/http"
	"time"
)

type api struct {
	router    *chi.Mux
	devices   map[string]DeviceInfo
	keys      []string
	scheduler *Scheduler
	groups    map[string][]string
	db        database.Database
}

func New(cfg *Config) (*api, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config can't be nil")
	}
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("error validating config: %w", err)
	}
	db, err := database.Factory(cfg.DatabaseBackend)
	if err != nil {
		return nil, fmt.Errorf("couldn't create database: %w", err)
	}

	a := &api{
		router:    chi.NewRouter(),
		devices:   make(map[string]DeviceInfo),
		scheduler: NewScheduler(),
		keys:      []string{},
		groups:    make(map[string][]string),
		db:        db,
	}

	rawKeys, err := db.Get("registry-keys")
	if err != nil {
		fmt.Printf("error retrieving keys: %v\n", err)
	} else {
		fmt.Println(rawKeys)
		rawKeysStr := rawKeys.(string)
		//rawKeysStr = strings.ReplaceAll(rawKeysStr, "\\", "")
		//rawKeysStr = strings.ReplaceAll(rawKeysStr, "\"", "")
		//rawKeysStr = strings.TrimPrefix(rawKeysStr, "[")
		//rawKeysStr = strings.TrimSuffix(rawKeysStr, "]")
		//
		//splitKeys := strings.Split(rawKeysStr, ",")
		//for _, sp := range splitKeys {
		//	a.keys = append(a.keys, strings.TrimSpace(sp))
		//}

		if err := json.Unmarshal([]byte(rawKeysStr), &a.keys); err != nil {

		}

		for _, k := range a.keys {
			fmt.Println(k)
			raw, err := a.db.Get(k)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(raw)
			var dev DeviceInfo
			if err := json.Unmarshal([]byte(raw.(string)), &dev); err != nil {
				log.Fatal(err)
			}
			fmt.Println(dev)
			a.devices[k] = dev
		}
	}
	// A good base middleware stack
	a.router.Use(middleware.RequestID)
	a.router.Use(middleware.RealIP)
	a.router.Use(middleware.Logger)
	a.router.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	a.router.Use(middleware.Timeout(60 * time.Second))

	a.router.Get("/new_id", a.getNewID)

	a.router.Get("/sensors", a.getSensors)
	a.router.Post("/sensors", a.addSensor)
	a.router.Delete("/sensors/{id}", a.deleteSensor)

	a.router.Get("/devices", a.getDevices)
	a.router.Post("/devices", a.postDevice)
	a.router.Delete("/devices/{id}", a.deleteDevice)

	a.router.Get("/schedules", a.getSchedules)
	a.router.Post("/schedules", a.createSchedule)

	a.router.Post("/group", a.createGroup)

	return a, nil
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

func (a *api) getNewID(w http.ResponseWriter, r *http.Request) {
	var resp struct {
		ID string `json:"id"`
	}
	resp.ID = uuid.New().String()

	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		log.Printf("error sending getDevices: %s\n", err.Error())
	}
}

func (a *api) getSensors(w http.ResponseWriter, r *http.Request)   {}
func (a *api) addSensor(w http.ResponseWriter, r *http.Request)    {}
func (a *api) deleteSensor(w http.ResponseWriter, r *http.Request) {}

func (a *api) postDevice(w http.ResponseWriter, r *http.Request) {
	var dev DeviceInfo
	if err := json.NewDecoder(r.Body).Decode(&dev); err != nil {
		log.Println("error decoding the body:", err)
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	a.devices[dev.ID] = dev
	jsonData, err := json.Marshal(&dev)
	if err != nil {
		log.Println(jsonData)
	}
	if err := a.db.Set(fmt.Sprintf("registry-%s", dev.ID), jsonData); err != nil {
		log.Println(err)
	}
	a.keys = append(a.keys, fmt.Sprintf("registry-%s", dev.ID))
	jsonKeys, err := json.Marshal(&a.keys)
	if err != nil {
		log.Println(jsonKeys)
	}
	if err := a.db.Set("registry-keys", jsonKeys); err != nil {
		log.Println(err)
	}

	a.getDevices(w, r)
}

func (a *api) getDevices(w http.ResponseWriter, r *http.Request) {
	var devices struct {
		Devices []DeviceInfo `json:"devices"`
	}
	devices.Devices = make([]DeviceInfo, 0)
	for _, v := range a.devices {
		devices.Devices = append(devices.Devices, v)
	}

	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(&devices); err != nil {
		log.Printf("error sending getDevices: %s\n", err.Error())
	}
}

func (a *api) deleteDevice(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.NotFound(w, r)
		return
	}
	fmt.Println(id)
	fmt.Println(a.devices)
	delete(a.devices, id)
	fmt.Println(a.devices)
	if err := a.db.Delete(fmt.Sprintf("registry-%s", id)); err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("error deleting device: %s", err.Error()), http.StatusInternalServerError)
	}
	for i, v := range a.keys {
		if v == fmt.Sprintf("registry-%s", id) {
			a.keys = append(a.keys[:i], a.keys[i+1:]...)
			fmt.Println("removed key from registry-keys")
			break
		}
	}
	fmt.Println(a.keys)

	jsonKeys, err := json.Marshal(&a.keys)
	if err != nil {
		log.Println(jsonKeys)
	}
	if err := a.db.Set("registry-keys", jsonKeys); err != nil {
		log.Println(err)
	}
	var msg struct {
		Message string `json:"message"`
	}
	msg.Message = fmt.Sprintf("deleted device: %s successfully", id)
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(&msg); err != nil {
		log.Printf("error sending deleteDevice msg: %s\n", err.Error())
	}
}

func (a *api) createSchedule(w http.ResponseWriter, r *http.Request) {
	a.scheduler.CreateSchedule(Schedule{})
}

func (a *api) getSchedules(w http.ResponseWriter, r *http.Request) {

}

func (a *api) createGroup(w http.ResponseWriter, r *http.Request) {
	var newGroup struct {
		Devices []string `json:"devices"`
	}
	if err := json.NewDecoder(r.Body).Decode(&newGroup); err != nil {
		log.Println("error decoding the body:", err)
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	id := uuid.New().String()
	a.groups[id] = newGroup.Devices

	var createdGroup struct {
		ID      string   `json:"id"`
		Devices []string `json:"devices"`
	}

	createdGroup.ID = id
	createdGroup.Devices = newGroup.Devices

	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(&createdGroup); err != nil {
		log.Printf("error sending createGroup: %s\n", err.Error())
	}
}

func (a *api) getGroups(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Groups map[string][]string `json:"groups"`
	}
	data.Groups = a.groups

	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(&data); err != nil {
		log.Printf("error sending getDevices: %s\n", err.Error())
	}
}
