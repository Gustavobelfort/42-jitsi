package config

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequiredConfigError(t *testing.T) {
	t.Run("ConfigSet", func(t *testing.T) {
		viper.Set("expected_field", "test")
		assert.NoError(t, checkRequired("expected_field"))

	})

	t.Run("ConfigNotSet", func(t *testing.T) {
		err := checkRequired("postgres.password")
		require.Error(t, err)
		assert.Equal(t, "missing required configuration: postgres.password", err.Error())
	})

}

func TestStringToMapstringFuncHook(t *testing.T) {

	t.Run("BadFormat", func(t *testing.T) {
		str := "badformat"
		mapstring, err := stringToMapstringHookFunc(reflect.TypeOf(str), reflect.MapOf(reflect.TypeOf(str), reflect.TypeOf(str)), str)
		assert.Nil(t, mapstring)
		require.Error(t, err)
		assert.Equal(t, "expected string of format 'key0:value0,key1:value1,...,keyN:valueN'", err.Error())
	})

	t.Run("BadFormat", func(t *testing.T) {
		str := "key:value,key2:value2"
		expected := map[string]string{"key": "value", "key2": "value2"}
		mapstring, err := stringToMapstringHookFunc(reflect.TypeOf(str), reflect.MapOf(reflect.TypeOf(str), reflect.TypeOf(str)), str)
		require.NoError(t, err)
		assert.Equal(t, expected, mapstring)
	})

}

func TestInitiate(t *testing.T) {
	t.Run("NoConfigFile", func(t *testing.T) {
		os.Unsetenv("CONFIG_FILE")

		// Setting required fields
		os.Setenv("INTRA_APP_ID", "intra_app_id")
		os.Setenv("INTRA_APP_SECRET", "intra_app_secret")
		os.Setenv("INTRA_WEBHOOKS", "key:value")
		os.Setenv("SLACK_THAT_WORKSPACE", "42born2code")
		os.Setenv("POSTGRES_PASSWORD", "changeme")

		defer os.Unsetenv("INTRA_APP_ID")
		defer os.Unsetenv("INTRA_APP_SECRET")
		defer os.Unsetenv("INTRA_WEBHOOKS")
		defer os.Unsetenv("SLACK_THAT_WORKSPACE")
		defer os.Unsetenv("POSTGRES_PASSWORD")

		expected := Configuration{
			SlackThat: SlackThatConfig{
				URL:       "http://localhost:8080",
				Workspace: "42born2code",
				Username:  "Evaluation Master",
			},
			EmailSuffix: "student.42campus.org",
			WarnBefore:  time.Minute * 15,
			HTTPAddr:    "0.0.0.0:5000",
			Timeout:     time.Second * 10,
			Intra: Intra{
				AppID:     "intra_app_id",
				AppSecret: "intra_app_secret",
				Webhooks: map[string]string{
					"key": "value",
				},
			},
			Postgres: Database{
				Host:     "localhost",
				Port:     "5432",
				DB:       "postgres",
				User:     "postgres",
				Password: "changeme",
			},
			RabbitMQ: RabbitMQ{
				Host:     "localhost",
				Port:     "5672",
				VHost:    "",
				User:     "guest",
				Password: "guest",
				Queue:    "webhooks_intra_42jitsi",
			},
		}

		assert.NoError(t, Initiate())
		assert.Equal(t, expected, Conf)
	})

	t.Run("MissingRequired", func(t *testing.T) {
		os.Unsetenv("CONFIG_FILE")

		assert.Error(t, Initiate())
	})

	t.Run("FullConfigFile", func(t *testing.T) {
		os.Setenv("CONFIG_FILE", "../../configs/configs.sample.yml")

		expected := Configuration{
			SlackThat: SlackThatConfig{
				URL:       "http://localhost:8080",
				Workspace: "42born2code",
				Username:  "Evaluation Master",
			},
			EmailSuffix: "student.42campus.org",
			WarnBefore:  time.Minute * 15,
			HTTPAddr:    "0.0.0.0:5000",
			Timeout:     time.Second * 10,
			Intra: Intra{
				AppID:     "--FILL ME--",
				AppSecret: "--FILL ME--",
				Webhooks: map[string]string{
					"--FILL": "ME--",
				},
			},
			Postgres: Database{
				Host:     "localhost",
				Port:     "5432",
				DB:       "postgres",
				User:     "postgres",
				Password: "--FILL ME--",
			},
			RabbitMQ: RabbitMQ{
				Host:     "localhost",
				Port:     "5672",
				VHost:    "",
				User:     "guest",
				Password: "guest",
				Queue:    "webhooks_intra_42jitsi",
			},
		}

		assert.NoError(t, Initiate())
		fmt.Println(viper.GetString("postgres.password"))
		assert.Equal(t, expected, Conf)
	})
}

func TestRabbitMQ_URL(t *testing.T) {
	rabbitmq := RabbitMQ{
		Host:     "localhost",
		Port:     "5672",
		VHost:    "vhost",
		User:     "user",
		Password: "password",
	}
	expected := "amqp://user:password@localhost:5672/vhost"
	assert.Equal(t, expected, rabbitmq.URL())
}
