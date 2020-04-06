package config

import (
	"testing"
)

func TestConfigLoad(t *testing.T) {

	expected := Configuration{
		SlackThatURL:      "localhost:8080",
		CampusSlug:        "42sp",
		WarnBefore:        "15m",
		IntraWebhooksAuth: "scale_team.create:create_secret,scale_team.update:update_secret,scale_team.destroy:destroy_secret",
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

	if got != expected {
		t.Errorf("Got: %v, Expected: %v", got, expected)
	}
}

func TestConfigCheck(t *testing.T) {
	err := check()
	if err == nil {
		t.Errorf("Got: nil, Expected: 'no such file or directory'")
	}
}
