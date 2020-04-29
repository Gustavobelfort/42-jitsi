package config

import "time"

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
