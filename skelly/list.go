package skelly

import (
	"fmt"

	"github.com/davidvader/skelly/db"
	"github.com/davidvader/skelly/types"
	"github.com/davidvader/skelly/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

// listReactions takes slash command configuration and responds
// to the triggering user with a list of the reactions for that channel
func listReactions(s *slack.SlashCommand, command string, args []string) error {

	// parse and validate input
	err := parseListSubCommandArgs(args)
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

	// attempt to retrieve an existing reaction
	reactions, err := db.GetChannelReactions(channel)
	if err != nil {
		err = errors.Wrap(err, "could not check for reaction")
		return err
	}

	// if reaction exists
	if len(*reactions) == 0 {

		logrus.Infof("no reactions exist for channel(%s)", channel)

		// notify user
		err = util.SendError("Sorry, no reactions exist for this channel.", channel, user)
		if err != nil {
			err = errors.Wrap(err, "could not send error")
			return err
		}

		return nil
	}

	logrus.Infof("listing (%v) reactions for channel(%s)", len(*reactions), channel)

	// build slack response
	response := listResponse(*reactions)

	// send response
	err = util.Respond(s.ResponseURL, response)
	if err != nil {
		err = errors.Wrap(err, "could not respond with reaction list")
		return err
	}
	return nil
}

// listResponse takes list of reactions and builds a slack message for listing them
func listResponse(reactions []types.Reaction) slack.Message {

	// header
	t := slack.NewTextBlockObject("mrkdwn",
		"Here are all of the reactions for this channel.",
		false, false)

	header := slack.NewSectionBlock(t, nil, nil)

	// init blocks for building the response
	blocks := []slack.Block{header}

	// adhere to the blocks limit
	if len(reactions) > 49 {
		reactions = reactions[:49]
	}

	// build blocks
	for _, r := range reactions {

		// build a view for the reaction
		t := slack.NewTextBlockObject("mrkdwn",
			fmt.Sprintf("*Emoji*: :%s: (%s) *Usergroup*: %s\n*Response*: %s",
				r.Emoji, r.Emoji, r.UserGroupFull, r.Response),
			false, false)

		block := slack.NewSectionBlock(t, nil, nil)

		// add block
		blocks = append(blocks, block)
	}

	// construct message
	msg := slack.NewBlockMessage(blocks...)

	return msg
}
