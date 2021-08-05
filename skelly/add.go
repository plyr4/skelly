package skelly

import (
	"fmt"
	"os"
	"strings"

	"github.com/davidvader/skelly/db"
	"github.com/davidvader/skelly/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

// openAddModal takes slash command configuration and responds
// to the triggering user with a dialog window for adding a new
// reaction to the skelly database
func openAddModal(s *slack.SlashCommand, command string, args []string) error {

	channel := s.ChannelID
	user := s.UserID
	triggerID := s.TriggerID

	// attempt to retrieve an existing reaction
	exists, _, err := db.ReactionExists(channel)
	if err != nil {
		err = errors.Wrap(err, "could not check for reaction")
		return err
	}

	// if reaction exists
	if exists {

		logrus.Infof("reaction already exists for channel(%s)", channel)

		// notify user
		err = util.SendError("Sorry, a typing reaction already exists for this channel. Did you mean to update?", channel, user)
		if err != nil {
			err = errors.Wrap(err, "could not send error")
			return err
		}

		return nil
	}

	// build default modal
	// uses channel and slash command as metadata
	metadata := strings.Join([]string{addSubCommand, channel}, " ")

	modal := modal(addSubCommand,
		"Add a reaction to this channel.",
		metadata, "")

	logrus.Infof("opening add modal for channel(%s) trigger_id(%s)", channel, triggerID)

	// create an api client
	bToken := os.Getenv("SKELLY_BOT_TOKEN")

	api := slack.New(bToken)

	// open modal view
	_, err = api.OpenView(triggerID, modal)
	if err != nil {
		err = errors.Wrap(err, "could not open view")
		return err
	}
	return nil
}

// handleAddSubmission takes slack view, extracts args, and attempts to add a reaction to the database
func handleAddSubmission(view *slack.View, user, responseURL string) error {

	// parse submission value
	response, err := parseViewResponse(view)
	if err != nil {
		err = errors.Wrap(err, "could not parse response")
		return err
	}

	// parse out args from private metadata
	// ex: META:CHANNEL_ID
	channel, err := parseViewMetadata(view)
	if err != nil {
		err = errors.Wrap(err, "could not parse metadata")
		return err
	}

	logrus.Infof("parsed metadata channel(%s)", channel)

	// check for reaction in the database
	exists, _, err := db.ReactionExists(channel)
	if err != nil {
		err = errors.Wrap(err, "could not check for reaction in db")
		return err
	}

	if exists {

		logrus.Infof("reaction already exists for channel(%s)", channel)

		// notify user
		err = util.SendError("Sorry, that reaction already exists. Did you mean to update?", channel, user)
		if err != nil {
			err = errors.Wrap(err, "could not send error")
			return err
		}

		return nil
	}

	// add reaction to the database
	err = db.AddReaction(channel, response)
	if err != nil {
		err = errors.Wrap(err, "could not add reaction to db")
		return err
	}

	// set response
	var text *slack.TextBlockObject = slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("Okay, I will respond to all users that type in this channel, once a day."), false, false)

	section := slack.NewSectionBlock(text, nil, nil)

	// create default msg options
	options := []slack.MsgOption{
		slack.MsgOptionBlocks(section),
	}

	// create an api client
	bToken := os.Getenv("SKELLY_BOT_TOKEN")

	api := slack.New(bToken)

	// post the confirmation
	_, err = api.PostEphemeral(channel, user, options...)
	if err != nil {
		err = errors.Wrap(err, "could not post response")
		return err
	}
	return nil
}
