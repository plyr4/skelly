package db

import (
	"fmt"

	"github.com/davidvader/skelly/types"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

// GetChannelReactions retrieve reactions for a channel from the db
func GetChannelReactions(channel string) (*[]types.Reaction, error) {

	logrus.Infof("getting reactions for channel(%s)", channel)

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
	err = col.Find(bson.M{"channel": channel}).All(&reactions)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not get reactions from db for channel(%s)", channel))
	}

	return &reactions, nil
}

// GetReactions retrieves reactions for a channel from the db
func GetReactions(channel string) ([]*types.Reaction, error) {

	logrus.Infof("getting reactions for channel(%s)", channel)

	// connect to mongo
	session, err := connect()
	if err != nil {
		return nil, errors.Wrap(err, "could not connect to db")
	}
	defer session.Close()

	// retrieve the collection
	col := session.DB(getConfig().DB).C(collection)

	reactions := []*types.Reaction{}

	// retrieve the reaction from the db
	err = col.Find(reactionsSelector(channel)).All(&reactions)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not get reactions from db for channel(%s)", channel))
	}

	return reactions, nil
}

// GetReaction retrieve reaction for a channel from the db
func GetReaction(channel string) (*types.Reaction, error) {

	logrus.Infof("getting a reaction for channel(%s)", channel)

	// connect to mongo
	session, err := connect()
	if err != nil {
		return nil, errors.Wrap(err, "could not connect to db")
	}
	defer session.Close()

	// retrieve the collection
	col := session.DB(getConfig().DB).C(collection)

	reaction := types.Reaction{}

	// retrieve the reaction from the db
	err = col.Find(reactionSelector(channel)).One(&reaction)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not get reaction from db for channel(%s)", channel))
	}

	return &reaction, nil
}

// AddReaction adds a reaction for a channel to the db
func AddReaction(channel, response string) error {

	logrus.Infof("adding a reaction for channel(%s) response(%s)", channel, response)

	// connect to mongo
	session, err := connect()
	if err != nil {
		return errors.Wrap(err, "could not connect to db")
	}
	defer session.Close()

	// retrieve the collection
	col := session.DB(getConfig().DB).C(collection)

	// TODO: improve the use of .All() as .One() check
	reactions := []types.Reaction{}

	// retrieve the reactions from the db
	err = col.Find(reactionSelector(channel)).All(&reactions)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("could not get reaction from db for channel(%s)", channel))
	}

	// if it exists, do not add it
	if len(reactions) != 0 {
		return fmt.Errorf("reaction already exists for channel(%s)", channel)
	}

	// save data into Reaction struct
	reaction := types.Reaction{
		Channel:  channel,
		Response: response,
	}

	// insert reaction into db
	err = col.Insert(reaction)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("could not insert reaction into db for channel(%s) response(%s)", channel, response))
	}

	return nil
}

// UpdateReaction retrieve and updates a reaction for a channel from the db
func UpdateReaction(channel, response string) error {

	logrus.Infof("updating a reaction for channel(%s) response(%s)", channel, response)

	// connect to mongo
	session, err := connect()
	if err != nil {
		return errors.Wrap(err, "could not connect to db")
	}
	defer session.Close()

	// retrieve the collection
	col := session.DB(getConfig().DB).C(collection)

	// TODO: improve the use of .All() as .One() check
	reactions := []types.Reaction{}

	// retrieve the reactions from the db
	err = col.Find(reactionSelector(channel)).All(&reactions)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("could not get reactions from db for channel(%s)", channel))
	}

	// if it does not exist, do not update it
	if len(reactions) == 0 {
		return fmt.Errorf("reaction does not exist for channel(%s)", channel)
	}

	reaction := reactions[0]

	if reaction.Response != response {
		reaction.Response = response
	}

	// update reaction in db
	err = col.Update(reactionSelector(channel), reaction)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("could not update reaction in db for channel(%s)", reaction.Channel))
	}

	return nil
}

// DeleteReactions retrieve and deletes reactions for a channel from the db
func DeleteReactions(channel string) (int, error) {

	logrus.Infof("removing reactions for channel(%s)", channel)

	// connect to mongo
	session, err := connect()
	if err != nil {
		return 0, errors.Wrap(err, "could not connect to db")
	}
	defer session.Close()

	// retrieve the collection
	col := session.DB(getConfig().DB).C(collection)

	reactions := []types.Reaction{}

	// retrieve the reactions from the db
	err = col.Find(reactionSelector(channel)).All(&reactions)
	if err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("could not get reaction from db for channel(%s)", channel))
	}

	// they do not exist, do not remove them
	if len(reactions) == 0 {
		return 0, fmt.Errorf("reactions do not exist for channel(%s)", channel)
	}

	// remove reactions from db
	_, err = col.RemoveAll(reactionSelector(channel))
	if err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("could not delete reactions from db for channel(%s)", channel))
	}

	return len(reactions), nil
}

// DeleteChannelReactions retrieve and deletes reactions for a channel from the db
func DeleteChannelReactions(channel string) (int, error) {

	logrus.Infof("removing reactions for channel(%s)", channel)

	// connect to mongo
	session, err := connect()
	if err != nil {
		return 0, errors.Wrap(err, "could not connect to db")
	}
	defer session.Close()

	// retrieve the collection
	col := session.DB(getConfig().DB).C(collection)

	reactions := []types.Reaction{}

	// retrieve the reactions from the db
	err = col.Find(channelSelector(channel)).All(&reactions)
	if err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("could not get reaction from db for channel(%s)", channel))
	}

	// they do not exist, do not remove them
	if len(reactions) == 0 {
		return 0, fmt.Errorf("reactions do not exist for channel(%s)", channel)
	}

	// remove reactions from db
	_, err = col.RemoveAll(channelSelector(channel))
	if err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("could not delete reactions from db for channel(%s)", channel))
	}

	return len(reactions), nil
}

// ReactionExists checks for reaction for a channel in the db
func ReactionExists(channel string) (bool, *types.Reaction, error) {

	logrus.Infof("checking for reaction channel(%s)", channel)

	// connect to mongo
	session, err := connect()
	if err != nil {
		return false, nil, errors.Wrap(err, "could not connect to db")
	}
	defer session.Close()

	// retrieve the collection
	col := session.DB(getConfig().DB).C(collection)

	reactions := []types.Reaction{}

	// retrieve the reaction from the db
	err = col.Find(reactionSelector(channel)).All(&reactions)
	if err != nil {
		return false, nil, errors.Wrap(err, fmt.Sprintf("could not get reaction from db for channel(%s)", channel))
	}

	// check for reaction
	exists := false
	var r *types.Reaction

	if len(reactions) > 0 {
		exists = true
		r = &reactions[0]
	}

	return exists, r, nil
}

// StoreResponse stores a response for a channel/user/timestamp in the db
func StoreResponse(channel, user, timestamp string) error {

	logrus.Infof("storing response for channel(%s) user(%s) timestamp(%s)", channel, user, timestamp)

	// connect to mongo
	session, err := connect()
	if err != nil {
		return errors.Wrap(err, "could not connect to db")
	}
	defer session.Close()

	// retrieve the collection
	col := session.DB(getConfig().DB).C(responseCollection)

	// TODO: improve the use of .All() as .One() check
	responses := []types.Response{}

	// retrieve the reactions from the db
	err = col.Find(responseSelector(channel, user, timestamp)).All(&responses)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("could not get response from db for channel(%s) user(%s) timestamp(%s)", channel, user, timestamp))
	}

	// if it exists, do not add it
	if len(responses) != 0 {
		return fmt.Errorf("response already exists for channel(%s) user(%s) timestamp(%s)", channel, user, timestamp)
	}

	// save data into Response struct
	response := types.Response{
		Channel:   channel,
		User:      user,
		Timestamp: timestamp,
	}

	// insert reaction into db
	err = col.Insert(response)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("could not insert response into db for channel(%s) user(%s) timestamp(%s)", channel, user, timestamp))
	}

	return nil
}

// CheckResponse checks to see if a response for a channel/user/timestamp exits in the db
func CheckResponse(channel, user, timestamp string) (bool, error) {

	logrus.Infof("checking for response for channel(%s) user(%s) timestamp(%s)", channel, user, timestamp)

	// connect to mongo
	session, err := connect()
	if err != nil {
		return false, errors.Wrap(err, "could not connect to db")
	}
	defer session.Close()

	// retrieve the collection
	col := session.DB(getConfig().DB).C(responseCollection)

	// TODO: improve the use of .All() as .One() check
	responses := []types.Response{}

	// retrieve the reactions from the db
	err = col.Find(responseSelector(channel, user, timestamp)).All(&responses)
	if err != nil {
		return false, errors.Wrap(err, fmt.Sprintf("could not get response from db for channel(%s) user(%s) timestamp(%s)", channel, user, timestamp))
	}

	return len(responses) != 0, nil
}
