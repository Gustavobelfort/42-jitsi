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

type SchedulerMock struct {
	mock.Mock
}

type TestSchedulerSuite struct {
	suite.Suite

	expected []int
	tasks    []scheduler.Task
	mock     *SchedulerMock
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
	s.mock = &SchedulerMock{}

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
	expected := []int{1, 2}
	time.Sleep(1 * time.Second)
	s.Require().Equal(expected, s.expected)
}

func (s *TestSchedulerSuite) task1() {
	s.expected = append(s.expected, 1)
}

func (s *TestSchedulerSuite) task2() {
	s.expected = append(s.expected, 2)
}
