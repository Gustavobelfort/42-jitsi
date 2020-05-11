package scheduler

import (
	"fmt"
	"reflect"
	"runtime"
	"time"

	"github.com/robfig/cron/v3"
)

// New returns a new cron job runner with the given tasks
func New(tasks []Task) (SchedulerInterface, error) {
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

// Add adds a func to the Cron to be run on the given schedule.
func (s Scheduler) Add(task func(), every time.Duration) error {

	formatedDuration := formatDuration(every)

	_, err := s.Scheduler.AddFunc(formatedDuration, task)
	if err != nil {
		errorMessage := fmt.Errorf("failed to add %s to the scheduler, err: %v", getFunctionName(task), err)
		return errorMessage
	}

	return nil

}

// Start the cron scheduler in its own goroutine, or no-op if already started.
func (s Scheduler) Start() error {
	s.Scheduler.Start()
	return nil
}

// Stop stops the cron scheduler if it is running; otherwise it does nothing. A context is returned so the caller can wait for running jobs to complete.
func (s Scheduler) Stop() error {
	s.Scheduler.Stop()
	return nil
}

func formatDuration(interval time.Duration) string {
	return fmt.Sprintf("@every %s", interval.String())
}

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
