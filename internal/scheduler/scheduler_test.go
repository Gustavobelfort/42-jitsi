package scheduler_test

import (
	"testing"
	"time"

	"github.com/gustavobelfort/42-jitsi/internal/scheduler"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestScheduler(t *testing.T) {
	suite.Run(t, new(TestSchedulerSuite))
}

type TestSchedulerSuite struct {
	suite.Suite

	tasks []scheduler.Task
	mock  *mock.Mock
	c     chan struct{}
}

func (s *TestSchedulerSuite) SetupSuite() {

	s.tasks = []scheduler.Task{
		{
			Task:     s.task1,
			Interval: time.Second * 1,
		},
		{
			Task:     s.task2,
			Interval: time.Second * 1,
		},
	}
	s.mock = &mock.Mock{}
	s.c = make(chan struct{})

}

func (s *TestSchedulerSuite) SetupTest() {

	s.mock.Calls = []mock.Call{}
	s.mock.ExpectedCalls = []*mock.Call{}
}

func (s *TestSchedulerSuite) Test00_NewScheduler() {
	scheduler, err := scheduler.New(s.tasks)
	s.Require().NoError(err)
	s.Require().NotNil(scheduler)
}

func (s *TestSchedulerSuite) Test01_ScheduledTasks() {
	scheduler, _ := scheduler.New(s.tasks)
	scheduler.Start()
	defer scheduler.Stop()
	s.mock.On("task1").Return().Once()
	s.mock.On("task2").Return().Once()

	ticker := time.NewTicker(time.Second * 2)
	defer ticker.Stop()
	for i := 0; i < 2; {
		select {
		case <-ticker.C:
			i = 2
		case <-s.c:
			i++
		}
	}
}

func (s *TestSchedulerSuite) TearDownTest() {
	s.mock.AssertExpectations(s.T())
}

func (s *TestSchedulerSuite) task1() {
	s.mock.MethodCalled("task1")
	s.c <- struct{}{}
}

func (s *TestSchedulerSuite) task2() {
	s.mock.MethodCalled("task2")
	s.c <- struct{}{}
}
