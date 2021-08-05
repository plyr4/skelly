package skelly

import (
	"fmt"

	"github.com/davidvader/skelly/db"
	"github.com/gosuri/uitable"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// List takes channel and message and adds a reaction to the database.
func List(bToken, channel string) error {

	// retrieve the appropriate reaction for the channel
	reactions, err := db.GetChannelReactions(channel)
	if err != nil {
		return err
	}

	// output reactions as a table
	table := uitable.New()
	table.MaxColWidth = 200
	table.Wrap = true // wrap columns

	table.AddRow(fmt.Sprintf("Reactions for channel(%s)", channel))
	table.AddRow("RESPONSE")

	for _, r := range *reactions {
		table.AddRow(r.Response)
	}

	// add a row of space at the bottom
	table.AddRow()

	// print the table
	fmt.Println(table)

	return nil
}

// Clear takes channel and removes reactions from the database.
func Clear(channel string) error {

	// delete reactions from the database
	n, err := db.DeleteChannelReactions(channel)
	if err != nil {
		err = errors.Wrap(err, "could not delete reactions from db")
		return err
	}

	logrus.Infof("removed (%v) reactions for channel(%s)", n, channel)

	return nil
}

// View takes channel and retrieves the appropriate response
func View(channel string) error {

	// retrieve reaction from db
	reaction, err := db.GetReaction(channel)
	if err != nil {
		err = errors.Wrap(err, "could not get reaction from db")
		return err
	}

	// use yaml as output format
	output, err := yaml.Marshal(reaction)
	if err != nil {
		err = errors.Wrap(err, "could not yaml marshal")
		return err
	}

	fmt.Println(string(output))
	return nil
}

// Add takes channel and response and adds a reaction to the database.
func Add(bToken, channel, response string) error {

	// add the appropriate reaction for the channel/msg
	err := db.AddReaction(channel, response)
	if err != nil {
		err = errors.Wrap(err, "could not add reaction to db")
		return err
	}

	logrus.Infof("reaction added for channel(%s) response(%s)", channel, response)
	return nil
}

// Update takes channel and response and updates a reaction in the database.
func Update(bToken, channel, response string) error {

	// update the appropriate reaction for the channel
	err := db.UpdateReaction(channel, response)
	if err != nil {
		err = errors.Wrap(err, "could not update reaction in db")
		return err
	}

	logrus.Infof("reaction updated for channel(%s) response(%s)", channel, response)
	return nil
}

// Delete takes channel and deletes reactions from the database.
func Delete(bToken, channel string) error {

	// delete the appropriate reactions for the channel
	n, err := db.DeleteReactions(channel)
	if err != nil {
		err = errors.Wrap(err, "could not delete reaction from db")
		return err
	}

	logrus.Infof("(%v) reactions deleted for channel(%s)", n, channel)
	return nil
}

// Trigger takes post parameters and posts a reaction following any rules specified for that channel.
func Trigger(bToken, channel, user, ts string) error {

	// post the appropriate reactions for the channel/ts
	err := React(bToken, channel, user, ts)
	if err != nil {
		logrus.Infof("could not post reaction for channel(%s) user(%s) ts(%s)", channel, user, ts)
		return err
	}

	logrus.Infof("reactions posted for channel(%s) user(%s) ts(%s)", channel, user, ts)
	return nil
}
