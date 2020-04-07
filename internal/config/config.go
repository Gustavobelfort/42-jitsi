package config

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"

	"github.com/gustavobelfort/42-jitsi/internal/utils"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

const configFile = "config.yml"

// Conf initializes a global variable
// where the configurations gathered from the config file will be stored
var Conf Configuration

// Configuration is the type that will hold the configuration informations
type Configuration struct {
	SlackThatURL      string        `mapstructure:"slack-that-url"`
	CampusSlug        string        `mapstructure:"campus-slug"`
	WarnBefore        time.Duration `mapstructure:"warn-before"`
	IntraWebhooksAuth []Webhook     `mapstructure:"intra-webhooks"`
	Postgres          Database
}

type Webhook struct {
	Hook   string
	Secret string
}

// Database is the type that will hold the database informations
type Database struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
}

func check() error {
	fd, err := filepath.Abs(configFile)
	if err != nil {
		return err
	}

	if _, err = ioutil.ReadFile(fd); err != nil {
		return err
	}

	return nil
}

func load(filename string) (Configuration, error) {

	viper.SetConfigType("yaml")

	fd, err := filepath.Abs(filename)
	if err != nil {
		return Conf, err
	}

	ymlFile, err := ioutil.ReadFile(fd)
	if err != nil {
		return Conf, err
	}

	viper.AutomaticEnv()
	viper.ReadConfig(bytes.NewBuffer(ymlFile))
	switch {
	case err != nil:
		log.Println(err)
	default:
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	err = viper.Unmarshal(&Conf)
	if err != nil {
		log.Println(errors.Wrap(err, "unmarshal config file"))
	}

	log.Println("config (info) config file successfully loaded.")
	return Conf, nil

}

// Initiate checks for the config file, and if its its found, try to load it into the program
func Initiate() {
	log.Println("config (info) loading config...")
	if err := check(); err != nil {
		log.Fatalf("%sconfig (error)%s can't access '%s'. (%v)\n", utils.Red, utils.Reset, configFile, err)
	}
	load(configFile)
}
