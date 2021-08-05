package skelly

import (
	"github.com/davidvader/skelly/util"
	"github.com/pkg/errors"
	"github.com/slack-go/slack"
)

// sendHelp responds to /skelly help with details on how to use Skelly via Slack
func sendHelp(command, responseURL string) error {

	// echo the given command
	given := slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", command, false, false), nil, nil)

	// header
	description := slack.NewTextBlockObject("mrkdwn", "I automatically react to typing.\nHere are some commands you can use the same as `/skelly help`.", false, false)
	header := slack.NewSectionBlock(description, nil, nil)

	// commands
	t := []*slack.TextBlockObject{
		slack.NewTextBlockObject("mrkdwn", "*Command*", false, false),
		slack.NewTextBlockObject("mrkdwn", "*Action*", false, false),
		slack.NewTextBlockObject("mrkdwn", "/skelly help", false, false),
		slack.NewTextBlockObject("mrkdwn", "prints commands and helpful information", false, false),
		slack.NewTextBlockObject("mrkdwn", "/skelly add", false, false),
		slack.NewTextBlockObject("mrkdwn", "trigger a response when users type in the channel", false, false),
	}
	commandsA := slack.NewSectionBlock(nil, t, nil)

	// split commands due to field limit
	t = []*slack.TextBlockObject{
		slack.NewTextBlockObject("mrkdwn", "/skelly update", false, false),
		slack.NewTextBlockObject("mrkdwn", "update a reaction in this channel", false, false),
		slack.NewTextBlockObject("mrkdwn", "/skelly delete", false, false),
		slack.NewTextBlockObject("mrkdwn", "delete a reaction in this channel", false, false),
		slack.NewTextBlockObject("mrkdwn", "/skelly list", false, false),
		slack.NewTextBlockObject("mrkdwn", "lists all reactions in this channel", false, false),
	}
	commandsB := slack.NewSectionBlock(nil, t, nil)

	// footer
	source := slack.NewTextBlockObject("mrkdwn", "Visit the <https://github.com/davidvader/skelly|GitHub  :github:>  repo for more info", false, false)
	footer := slack.NewContextBlock("help_footer", source)

	// construct blocks
	blocks := []slack.Block{
		given,
		header,
		commandsA,
		commandsB,
		footer,
	}

	// construct message
	msg := slack.NewBlockMessage(blocks...)

	// respond using response url
	err := util.Respond(responseURL, msg)
	if err != nil {
		err = errors.Wrap(err, "could not respond")
		return err
	}
	return nil
}
