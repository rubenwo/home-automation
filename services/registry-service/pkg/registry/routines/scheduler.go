package routines

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/gorhill/cronexpr"
	"github.com/rubenwo/home-automation/registry-service/pkg/registry/models"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"sync"
	"time"
)

//Scheduler holds the routines and their triggers. When a trigger event is raised
type Scheduler struct {
	sync.Mutex
	routines []models.Routine
	jobs     chan models.Action
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
		jobs:    make(chan models.Action, 100),
		results: make(chan error, 100),
		db:      db}
	for i := 0; i < maxConcurrentWorkers; i++ {
		go s.worker()
	}
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

			if checkIfRoutineShouldRun(routine.Trigger, currentTime, interval.Nanoseconds()) {
				for _, action := range routine.Actions {
					s.jobs <- action
				}
			}
		}
		s.Unlock()
	}
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
	for err := range s.results {
		log.Println(err)
	}
}

func (s *Scheduler) worker() {
	for action := range s.jobs {

		fmt.Println(action)
		client := &http.Client{}
		var (
			req *http.Request
			err error
		)

		//vm := otto.New()
		//_, err = vm.Run(action.Script)
		//if err != nil {
		//	s.results <- err
		//	continue
		//}

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
		raw, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			s.results <- err
			continue
		}
		if err := resp.Body.Close(); err != nil {
			s.results <- err
			continue
		}
		fmt.Println(string(raw))
	}
}
