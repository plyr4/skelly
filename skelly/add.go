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

	// parse and validate input
	emoji, usergroup, err := parseReactionSubCommandArgs(args)
	if err != nil {

		// invalid command args, send help
		err := sendHelp(command, s.ResponseURL)
		if err != nil {
			err = errors.Wrap(err, "could not send help")
		}

		return err
	}

	channel := s.ChannelID
	user := s.UserID
	triggerID := s.TriggerID

	// parse emoji input
	eID, err := parseEmoji(emoji)
	if err != nil {

		err = errors.Wrap(err, "could not parse emoji id")

		// notify user
		e := util.SendError(fmt.Sprintf("Sorry, is %s a valid emoji? Try using the autocomplete, starting with `:` !", emoji), channel, user)
		if e != nil {
			err = errors.Wrap(err, "could not send error: "+e.Error())
		}
		return err
	}

	// parse usergroup input
	ugID, _, err := parseUserGroup(usergroup)
	if err != nil {

		err = errors.Wrap(err, "could not parse usergroup id")

		// notify user
		e := util.SendError(fmt.Sprintf("Sorry, is %s a valid usergroup? Try using the autocomplete. starting with `@` !", usergroup), channel, user)
		if e != nil {
			err = errors.Wrap(err, "could not send error: "+e.Error())
		}
		return err
	}

	// attempt to retrieve an existing reaction
	exists, _, err := db.ReactionExists(channel, eID, ugID)
	if err != nil {
		err = errors.Wrap(err, "could not check for reaction")
		return err
	}

	// if reaction exists
	if exists {

		logrus.Infof("reaction already exists for channel(%s) emoji(%s) usergroup(%s)", channel, eID, ugID)

		// notify user
		err = util.SendError("Sorry, that reaction already exists. Did you mean to update?", channel, user)
		if err != nil {
			err = errors.Wrap(err, "could not send error")
			return err
		}

		return nil
	}

	// build default modal
	// uses channel and slash command as metadata
	metadata := strings.Join([]string{addSubCommand, emoji, usergroup, channel}, " ")

	modal := modal(addSubCommand,
		"Add a reaction to this channel.",
		metadata, emoji, usergroup, "")

	logrus.Infof("opening add modal for channel(%s) emoji(%s) usergroup(%s) trigger_id(%s)", channel, emoji, usergroup, triggerID)

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
	// ex: EMOJI:CHANNEL_ID
	channel, emoji, usergroup, err := parseViewMetadata(view)
	if err != nil {
		err = errors.Wrap(err, "could not parse metadata")
		return err
	}

	logrus.Infof("parsed metadata channel(%s) emoji(%s) usergroup(%s)", channel, emoji, usergroup)

	// parse out information from usergroup
	eID, err := parseEmoji(emoji)
	if err != nil {
		err = errors.Wrap(err, "could not parse emoji")
		return err
	}

	// parse out information from usergroup
	ugID, ugHandle, err := parseUserGroup(usergroup)
	if err != nil {
		err = errors.Wrap(err, "could not parse usergroup")
		return err
	}

	// check if user entered "none"
	if usergroup == "none" {
		ugID = "none"
	}

	// check for reaction in the database
	exists, _, err := db.ReactionExists(channel, eID, ugID)
	if err != nil {
		err = errors.Wrap(err, "could not check for reaction in db")
		return err
	}

	if exists {

		logrus.Infof("reaction already exists for channel(%s) emoji(%s) usergroup(%s)", channel, eID, ugID)

		// notify user
		err = util.SendError("Sorry, that reaction already exists. Did you mean to update?", channel, user)
		if err != nil {
			err = errors.Wrap(err, "could not send error")
			return err
		}

		return nil
	}

	// add reaction to the database
	err = db.AddReaction(channel, eID, ugID, usergroup, response)
	if err != nil {
		err = errors.Wrap(err, "could not add reaction to db")
		return err
	}

	// set response based on usergroup provided
	var text *slack.TextBlockObject
	if usergroup != "none" {
		text = slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("Okay, I will respond to %s for users in @%s.", emoji, ugHandle), false, false)
	} else {
		text = slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("Okay, I will respond to %s for *all users* in this channel.", emoji), false, false)
	}

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
