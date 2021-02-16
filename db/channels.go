package db

import (
	"github.com/davidvader/skelly/types"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// GetChannels retrieve a map for channels to reactions from the db
func GetChannels() (*map[string]int, error) {

	logrus.Infof("getting all channels")

	// connect to mongo
	session, err := connect()
	if err != nil {
		return nil, errors.Wrap(err, "could not connect to db")
	}
	defer session.Close()

	// retrieve the collection
	col := session.DB(getConfig().DB).C(collection)

	reactions := []types.Reaction{}

	// retrieve the reactions from the db
	err = col.Find(nil).All(&reactions)
	if err != nil {
		return nil, errors.Wrap(err, "could not get reactions from db for all channels")
	}

	logrus.Tracef("retrieved %v reactions from the db", len(reactions))

	// create a map from id to number of rules
	channelRules := map[string]int{}

	for _, reaction := range reactions {

		// set or increment channel rules
		_, ok := channelRules[reaction.Channel]
		if !ok {
			channelRules[reaction.Channel] = 1
		} else {
			channelRules[reaction.Channel]++
		}
	}

	return &channelRules, nil
}
