package handler

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gustavobelfort/42-jitsi/internal/config"
)

// customTime will allow us to support unmarshalling time that is not RFC 3339
type customTime struct {
	time.Time
}

var timeLayout string

func (ct *customTime) UnmarshalJSON(d []byte) error {
	if ct.Time.UnmarshalJSON(d) == nil {
		return nil
	}
	if timeLayout == "" {
		timeLayout = config.Conf.BeginAtTimeLayout
	}
	var err error
	ct.Time, err = time.Parse(fmt.Sprintf(`"%s"`, timeLayout), string(d))
	return err
}

type scaleTeam struct {
	ID         int        `json:"id"`
	BeginAt    customTime `json:"begin_at"`
	Corrector  string
	Correcteds []string
	TeamID     int
}

func (st *scaleTeam) validate() error {
	missing := make([]string, 0)
	if st.ID == 0 {
		missing = append(missing, "id")
	}
	if st.BeginAt.IsZero() {
		missing = append(missing, "begin_at")
	}
	if st.TeamID == 0 {
		missing = append(missing, "team.id")
	}

	if len(missing) != 0 {
		return &MissingFieldsError{missing: missing}
	}

	if st.Corrector == "" {
		return NoCorrectorError
	}

	return nil
}

// UnmarshalJSON will unmarshal the evaluation's payload into the scaleTeam structure.
func (st *scaleTeam) UnmarshalJSON(d []byte) error {
	type scaleTeamTwin scaleTeam
	type scaleTeamUnmarshaller struct {
		scaleTeamTwin
		User struct {
			Login string `json:"login"`
		} `json:"user"`
		Team struct {
			ID int `json:"id"`
		} `json:"team"`
	}
	unmarshaller := &scaleTeamUnmarshaller{}
	if err := json.Unmarshal(d, unmarshaller); err != nil {
		return err
	}

	*st = scaleTeam(unmarshaller.scaleTeamTwin)
	st.Corrector = unmarshaller.User.Login
	st.TeamID = unmarshaller.Team.ID

	return st.validate()
}
