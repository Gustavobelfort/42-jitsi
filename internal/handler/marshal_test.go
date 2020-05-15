package handler

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCustomTime_UnmarshalJSON(t *testing.T) {
	timeLayout = "2006-01-02 15:04:05 UTC" // As configuration is not loaded

	t.Run("RFC_3339", func(t *testing.T) {
		timeString := "2020-07-15T21:00:00.000Z"
		expectedTime, _ := time.Parse(time.RFC3339, timeString)
		ct := &customTime{}
		assert.NoError(t, ct.UnmarshalJSON([]byte(fmt.Sprintf(`"%s"`, timeString))))
		assert.True(t, expectedTime.Equal(ct.Time))
	})

	t.Run("No_RFC_3339", func(t *testing.T) {
		timeString := "2020-05-14 22:30:00 UTC"
		expectedTime, _ := time.Parse(timeLayout, timeString)
		ct := &customTime{}
		assert.NoError(t, ct.UnmarshalJSON([]byte(fmt.Sprintf(`"%s"`, timeString))))
		assert.True(t, expectedTime.Equal(ct.Time))
	})
}

func TestScaleTeamMarshal(t *testing.T) {
	t.Run("ValidPayload", func(t *testing.T) {
		expected := &scaleTeam{
			ID:        21,
			BeginAt:   customTime{time.Now()},
			Corrector: "xlogin",
			TeamID:    42,
		}

		payload := []byte(fmt.Sprintf(`{
	"id": %d,
	"begin_at": "%s",
	"user": {"login": "%s"},
	"team": {"id": %d}
}`, expected.ID, expected.BeginAt.Format(time.RFC3339), expected.Corrector, expected.TeamID))

		st := &scaleTeam{}
		assert.NoError(t, json.Unmarshal(payload, &st))
		assert.Equal(t, expected.BeginAt.Format(time.RFC3339), st.BeginAt.Format(time.RFC3339))
		st.BeginAt = expected.BeginAt
		assert.Equal(t, expected, st)
	})

	t.Run("InvalidPayload", func(t *testing.T) {
		payload := []byte(`{"id": "bad"}`)

		st := scaleTeam{}
		err := json.Unmarshal(payload, &st)
		assert.Error(t, err)
		assert.IsType(t, &json.UnmarshalTypeError{}, err)
	})

	t.Run("IncompletePayload", func(t *testing.T) {
		st := &scaleTeam{}
		err := json.Unmarshal([]byte(`{}`), &st)
		assert.Error(t, err)
		assert.IsType(t, &MissingFieldsError{}, err)
		assert.Equal(t, "missing required fields: id,begin_at,team.id", err.Error())
	})

	t.Run("MissingCorrector", func(t *testing.T) {
		st := &scaleTeam{}
		err := json.Unmarshal([]byte(`{
	"id": 21,
	"begin_at": "2020-07-15T21:00:00.000Z",
	"team": {"id": 42}
}`), &st)
		assert.Error(t, err)
		assert.Equal(t, NoCorrectorError, err)
	})
}
