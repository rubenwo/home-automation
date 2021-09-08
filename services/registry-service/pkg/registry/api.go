package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"github.com/rubenwo/home-automation/services/registry-service/pkg/registry/models"
	"github.com/rubenwo/home-automation/services/registry-service/pkg/registry/routines"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

type api struct {
	router    *chi.Mux
	groups    map[string][]string
	db        *pg.DB
	scheduler *routines.Scheduler
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
		router:    chi.NewRouter(),
		groups:    make(map[string][]string),
		db:        db,
		scheduler: routines.NewScheduler(db, runtime.NumCPU(), cfg.MqttHost, 10),
	}
	go a.scheduler.Run(time.Second * 1)

	// A good base middleware stack
	a.router.Use(middleware.RequestID)
	a.router.Use(middleware.RealIP)
	a.router.Use(middleware.Logger)
	a.router.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	a.router.Use(middleware.Timeout(60 * time.Second))
	//
	//a.router.Use(cors.Handler(cors.Options{
	//	// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
	//	AllowedOrigins:   []string{"*"},
	//	// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
	//	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	//	AllowCredentials: false,
	//}))

	a.router.Get("/healthz", a.healthz)
	a.router.Get("/new_id", a.getNewID)

	a.router.Get("/sensors", a.getSensors)
	a.router.Post("/sensors", a.addSensor)
	a.router.Delete("/sensors/{id}", a.deleteSensor)

	a.router.Get("/devices", a.getDevices)
	a.router.Post("/devices", a.postDevice)
	a.router.Get("/devices/{id}", a.getDevice)
	a.router.Delete("/devices/{id}", a.deleteDevice)

	a.router.Get("/routines", a.getRoutines)
	a.router.Post("/routines", a.createRoutine)
	a.router.Get("/routines/{id}", a.getRoutine)
	a.router.Delete("/routines/{id}", a.deleteRoutine)
	a.router.Put("/routines/{id}", a.updateRoutine)
	a.router.Get("/routines/logs", a.getRoutinesLogs)
	a.router.Get("/routines/logs/{id}", a.getRoutineLogs)

	a.router.Post("/group", a.createGroup)

	return a, nil
}
func (a *api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}
func (a *api) healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(&models.HealthzModel{
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
	var sensors []models.SensorDevice
	if err := a.db.Model(&sensors).Select(); err != nil {
		http.Error(w, fmt.Sprintf("couldn't load item from database: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	if sensors == nil {
		sensors = []models.SensorDevice{}
	}

	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(&sensors); err != nil {
		log.Printf("error sending sensors: %s\n", err.Error())
	}
}

func (a *api) addSensor(w http.ResponseWriter, r *http.Request) {
	var sensor models.SensorDevice
	if err := json.NewDecoder(r.Body).Decode(&sensor); err != nil {
		http.Error(w, fmt.Sprintf("couldn't decode body: %s", err.Error()), http.StatusBadRequest)
		return
	}

	result, err := a.db.Model(&sensor).Insert()
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

	result, err := a.db.Model(&models.SensorDevice{Id: int64(id)}).Where("item.id = ?", id).Delete()
	if err != nil {
		http.Error(w, fmt.Sprintf("an error occured when deleteing item with id: %d, error: %s", id, err.Error()), http.StatusInternalServerError)
		return
	}
	fmt.Println(result)
	a.getSensors(w, r)
}

func (a *api) getDevices(w http.ResponseWriter, r *http.Request) {
	var deviceInfos []models.DeviceInfo
	if err := a.db.Model(&deviceInfos).Select(); err != nil {
		http.Error(w, fmt.Sprintf("couldn't load item from database: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	if deviceInfos == nil {
		deviceInfos = []models.DeviceInfo{}
	}

	var resp struct {
		Devices []models.DeviceInfo `json:"devices"`
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

	var device models.DeviceInfo
	if err := a.db.Model(&device).Where("device_info.id = ?", id).Select(); err != nil {
		http.Error(w, fmt.Sprintf("couldn't load item from database: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	if device.ID == "" {
		device = models.DeviceInfo{}
	}

	var resp struct {
		Device models.DeviceInfo `json:"device"`
	}
	resp.Device = device

	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		log.Printf("error sending devices: %s\n", err.Error())
	}
}
func (a *api) postDevice(w http.ResponseWriter, r *http.Request) {
	var device models.DeviceInfo
	if err := json.NewDecoder(r.Body).Decode(&device); err != nil {
		http.Error(w, fmt.Sprintf("couldn't decode body: %s", err.Error()), http.StatusBadRequest)
		return
	}

	result, err := a.db.Model(&device).Insert()
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

	result, err := a.db.Model(&models.DeviceInfo{ID: id}).Where("device_info.id = ?", id).Delete()
	if err != nil {
		http.Error(w, fmt.Sprintf("an error occured when deleteing item with id: %s, error: %s", id, err.Error()), http.StatusInternalServerError)
		return
	}
	fmt.Println(result)
	a.getDevices(w, r)
}

func (a *api) getRoutines(w http.ResponseWriter, r *http.Request) {
	var routins []models.Routine
	if err := a.db.Model(&routins).Select(); err != nil {
		http.Error(w, fmt.Sprintf("couldn't load item from database: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	if routins == nil {
		routins = []models.Routine{}
	}

	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(&routins); err != nil {
		log.Printf("error sending sensors: %s\n", err.Error())
	}
}

func (a *api) createRoutine(w http.ResponseWriter, r *http.Request) {
	var routine models.Routine
	if err := json.NewDecoder(r.Body).Decode(&routine); err != nil {
		http.Error(w, fmt.Sprintf("couldn't decode body: %s", err.Error()), http.StatusBadRequest)
		return
	}

	result, err := a.db.Model(&routine).Insert()
	if err != nil {
		http.Error(w, fmt.Sprintf("couldn't insert model into database: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	if err := a.scheduler.UpdateRoutines(); err != nil {
		log.Println(err)
	}
	fmt.Println(result)
	// return the entire inventory now
	a.getRoutines(w, r)
}
func (a *api) getRoutine(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "no valid id provided", http.StatusBadRequest)
		return
	}

	var routine models.Routine
	if err := a.db.Model(&routine).Where("routine.id = ?", id).Select(); err != nil {
		http.Error(w, fmt.Sprintf("couldn't load item from database: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	var resp struct {
		Routine models.Routine `json:"routine"`
	}
	resp.Routine = routine

	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		log.Printf("error sending devices: %s\n", err.Error())
	}
}
func (a *api) deleteRoutine(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "provided id was not a number, thus couldn't be parsed", http.StatusBadRequest)
		return
	}

	result, err := a.db.Model(&models.Routine{Id: int64(id)}).Where("routine.id = ?", id).Delete()
	if err != nil {
		http.Error(w, fmt.Sprintf("an error occured when deleteing item with id: %d, error: %s", id, err.Error()), http.StatusInternalServerError)
		return
	}

	result, err = a.db.Model(&models.RoutineLog{RoutineId: int64(id)}).Where("routine_log.routine_id = ?", id).Delete()
	if err != nil {
		http.Error(w, fmt.Sprintf("an error occured when deleteing item with id: %d, error: %s", id, err.Error()), http.StatusInternalServerError)
		return
	}

	if err := a.scheduler.UpdateRoutines(); err != nil {
		log.Println(err)
	}
	fmt.Println(result)
	a.getRoutines(w, r)
}

func (a *api) updateRoutine(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "provided id was not a number, thus couldn't be parsed", http.StatusBadRequest)
		return
	}

	var routine models.Routine
	if err := json.NewDecoder(r.Body).Decode(&routine); err != nil {
		http.Error(w, fmt.Sprintf("couldn't decode body: %s", err.Error()), http.StatusBadRequest)
		return
	}

	if int64(id) != routine.Id {
		http.Error(w, fmt.Sprintf("id in url does not correspond to the provided body"), http.StatusBadRequest)
		return
	}

	result, err := a.db.Model(&routine).Where("routine.id = ?", id).Update(&routine)

	if err != nil {
		http.Error(w, fmt.Sprintf("an error occured when deleteing item with id: %d, error: %s", id, err.Error()), http.StatusInternalServerError)
		return
	}
	if err := a.scheduler.UpdateRoutines(); err != nil {
		log.Println(err)
	}
	fmt.Println(result)
	var resp struct {
		Routine models.Routine `json:"routine"`
	}
	resp.Routine = routine

	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		log.Printf("error sending devices: %s\n", err.Error())
	}
}

func (a *api) getRoutinesLogs(w http.ResponseWriter, r *http.Request) {
	var logs []models.RoutineLog
	if err := a.db.Model(&logs).Select(); err != nil {
		http.Error(w, fmt.Sprintf("couldn't load item from database: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	var resp struct {
		Logs []models.RoutineLog `json:"logs"`
	}
	resp.Logs = logs

	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		log.Printf("error sending devices: %s\n", err.Error())
	}
}

func (a *api) getRoutineLogs(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "provided id was not a number, thus couldn't be parsed", http.StatusBadRequest)
		return
	}
	log.Println(id)

	var logs []models.RoutineLog
	if err := a.db.Model(&logs).Where("routine_log.routine_id = ?", id).Select(); err != nil {
		http.Error(w, fmt.Sprintf("couldn't load item from database: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	var resp struct {
		Logs []models.RoutineLog `json:"logs"`
	}
	resp.Logs = logs

	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		log.Printf("error sending devices: %s\n", err.Error())
	}
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
