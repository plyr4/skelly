package skelly

import (
	"os"
	"strings"

	"github.com/davidvader/skelly/db"
	"github.com/davidvader/skelly/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

// openDeleteModal takes slash command configuration and responds
// to the triggering user with a dialog window for deleting an existing
// reaction from the skelly database
func openDeleteModal(s *slack.SlashCommand, command string, args []string) error {

	channel := s.ChannelID
	user := s.UserID
	triggerID := s.TriggerID

	// attempt to retrieve an existing reaction
	exists, _, err := db.ReactionExists(channel)
	if err != nil {
		err = errors.Wrap(err, "could not check for reaction")
		return err
	}

	// if reaction does not exist
	if !exists {

		logrus.Infof("reaction does not exist for channel(%s)", channel)

		// notify user
		err = util.SendError("Sorry, that reaction does not exist for this channel. Did you mean to add?", channel, user)
		if err != nil {
			err = errors.Wrap(err, "could not send error")
			return err
		}

		return nil
	}

	// build default modal
	// uses channel and slash command as metadata
	metadata := strings.Join([]string{deleteSubCommand, channel}, " ")

	modal := deleteModal(deleteSubCommand,
		metadata)

	logrus.Infof("opening delete modal for channel(%s) trigger_id(%s)", channel, triggerID)

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

// handleDeleteSubmission takes slack view, extracts args, and attempts to delete a reaction from the database
func handleDeleteSubmission(view *slack.View, user, responseURL string) error {

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

	if !exists {

		logrus.Infof("reaction does not exist for channel(%s)", channel)

		// notify user
		err = util.SendError("Sorry, that reaction does not exist for this channel. Did you mean to add?", channel, user)
		if err != nil {
			err = errors.Wrap(err, "could not send error")
			return err
		}

		return nil
	}

	// delete reaction in the database
	n, err := db.DeleteReactions(channel)
	if err != nil {
		err = errors.Wrap(err, "could not delete reactions from db")
		return err
	}

	logrus.Infof("removed (%v) reactions for channel(%s)", n, channel)

	text := slack.NewTextBlockObject("mrkdwn", "I've deleted the reaction for this channel!", false, false)

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
