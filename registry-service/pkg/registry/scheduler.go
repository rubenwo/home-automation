package registry

import (
	"fmt"
	"github.com/jasonlvhit/gocron"
	"github.com/rubenwo/home-automation/registry-service/pkg/rlang"
)

type Schedule struct {
	Crontab string `json:"crontab"`
	Every   int    `json:"every"`
	At      string `json:"at"`
	Script  string `json:"script"`
}

type Scheduler struct {
}

func NewScheduler() *Scheduler { return &Scheduler{} }

func (s *Scheduler) CreateSchedule(schedule Schedule) error {
	fn, err := rlang.Parse(schedule.Script)
	if err != nil {
		return fmt.Errorf("couldn't parse script: %w", err)
	}
	err = gocron.Every(uint64(schedule.Every)).Day().At(schedule.At).Do(fn)
	if err != nil {
		return fmt.Errorf("error scheduling job: %w", err)
	}
	go gocron.Every(1).Second().Do(fn)
	return nil
}
