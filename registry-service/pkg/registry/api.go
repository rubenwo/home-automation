package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"github.com/rubenwo/home-automation/libraries/go/pkg/database"
	"log"
	"net/http"
	"strconv"
	"time"
)

type api struct {
	router    *chi.Mux
	devices   map[string]DeviceInfo
	keys      []string
	scheduler *Scheduler
	groups    map[string][]string
	db        database.Database
	pgDb      *pg.DB
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

	pgDb := pg.Connect(&pg.Options{
		Addr:     cfg.PgDatabaseBackend,
		User:     cfg.PgDatabaseUser,
		Password: cfg.PgDatabasePassword,
		Database: cfg.PgDatabaseName,
	})

	if err := pgDb.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	if err := createSchema(pgDb); err != nil {
		return nil, fmt.Errorf("couldn't create schema: %w", err)
	}

	a := &api{
		router:    chi.NewRouter(),
		devices:   make(map[string]DeviceInfo),
		scheduler: NewScheduler(),
		keys:      []string{},
		groups:    make(map[string][]string),
		db:        db,
		pgDb:      pgDb,
	}

	rawKeys, err := db.Get("registry-keys")
	if err != nil {
		fmt.Printf("error retrieving keys: %v\n", err)
	} else if rawKeys == nil || rawKeys.(string) == "" {
		fmt.Println("rawKeys is nil")
	} else {
		fmt.Println(rawKeys)
		rawKeysStr := rawKeys.(string)

		fmt.Println(rawKeysStr)

		if err := json.Unmarshal([]byte(rawKeysStr), &a.keys); err != nil {
			return nil, fmt.Errorf("couldn't unmarshal the keys from the database")
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
	a.router.Get("/devices/{id}", a.getDevice)
	a.router.Delete("/devices/{id}", a.deleteDevice)

	a.router.Get("/schedules", a.getSchedules)
	a.router.Post("/schedules", a.createSchedule)
	a.router.Get("/schedules/{id}", a.getSchedule)
	a.router.Delete("/schedules/{id}", a.deleteSchedule)

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

func (a *api) getSensors(w http.ResponseWriter, r *http.Request) {
	var sensors []SensorDevice
	if err := a.pgDb.Model(&sensors).Select(); err != nil {
		http.Error(w, fmt.Sprintf("couldn't load item from database: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	if sensors == nil {
		sensors = []SensorDevice{}
	}

	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(&sensors); err != nil {
		log.Printf("error sending sensors: %s\n", err.Error())
	}
}

func (a *api) addSensor(w http.ResponseWriter, r *http.Request) {
	var sensor SensorDevice
	if err := json.NewDecoder(r.Body).Decode(&sensor); err != nil {
		http.Error(w, fmt.Sprintf("couldn't decode body: %s", err.Error()), http.StatusBadRequest)
		return
	}

	result, err := a.pgDb.Model(&sensor).Insert()
	if err != nil {
		http.Error(w, fmt.Sprintf("couldn't insert model into database: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	fmt.Println(result)
	// return the entire inventory now
	a.getSensors(w, r)

}
func (a *api) deleteSensor(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "provided id was not a number, thus couldn't be parsed", http.StatusBadRequest)
		return
	}

	result, err := a.pgDb.Model(&SensorDevice{Id: int64(id)}).Where("item.id = ?", id).Delete()
	if err != nil {
		http.Error(w, fmt.Sprintf("an error occured when deleteing item with id: %d, error: %s", id, err.Error()), http.StatusInternalServerError)
		return
	}
	fmt.Println(result)
	a.getSensors(w, r)
}

func (a *api) getDevices(w http.ResponseWriter, r *http.Request) {
	var deviceInfos []DeviceInfo
	if err := a.pgDb.Model(&deviceInfos).Select(); err != nil {
		http.Error(w, fmt.Sprintf("couldn't load item from database: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	if deviceInfos == nil {
		deviceInfos = []DeviceInfo{}
	}

	var resp struct {
		Devices []DeviceInfo `json:"devices"`
	}
	resp.Devices = deviceInfos

	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		log.Printf("error sending devices: %s\n", err.Error())
	}
}
func (a *api) getDevice(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "no id provided", http.StatusBadRequest)
		return
	}

	var device DeviceInfo
	if err := a.pgDb.Model(&device).Where("device_info.id = ?", id).Select(); err != nil {
		http.Error(w, fmt.Sprintf("couldn't load item from database: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	if device.ID == "" {
		device = DeviceInfo{}
	}

	var resp struct {
		Device DeviceInfo `json:"device"`
	}
	resp.Device = device

	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		log.Printf("error sending devices: %s\n", err.Error())
	}
}
func (a *api) postDevice(w http.ResponseWriter, r *http.Request) {
	var device DeviceInfo
	if err := json.NewDecoder(r.Body).Decode(&device); err != nil {
		http.Error(w, fmt.Sprintf("couldn't decode body: %s", err.Error()), http.StatusBadRequest)
		return
	}

	result, err := a.pgDb.Model(&device).Insert()
	if err != nil {
		http.Error(w, fmt.Sprintf("couldn't insert model into database: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	fmt.Println(result)
	// return the entire inventory now
	a.getDevices(w, r)
}

func (a *api) deleteDevice(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "no id provided", http.StatusBadRequest)
		return
	}

	result, err := a.pgDb.Model(&DeviceInfo{ID: id}).Where("device_info.id = ?", id).Delete()
	if err != nil {
		http.Error(w, fmt.Sprintf("an error occured when deleteing item with id: %s, error: %s", id, err.Error()), http.StatusInternalServerError)
		return
	}
	fmt.Println(result)
	a.getDevices(w, r)
}

//func (a *api) postDevice(w http.ResponseWriter, r *http.Request) {
//	var dev DeviceInfo
//	if err := json.NewDecoder(r.Body).Decode(&dev); err != nil {
//		log.Println("error decoding the body:", err)
//		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
//		return
//	}
//	a.devices[dev.ID] = dev
//	jsonData, err := json.Marshal(&dev)
//	if err != nil {
//		log.Println(jsonData)
//	}
//	if err := a.db.Set(fmt.Sprintf("registry-%s", dev.ID), jsonData); err != nil {
//		log.Println(err)
//	}
//	a.keys = append(a.keys, fmt.Sprintf("registry-%s", dev.ID))
//	jsonKeys, err := json.Marshal(&a.keys)
//	if err != nil {
//		log.Println(jsonKeys)
//	}
//	if err := a.db.Set("registry-keys", jsonKeys); err != nil {
//		log.Println(err)
//	}
//
//	a.getDevices(w, r)
//}
//
//func (a *api) getDevices(w http.ResponseWriter, r *http.Request) {
//	var devices struct {
//		Devices []DeviceInfo `json:"devices"`
//	}
//	devices.Devices = make([]DeviceInfo, 0)
//	for _, v := range a.devices {
//		devices.Devices = append(devices.Devices, v)
//	}
//
//	w.Header().Set("content-type", "application/json")
//	if err := json.NewEncoder(w).Encode(&devices); err != nil {
//		log.Printf("error sending getDevices: %s\n", err.Error())
//	}
//}
//
//func (a *api) deleteDevice(w http.ResponseWriter, r *http.Request) {
//	id := chi.URLParam(r, "id")
//	if id == "" {
//		http.NotFound(w, r)
//		return
//	}
//	fmt.Println(id)
//	fmt.Println(a.devices)
//	delete(a.devices, id)
//	fmt.Println(a.devices)
//	if err := a.db.Delete(fmt.Sprintf("registry-%s", id)); err != nil {
//		log.Println(err)
//		http.Error(w, fmt.Sprintf("error deleting device: %s", err.Error()), http.StatusInternalServerError)
//	}
//	for i, v := range a.keys {
//		if v == fmt.Sprintf("registry-%s", id) {
//			a.keys = append(a.keys[:i], a.keys[i+1:]...)
//			fmt.Println("removed key from registry-keys")
//			break
//		}
//	}
//	fmt.Println(a.keys)
//
//	jsonKeys, err := json.Marshal(&a.keys)
//	if err != nil {
//		log.Println(jsonKeys)
//	}
//	if err := a.db.Set("registry-keys", jsonKeys); err != nil {
//		log.Println(err)
//	}
//	var msg struct {
//		Message string `json:"message"`
//	}
//	msg.Message = fmt.Sprintf("deleted device: %s successfully", id)
//	w.Header().Set("content-type", "application/json")
//	if err := json.NewEncoder(w).Encode(&msg); err != nil {
//		log.Printf("error sending deleteDevice msg: %s\n", err.Error())
//	}
//}

func (a *api) createSchedule(w http.ResponseWriter, r *http.Request) {
	a.scheduler.CreateSchedule(Schedule{})
}

func (a *api) getSchedules(w http.ResponseWriter, r *http.Request) {

}

func (a *api) getSchedule(w http.ResponseWriter, r *http.Request) {

}

func (a *api) deleteSchedule(w http.ResponseWriter, r *http.Request) {

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
