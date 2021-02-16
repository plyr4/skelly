package db

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
)

// connect starts a session with the mongo db
func connect() (*mgo.Session, error) {

	// retrieving database config
	c := getConfig()

	logrus.Infof("connecting to mongo db(%s:%s:%s)", c.Host, c.DB, c.Username)

	// connect to mongo db
	s, err := mgo.Dial(c.toURI())
	if err != nil {
		return nil, errors.Wrap(err, "could not connect to mongo db")
	}
	return s, nil
}
