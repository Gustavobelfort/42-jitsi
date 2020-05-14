package config

import (
	"time"

	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Conf initializes a global variable
// where the configurations gathered from the config file will be stored
var Conf Configuration

var requiredConf = []string{"postgres.password"}

// AddRequired will add more required variables to the list of required configuration.
func AddRequired(required ...string) {
	requiredConf = append(requiredConf, required...)
}

func setDefaults() {
	log.Debugln("setting config defaults value")
	viper.SetTypeByDefaultValue(true)

	viper.SetDefault("config_file", "./config.yml")

	viper.SetDefault("environment", "development")
	viper.SetDefault("service", "42-jitsi")

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

	viper.SetDefault("rabbitmq.host", "localhost")
	viper.SetDefault("rabbitmq.port", "5672")
	viper.SetDefault("rabbitmq.vhost", "")
	viper.SetDefault("rabbitmq.user", "guest")
	viper.SetDefault("rabbitmq.password", "guest")
	viper.SetDefault("rabbitmq.queue", "webhooks_intra_42jitsi")

	viper.SetDefault("log_level", log.DebugLevel)

	viper.SetDefault("logstash.host", "localhost")
	viper.SetDefault("logstash.port", "5000")
	viper.SetDefault("logstash.protocol", "tcp")
	viper.SetDefault("logstash.levels", []log.Level{log.InfoLevel, log.WarnLevel, log.ErrorLevel, log.FatalLevel, log.PanicLevel})
	viper.SetDefault("logstash.enabled", false)

	viper.SetDefault("sentry.levels", []log.Level{log.ErrorLevel, log.FatalLevel, log.PanicLevel})
	viper.SetDefault("sentry.enabled", false)
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

	// Workaround until unmarshal is supported for unbound environmental variables
	// See github issue: https://github.com/spf13/viper/issues/761
	logBinding("environment", "ENVIRONMENT")
	logBinding("service", "SERVICE")

	logBinding("email_suffix", "EMAIL_SUFFIX")

	logBinding("warn_before", "WARN_BEFORE")

	logBinding("timeout", "TIMEOUT")

	logBinding("postgres.host", "POSTGRES_HOST")
	logBinding("postgres.port", "POSTGRES_PORT")
	logBinding("postgres.db", "POSTGRES_DB")
	logBinding("postgres.user", "POSTGRES_USER")

	logBinding("slack_that.url", "SLACK_THAT_URL")
	logBinding("slack_that.username", "SLACK_THAT_USERNAME")

	logBinding("rabbitmq.host", "RABBITMQ_HOST")
	logBinding("rabbitmq.port", "RABBITMQ_PORT")
	logBinding("rabbitmq.vhost", "RABBITMQ_VHOST")
	logBinding("rabbitmq.user", "RABBITMQ_USER")
	logBinding("rabbitmq.password", "RABBITMQ_PASSWORD")
	logBinding("rabbitmq.queue", "RABBITMQ_QUEUE")

	logBinding("log_level", "LOG_LEVEL")

	logBinding("logstash.host", "LOGSTASH_HOST")
	logBinding("logstash.port", "LOGSTASH_PORT")
	logBinding("logstash.protocol", "LOGSTASH_PROTOCOL")
	logBinding("logstash.levels", "LOGSTASH_LEVELS")
	logBinding("logstash.enabled", "LOGSTASH_ENABLED")

	logBinding("sentry.levels", "SENTRY_LEVELS")
	logBinding("sentry.enabled", "SENTRY_ENABLED")
	// End workaround

	logBinding("sentry.dsn", "SENTRY_DSN")

	logBinding("postgres.password", "POSTGRES_PASSWORD")
	logBinding("intra.app_id", "INTRA_APP_ID")
	logBinding("intra.app_secret", "INTRA_APP_SECRET")
	logBinding("intra.webhooks", "INTRA_WEBHOOKS")

	logBinding("slack_that.workspace", "SLACK_THAT_WORKSPACE")

}

func loadFile() {
	viper.SetConfigType("yaml")

	filename := viper.GetString("config_file")
	if filename == "" {
		log.Info("'config_file' is empty")
		return
	}
	viper.SetConfigFile(viper.GetString("config_file"))
	if err := viper.ReadInConfig(); err != nil {
		log.WithField("config_file", filename).WithError(err).Errorf("loading config file: %v", err)
		return
	}
	log.WithField("config_file", filename).Infof("loaded config from: '%s'", filename)
}

func unmarshalConfig() error {
	decodeHook := mapstructure.ComposeDecodeHookFunc(
		stringToMapstringHookFunc,
		stringToLogLevelHookFunc,
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
	loadFile()

	if err := checkRequired(requiredConf...); err != nil {
		return err
	}

	Conf = Configuration{}
	return unmarshalConfig()
}
