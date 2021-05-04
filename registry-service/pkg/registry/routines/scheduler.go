package routines

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/rubenwo/home-automation/registry-service/pkg/registry/models"
	"io/ioutil"
	"log"
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
			if checkIfRoutineShouldRun(routine.Trigger, currentTime, 1e+10) {
				for _, action := range routine.Actions {
					s.jobs <- action
				}
			}
		}
		s.Unlock()
	}
}

func checkIfTimerRoutineShouldRun(schedule models.Schedule, month, day, hour, minute int) bool {
	if schedule.DayOfWeek != models.Star {
		if schedule.DayOfWeek.Int() == day {

		}
	}

	return false
}

func checkIfRoutineShouldRun(trigger models.Trigger, currentTime time.Time, diffTimeInNanoS float64) bool {
	switch trigger.Type {
	case models.TimerTriggerType:
		//nextTime := cronexpr.MustParse("0 0 29 2 *").Next(time.Now())
		//return  math.Abs(float64(time.Now().UnixNano() - nextTime.UnixNano())) < diffTimeInNanoS
		//

		_, month, day := currentTime.Date()
		hour, minute, _ := currentTime.Clock()
		return checkIfTimerRoutineShouldRun(trigger.Schedule, int(month), day, hour, minute)
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
