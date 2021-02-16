package skelly

import (
	"fmt"

	"github.com/davidvader/skelly/db"
	"github.com/davidvader/skelly/stats"
	"github.com/gosuri/uitable"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// List takes channel, emoji and message and adds a reaction to the database.
func List(bToken, channel string) error {

	// retrieve the appropriate reaction for the channel/emoji
	reactions, err := db.GetChannelReactions(channel)
	if err != nil {
		return err
	}

	// output reactions as a table
	table := uitable.New()
	table.MaxColWidth = 200
	table.Wrap = true // wrap columns

	table.AddRow(fmt.Sprintf("Reactions for channel(%s)", channel))
	table.AddRow("EMOJI", "USERGROUP", "RESPONSE")

	for _, r := range *reactions {
		table.AddRow(r.Emoji, r.UserGroup, r.Response)
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

// View takes channel and emoji and retrieves the appropriate response
func View(channel, emoji, usergroup string) error {

	// retrieve reaction from db
	reaction, err := db.GetReaction(channel, emoji, usergroup)
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

// Add takes channel, emoji and response and adds a reaction to the database.
func Add(bToken, channel, emoji, usergroup, response string) error {

	// add the appropriate reaction for the channel/emoji/msg
	err := db.AddReaction(channel, emoji, usergroup, usergroup, response)
	if err != nil {
		err = errors.Wrap(err, "could not add reaction to db")
		return err
	}

	logrus.Infof("reaction added for channel(%s) emoji(%s) usergroup(%s) response(%s)", channel, emoji, usergroup, response)
	return nil
}

// Update takes channel, emoji and response and updates a reaction in the database.
func Update(bToken, channel, emoji, usergroup, response string) error {

	// update the appropriate reaction for the channel/emoji
	err := db.UpdateReaction(channel, emoji, usergroup, response)
	if err != nil {
		err = errors.Wrap(err, "could not update reaction in db")
		return err
	}

	logrus.Infof("reaction updated for channel(%s) emoji(%s) usergroup(%s) response(%s)", channel, emoji, usergroup, response)
	return nil
}

// Delete takes channel and emoji and deletes reactions from the database.
func Delete(bToken, channel, emoji, usergroup string) error {

	// delete the appropriate reactions for the channel/emoji
	n, err := db.DeleteReactions(channel, emoji, usergroup)
	if err != nil {
		err = errors.Wrap(err, "could not delete reaction from db")
		return err
	}

	logrus.Infof("(%v) reactions deleted for channel(%s) emoji(%s)", n, channel, emoji)
	return nil
}

// Trigger takes post parameters and posts a reaction following any rules specified for that channel, emoji and usergroup.
func Trigger(bToken, channel, emoji, user, ts string) error {

	// post the appropriate reactions for the channel/emoji/usergroup/ts
	err := React(bToken, channel, emoji, user, ts)
	if err != nil {
		logrus.Infof("could not post reaction for channel(%s) emoji(%s) user(%s) ts(%s)", channel, emoji, user, ts)
		return err
	}

	logrus.Infof("reactions posted for channel(%s) emoji(%s) user(%s) ts(%s)", channel, emoji, user, ts)
	return nil
}

// ChannelStats takes channel and prints skelly statistics for the specified channel.
func ChannelStats(bToken, channel string) error {

	// retrieve stats for the specified channel
	stats, err := stats.GetChannelStats(channel)
	if err != nil {
		logrus.Infof("could not get channel stats for channel(%s)", channel)
		return err
	}

	logrus.Infof("stats for channel(%s)", channel)

	// output stats as a table
	table := uitable.New()
	table.MaxColWidth = 200
	table.Wrap = true // wrap columns

	table.AddRow(fmt.Sprintf("Stats for channel(%s)", channel))
	table.AddRow("TOTAL RULES", stats.TotalRules)

	// add a row of space at the bottom
	table.AddRow()

	// print the table
	fmt.Println(table)

	return nil
}

// GenericStats prints skelly statistics for the configured workspace.
func GenericStats(bToken string) error {

	// retrieve stats for the specified channel
	stats, err := stats.GetSkellyStats(bToken)
	if err != nil {
		logrus.Infof("could not get Skelly stats")
		return err
	}

	logrus.Infof("stats for all channels")

	// output stats as a table
	table := uitable.New()
	table.MaxColWidth = 200
	table.Wrap = true // wrap columns

	table.AddRow("Stats for Skelly")
	table.AddRow("TOTAL CHANNELS", stats.TotalChannels)
	table.AddRow("CHANNEL NAMES")

	for _, channel := range stats.Channels {
		table.AddRow(channel.NameNormalized, channel.ID)
	}
	// add a row of space at the bottom
	table.AddRow()

	// print the table
	fmt.Println(table)

	return nil
}
