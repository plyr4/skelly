package db

import (
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	// collection is the primary mongo db collection used for storing reactions
	collection = "reactions"
	// responseCollection is the mongo db collection to store reponses
	responseCollection = "responses"
	// dbTimeout is the primary mongo db collection used for storing reactions
	dbTimeout = 60 * time.Second
)

// global var for storing db Config
var config *Config
var once sync.Once

// Config is the struct representation for a monogodb connection configuration
type Config struct {
	Timeout  time.Duration
	Host     string
	DB       string
	Username string
	Password string
}

// setup uses environment to intialize the db configuration
func setup() {

	// retrieve db configurations from the environment

	// host
	host := os.Getenv("SKELLY_MONGO_HOST")

	// db
	database := os.Getenv("SKELLY_MONGO_DB")

	// auth
	username := os.Getenv("SKELLY_MONGO_USERNAME")
	password := os.Getenv("SKELLY_MONGO_PASSWORD")

	// set the mongo db configurations
	config = &Config{
		Timeout:  dbTimeout,
		Host:     host,
		DB:       database,
		Username: username,
		Password: password,
	}
}

// getConfig is a wrapper for retrieving the db config
// calls setup one time per application instance
func getConfig() *Config {
	once.Do(setup)
	return config
}

// toURI takes mongo config and returns the connection string
func (c *Config) toURI() string {
	return "mongodb://" + c.Username + ":" + url.QueryEscape(c.Password) + "@" + c.Host + "/" + c.DB
}

// Verify takes mongo connection config and verifies that it can connect to the database
func Verify() error {

	// retrieving database config
	c := getConfig()

	logrus.Infof("verifying mongo config(%s:%s:%s)", c.Host, c.DB, c.Username)

	// attempt to connect to the database
	_, err := connect()
	if err != nil {
		return errors.Wrap(err, "could not verify mongo config")
	}

	return nil
}
