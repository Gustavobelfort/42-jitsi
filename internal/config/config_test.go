package config

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestConfigLoad(t *testing.T) {

	expected := Configuration{
		SlackThatURL: "localhost:8080",
		CampusSlug:   "42sp",
		WarnBefore:   time.Minute * 15,
		IntraWebhooksAuth: []Webhook{
			Webhook{
				Hook:   "scale_team.create",
				Secret: "create_secret",
			},
			Webhook{
				Hook:   "scale_team.update",
				Secret: "update_secret",
			},
			Webhook{
				Hook:   "scale_team.destroy",
				Secret: "destroy_secret",
			},
		},
		Postgres: Database{
			Host:     "1",
			Port:     "2",
			Database: "3",
			Username: "4",
			Password: "5",
		},
	}

	got, err := load("../../configs/configs.sample.yml")

	if err != nil {
		t.Errorf("load failed with error: %v", err)
	}
	if !cmp.Equal(got, expected) {
		t.Errorf("Got: %#v, Expected: %#v", got, expected)
	}
}

func TestConfigCheck(t *testing.T) {
	err := check()
	if err == nil {
		t.Errorf("Got: nil, Expected: 'no such file or directory'")
	}
}
