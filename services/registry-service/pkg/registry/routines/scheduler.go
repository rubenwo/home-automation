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

type job struct {
	RoutineId int64
	Action    models.Action
}

//Scheduler holds the routines and their triggers. When a trigger event is raised
type Scheduler struct {
	sync.Mutex
	routines []models.Routine
	jobs     chan job
	results  chan error
	db       *pg.DB
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

func NewScheduler(db *pg.DB, maxConcurrentWorkers int) *Scheduler {
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
		s.Lock()
		for _, routine := range s.routines {
			if !routine.IsActive {
				continue
			}
			if checkIfRoutineShouldRun(routine.Trigger, currentTime, interval.Nanoseconds()) {
				for _, action := range routine.Actions {
					s.jobs <- job{
						RoutineId: routine.Id,
						Action:    action,
					}
				}
			}
		}
		s.Unlock()
	}
}

func (s *Scheduler) eventWorker(path, host string, retry int) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("%s:1883", host))
	opts.SetClientID(uuid.New().String())
	client := mqtt.NewClient(opts)

	var err error
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
		return
	}

	client.Subscribe(path, 0, func(cl mqtt.Client, msg mqtt.Message) {
		msg.Ack()
		payload := msg.Payload()
		fmt.Println(string(payload))
	})

	s.Lock()
	for _, routine := range s.routines {
		if !routine.IsActive {
			continue
		}
		if routine.Trigger.Type != models.MqttEventTriggerType {
			continue
		}

	}
	s.Unlock()

	//TODO: Push job to worker
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
		log.Printf("Scheduler()->resultWorker() received an error in channel: %s\n", err.Error())
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
