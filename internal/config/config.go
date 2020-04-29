package config

import (
	"errors"
	"reflect"
	"strings"
	"time"
)

// Configuration is the type that will hold the configuration informations
type Configuration struct {
	EmailSuffix string          `mapstructure:"email_suffix"`
	SlackThat   SlackThatConfig `mapstructure:"slack_that"`
	WarnBefore  time.Duration   `mapstructure:"warn_before"`
	Intra       Intra
	Postgres    Database
}

// Database is the type that will hold the database informations
type Database struct {
	Host     string
	Port     string
	DB       string
	User     string
	Password string
}

// Intra is the type that will hold the Intranet configurations
type Intra struct {
	AppID     string `mapstructure:"app_id"`
	AppSecret string `mapstructure:"app_secret"`
	Webhooks  map[string]string
}

// Configurations for the Slackthat Microsservice
type SlackThatConfig struct {
	URL       string
	Workspace string
	Username  string
}

// stringToMapstringHookFunc will decode a string to a mapstring.
func stringToMapstringHookFunc(f, t reflect.Type, data interface{}) (interface{}, error) {
	if f.Kind() != reflect.String || t != reflect.MapOf(reflect.TypeOf(""), reflect.TypeOf("")) {
		return data, nil
	}

	mapstring := make(map[string]string)
	for _, elements := range strings.Split(data.(string), ",") {
		config := strings.Split(elements, ":")
		if len(config) != 2 {
			return nil, errors.New("expected string of format 'key0:value0,key1:value1,...,keyN:valueN'")
		}
		mapstring[config[0]] = config[1]
	}

	return mapstring, nil
}
