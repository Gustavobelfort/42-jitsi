package scheduler

import (
	"time"

	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	Scheduler *cron.Cron
}

type Task struct {
	Task     func()
	Interval time.Duration
}

type SchedulerInterface interface {
	Add(task func(), every time.Duration) error
	Start()
}
