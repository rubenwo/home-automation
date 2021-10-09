package routines

import (
	"bytes"
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"github.com/gorhill/cronexpr"
	"github.com/robertkrimen/otto"
	"github.com/rubenwo/home-automation/services/registry-service/pkg/registry/models"
	"github.com/rubenwo/home-automation/services/registry-service/pkg/registry/routines/script"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"sync"
	"time"
)

const (
	NotificationsMqttPath = "notifications"
	EventsMqttPath        = "events"

	QosAtMostOnce  = 0
	QosAtLeastOnce = 1
	QosExactlyOnce = 2
)

type job struct {
	RoutineId int64
	Action    models.Action
}

//Scheduler holds the routines and their triggers. When a trigger event is raised
type Scheduler struct {
	sync.RWMutex
	routines []models.Routine
	jobs     chan job
	results  chan error
	db       *pg.DB

	mqttClient mqtt.Client
}

func initRoutines(db *pg.DB) ([]models.Routine, error) {
	var routines []models.Routine
	if err := db.Model(&routines).Select(); err != nil {
		return nil, err
	}

	if routines == nil {
		routines = []models.Routine{}
	}

	return routines, nil
}

func NewScheduler(db *pg.DB, maxConcurrentWorkers int, host string, retry int) *Scheduler {
	routines, err := initRoutines(db)
	if err != nil {
		log.Fatal(err)
	}

	s := &Scheduler{routines: routines,
		jobs:    make(chan job, 100),
		results: make(chan error, 100),
		db:      db}
	for i := 0; i < maxConcurrentWorkers; i++ {
		go s.worker()
	}
	go s.resultWorker()

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("%s:1883", host))
	opts.SetClientID(uuid.New().String())
	client := mqtt.NewClient(opts)

	for i := 0; i < retry; i++ {
		token := client.Connect()
		token.Wait()
		err = token.Error()
		if err == nil {
			break
		}
		time.Sleep(time.Second * 1)
	}
	if err != nil {
		log.Fatalf("could not connect to mqtt broker: %s", err.Error())
	}

	s.mqttClient = client
	go s.eventWorker()

	return s
}

func (s *Scheduler) UpdateRoutines() error {
	routines, err := initRoutines(s.db)
	if err != nil {
		return err
	}
	s.Lock()
	s.routines = routines
	s.Unlock()
	return nil
}

//Run will periodically check at the specified interval rate if actions should run.
//Let's say interval is 10 seconds, this means that every 10 seconds we'll check if an action should run.
func (s *Scheduler) Run(interval time.Duration) {
	for range time.Tick(interval) {
		currentTime := time.Now()
		s.RLock()
		for _, routine := range s.routines {
			if !routine.IsActive {
				continue
			}
			if checkIfRoutineShouldRun(routine.Trigger, currentTime, interval.Nanoseconds()) {
				data, _ := json.Marshal(models.Notification{
					Title: fmt.Sprintf("Running routine: %s\n", routine.Name),
					Body:  fmt.Sprintf("Starting routine: %s with %d actions\n", routine.Name, len(routine.Actions)),
				})
				token := s.mqttClient.Publish(NotificationsMqttPath, 0, false, data)
				token.Wait()
				for _, action := range routine.Actions {
					s.jobs <- job{
						RoutineId: routine.Id,
						Action:    action,
					}
				}
			}
		}
		s.RUnlock()
	}
}

func (s *Scheduler) eventWorker() {
	s.mqttClient.Subscribe(NotificationsMqttPath, 0, func(cl mqtt.Client, msg mqtt.Message) {
		msg.Ack()
		payload := msg.Payload()
		fmt.Println(string(payload))
	})

	type eventMsg struct {
		Name string `json:"name"`
	}

	s.mqttClient.Subscribe(EventsMqttPath, 0, func(cl mqtt.Client, msg mqtt.Message) {
		msg.Ack()
		payload := msg.Payload()
		fmt.Println(string(payload))
		var eventMsg eventMsg
		if err := json.Unmarshal(payload, &msg); err != nil {
			s.results <- err
			return
		}
		s.RLock()
		for _, routine := range s.routines {
			if routine.IsActive && routine.Trigger.Type == models.MqttEventTriggerType && routine.Trigger.OnEvent == eventMsg.Name {
				if err := s.Trigger(routine.Id); err != nil {
					s.results <- err
				}
			}
		}
		s.RUnlock()
	})
}

func checkIfRoutineShouldRun(trigger models.Trigger, currentTime time.Time, diffTimeInNanoS int64) bool {
	switch trigger.Type {
	case models.TimerTriggerType:
		nextTime := cronexpr.MustParse(trigger.CronExpr).Next(time.Now())
		return math.Abs(float64(currentTime.UnixNano()-nextTime.UnixNano())) < float64(diffTimeInNanoS)

	default:
		return false
	}
}

func (s *Scheduler) resultWorker() {
	log.Println("Scheduler()->resultWorker() started")
	for err := range s.results {
		if err == nil {
			continue
		}
		log.Printf("Scheduler()->resultWorker() received an error in channel: %s\n", err.Error())
		data, _ := json.Marshal(models.Notification{
			Title: "Scheduler encountered an error",
			Body:  fmt.Sprintf("error: %s\n", err.Error()),
		})
		token := s.mqttClient.Publish(NotificationsMqttPath, 0, false, data)
		token.Wait()
	}
	log.Println("Scheduler()->resultWorker() finished")
}

func (s *Scheduler) worker() {
	log.Println("Scheduler()->worker() started")
	for job := range s.jobs {
		action := job.Action
		if action.Script != "" {
			vm := otto.New()

			if err := vm.Set("HttpGet", script.HttpGet); err != nil {
				s.results <- err
				continue
			}
			if err := vm.Set("HttpPost", script.HttpPost); err != nil {
				s.results <- err
				continue
			}
			if err := vm.Set("HttpDelete", script.HttpDelete); err != nil {
				s.results <- err
				continue
			}
			if err := vm.Set("HttpPut", script.HttpPut); err != nil {
				s.results <- err
				continue
			}
			if err := vm.Set("__log__", script.Log(job.RoutineId, s.db)); err != nil {
				s.results <- err
				continue
			}

			sc := fmt.Sprintf("console.log = __log__;\n%s", action.Script)

			result, err := vm.Run(sc)
			if err != nil {
				s.results <- err
				continue
			}
			_, _ = s.db.Model(&models.RoutineLog{
				RoutineId: job.RoutineId,
				LoggedAt:  time.Now(),
				Message:   result.String(),
			}).Insert()

			log.Printf("Scheduler()->worker() finished executing script: [%s]\n", action.Script)
		}

		if action.Method == "" {
			continue
		}
		if action.Addr == "" {
			continue
		}

		client := &http.Client{}
		var (
			req *http.Request
			err error
		)

		if action.Data == nil {
			req, err = http.NewRequest(action.Method, action.Addr, nil)
			if err != nil {
				s.results <- err
				continue
			}
		} else {
			b, err := json.Marshal(action.Data)
			if err != nil {
				s.results <- err
				continue
			}
			req, err = http.NewRequest(action.Method, action.Addr, bytes.NewBuffer(b))
			if err != nil {
				s.results <- err
				continue
			}
		}
		resp, err := client.Do(req)
		if err != nil {
			s.results <- err
			continue
		}
		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			s.results <- err
			continue
		}
		if err := resp.Body.Close(); err != nil {
			s.results <- err
			continue
		}

		log.Printf("Scheduler()->worker() finished executing request: [%s] - [%s]\n", action.Method, action.Addr)
		_, _ = s.db.Model(&models.RoutineLog{
			RoutineId: job.RoutineId,
			LoggedAt:  time.Now(),
			Message:   fmt.Sprintf("Scheduler()->worker() finished executing request: [%s] - [%s]\n", action.Method, action.Addr),
		}).Insert()
	}

	log.Println("Scheduler()->worker() finished")
}

func (s *Scheduler) Trigger(id int64) error {
	var routine *models.Routine
	s.RLock()
	for _, m := range s.routines {
		if m.Id == id {
			routine = &m
			break
		}
	}
	s.RUnlock()

	if routine == nil {
		return ErrRoutineNotFound
	}

	for _, action := range routine.Actions {
		s.jobs <- job{
			RoutineId: routine.Id,
			Action:    action,
		}
	}

	return nil
}
