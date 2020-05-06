package config

import (
	"errors"
	"time"

	"github.com/gustavobelfort/42-jitsi/internal/logging"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Conf initializes a global variable
// where the configurations gathered from the config file will be stored
var Conf Configuration

func setDefaults() {
	log.Debugln("setting config defaults value")
	viper.SetTypeByDefaultValue(true)

	viper.SetDefault("config_file", "./config.yml")

	viper.SetDefault("email_suffix", "student.42campus.org")

	viper.SetDefault("warn_before", time.Minute*15)

	viper.SetDefault("http_addr", "0.0.0.0:5000")

	viper.SetDefault("timeout", time.Second*10)

	viper.SetDefault("postgres.host", "localhost")
	viper.SetDefault("postgres.port", 5432)
	viper.SetDefault("postgres.db", "postgres")
	viper.SetDefault("postgres.user", "postgres")

	viper.SetDefault("slack_that.url", "http://localhost:8080")
	viper.SetDefault("slack_that.username", "Evaluation Master")
}

func bindEnv() {
	log.Debugln("binding env vars to their configuration")

	viper.AllowEmptyEnv(true)
	viper.AutomaticEnv() // It is not necessary to bind what already has been defined in defaults

	logBinding := func(config, env string) {
		if err := viper.BindEnv(config, env); err != nil {
			log.Warnf("binding env vars: could not bind '%s' to '%s'\n", env, config)
			return
		}
		log.Debugf("bound '%s' to '%s'", env, config)
	}

	logBinding("postgres.password", "POSTGRES_PASSWORD")
	logBinding("intra.app_id", "INTRA_APP_ID")
	logBinding("intra.app_secret", "INTRA_APP_SECRET")
	logBinding("intra.webhooks", "INTRA_WEBHOOKS")

	logBinding("slack_that.workspace", "SLACK_THAT_WORKSPACE")

}

func loadFile() error {
	viper.SetConfigType("yaml")

	filename := viper.GetString("config_file")
	if filename == "" {
		return logging.WithLog(errors.New("'config_file' is empty"), log.InfoLevel, nil)
	}
	viper.SetConfigFile(viper.GetString("config_file"))
	if err := viper.ReadInConfig(); err != nil {
		return logging.WithLog(err, log.ErrorLevel, log.Fields{"config_file": filename})
	}
	log.WithField("config_file", filename).Infof("loaded config from: '%s'\n", filename)
	return nil
}

func unmarshalConfig() error {
	decodeHook := mapstructure.ComposeDecodeHookFunc(
		stringToMapstringHookFunc,
		mapstructure.StringToSliceHookFunc(","),
		mapstructure.StringToTimeDurationHookFunc(),
	)
	return viper.Unmarshal(&Conf, viper.DecodeHook(decodeHook))
}

// Initiate checks for the config file, and if its its found, try to load it into the program
func Initiate() error {
	log.Infoln("loading configuration")
	setDefaults()
	bindEnv()
	err := loadFile()
	// logging.LogError won't log if err is nil.
	logging.LogError(log.StandardLogger(), err, "loading config file")

	if err := checkRequired("intra.app_id", "intra.app_secret", "intra.webhooks", "postgres.password", "slack_that.workspace"); err != nil {
		return err
	}

	Conf = Configuration{}
	return unmarshalConfig()
}
