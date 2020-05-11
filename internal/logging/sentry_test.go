package logging

import (
	"context"
	"errors"
	"io/ioutil"
	"testing"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/magiconair/properties/assert"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type HubMock struct {
	mock.Mock
}

func (m *HubMock) AddBreadcrumb(breadcrumb *sentry.Breadcrumb, _ *sentry.BreadcrumbHint) {
	m.Called(breadcrumb.Level, breadcrumb.Message, logrus.Fields(breadcrumb.Data), breadcrumb.Category, breadcrumb.Type)
}

func (m *HubMock) CaptureException(exception error) *sentry.EventID {
	return m.Called(exception).Get(0).(*sentry.EventID)
}

func (m *HubMock) Flush(timeout time.Duration) bool {
	return m.Called(timeout).Bool(0)
}

func TestSentryHook(t *testing.T) {
	suite.Run(t, new(SentryHookSuite))
}

type SentryHookSuite struct {
	suite.Suite

	mock *HubMock

	logger   *logrus.Logger
	testHook *test.Hook

	sentryHook *SentryHook
}

func (s *SentryHookSuite) SetupSuite() {
	s.mock = &HubMock{}

	s.logger, s.testHook = test.NewNullLogger()

	s.sentryHook = &SentryHook{hub: s.mock}
	s.Equal(logrus.AllLevels, s.sentryHook.Levels())
	s.logger.AddHook(s.sentryHook)
	s.logger.SetLevel(logrus.DebugLevel)
}

func (s *SentryHookSuite) SetupTest() {
	s.mock.Calls = []mock.Call{}
	s.mock.ExpectedCalls = []*mock.Call{}

	s.testHook.Reset()
}

func (s *SentryHookSuite) Test00_addErrorLevels() {
	expected := map[logrus.Level]interface{}{
		logrus.ErrorLevel: struct{}{},
		logrus.FatalLevel: struct{}{},
	}
	s.sentryHook.addErrorLevels([]logrus.Level{logrus.ErrorLevel, logrus.FatalLevel})
	s.Equal(expected, s.sentryHook.errorLevels)
}

func (s *SentryHookSuite) Test01_InfoLevel() {
	expectedLevel := sentry.LevelInfo
	expectedMessage := "testing message"
	expectedFields := logrus.Fields{"testing": "data"}
	expectedCategory := "default"
	expectedType := "default"

	s.mock.On("AddBreadcrumb", expectedLevel, expectedMessage, expectedFields, expectedCategory, expectedType).Return().Once()

	s.logger.WithFields(expectedFields).Info(expectedMessage)
}

func (s *SentryHookSuite) Test02_DebugLevel_ChangeCategory() {
	expectedLevel := sentry.LevelDebug
	expectedMessage := "testing message"
	expectedFields := logrus.Fields{"testing": "data"}
	expectedCategory := "testing"
	expectedType := "default"

	s.mock.On("AddBreadcrumb", expectedLevel, expectedMessage, expectedFields, expectedCategory, expectedType).Return().Once()

	ctx := ContextWithSentryCategory(context.Background(), expectedCategory)
	s.logger.WithContext(ctx).WithFields(expectedFields).Debug(expectedMessage)
}

func (s *SentryHookSuite) Test03_ErrorLevel() {
	expectedError := errors.New("testing error")
	expectedLevel := sentry.LevelError
	expectedMessage := "testing message"
	expectedFields := logrus.Fields{"testing": "data", "error": expectedError}
	expectedCategory := "default"
	expectedType := "error"

	expectedID := "id"

	s.mock.On("CaptureException", expectedError).Return((*sentry.EventID)(&expectedID))
	s.mock.On("AddBreadcrumb", expectedLevel, expectedMessage, expectedFields, expectedCategory, expectedType).Return().Once()

	s.logger.WithFields(expectedFields).WithError(expectedError).Error(expectedMessage)
}

func (s *SentryHookSuite) Test04_ErrorLevel_StringError() {
	expectedError := errors.New("testing error")
	expectedLevel := sentry.LevelError
	expectedMessage := "testing message"
	expectedFields := logrus.Fields{"testing": "data", "error": expectedError.Error()}
	expectedCategory := "default"
	expectedType := "error"

	expectedID := "id"

	s.mock.On("CaptureException", expectedError).Return((*sentry.EventID)(&expectedID))
	s.mock.On("AddBreadcrumb", expectedLevel, expectedMessage, expectedFields, expectedCategory, expectedType).Return().Once()

	s.logger.WithFields(expectedFields).Error(expectedMessage)
}

type stringer struct{}

func (*stringer) String() string {
	return "testing error"
}

func (s *SentryHookSuite) Test05_ErrorLevel_StringerError() {
	expectedError := errors.New("testing error")
	expectedLevel := sentry.LevelError
	expectedMessage := "testing message"
	expectedFields := logrus.Fields{"testing": "data", "error": &stringer{}}
	expectedCategory := "default"
	expectedType := "error"

	expectedID := "id"

	s.mock.On("CaptureException", expectedError).Return((*sentry.EventID)(&expectedID))
	s.mock.On("AddBreadcrumb", expectedLevel, expectedMessage, expectedFields, expectedCategory, expectedType).Return().Once()

	s.logger.WithFields(expectedFields).Error(expectedMessage)
}

func (s *SentryHookSuite) Test06_ErrorLevel_NoError() {
	expectedLevel := sentry.LevelError
	expectedMessage := "testing message"
	expectedFields := logrus.Fields{"testing": "data"}
	expectedCategory := "default"
	expectedType := "error"
	expectedError := errors.New(expectedMessage)

	expectedID := "id"

	s.mock.On("CaptureException", expectedError).Return((*sentry.EventID)(&expectedID))
	s.mock.On("AddBreadcrumb", expectedLevel, expectedMessage, expectedFields, expectedCategory, expectedType).Return().Once()

	s.logger.WithFields(expectedFields).Error(expectedMessage)
}

func (s *SentryHookSuite) TearDownTest() {
	s.mock.AssertExpectations(s.T())
}

func TestAddSentryHook(t *testing.T) {
	logger := logrus.New()
	logger.SetOutput(ioutil.Discard)

	expectedHub := &sentry.Hub{}
	expectedErrorLevels := map[logrus.Level]interface{}{
		logrus.ErrorLevel: struct{}{},
		logrus.FatalLevel: struct{}{},
	}

	AddSentryHook(expectedHub, logger, []logrus.Level{logrus.ErrorLevel, logrus.FatalLevel})

	for _, hooks := range logger.Hooks {
		require.IsType(t, &SentryHook{}, hooks[0])
	}
	hook := logger.Hooks[logrus.ErrorLevel][0].(*SentryHook)

	assert.Equal(t, expectedHub, hook.hub)
	assert.Equal(t, expectedErrorLevels, hook.errorLevels)
}
