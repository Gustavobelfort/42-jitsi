package scheduler

import (
	"fmt"
	"reflect"
	"runtime"
	"time"

	"github.com/robfig/cron/v3"
)

func New(tasks []Task) (Scheduler, error) {
	s := Scheduler{}

	s.Scheduler = cron.New()

	if tasks == nil {
		errorMessage := fmt.Errorf("empty task list")
		return s, errorMessage
	}

	for _, t := range tasks {
		err := s.Add(t.Task, t.Interval)
		if err != nil {
			return s, err
		}
	}

	return s, nil
}

func (s *Scheduler) Add(task func(), every time.Duration) error {

	formatedDuration := FormatDuration(every)

	_, err := s.Scheduler.AddFunc(formatedDuration, task)
	if err != nil {
		errorMessage := fmt.Errorf("failed to add %s to the scheduler, err: %v", getFunctionName(task), err)
		return errorMessage
	}

	return nil

}

func (s *Scheduler) Start() {
	s.Scheduler.Start()
}

func FormatDuration(interval time.Duration) string {
	return fmt.Sprintf("@every %s", interval.String())
}

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
