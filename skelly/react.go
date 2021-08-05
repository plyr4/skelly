package skelly

import (
	"github.com/davidvader/skelly/db"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

// React takes channel and reacts with the appropriate response based on application configuration.
func React(bToken, channel, user, ts string) error {

	// retrieve all of the reactions for the channel
	reactions, err := db.GetReactions(channel)
	if err != nil {
		err = errors.Wrap(err, "could not get reaction from db")
		return err
	}

	logrus.Infof("retreived (%v) reactions for channel(%s)", len(reactions), channel)

	// create an api client
	api := slack.New(bToken)

	// filter the reactions based on user id and channel
	logrus.Infof("filtering reactions for channel(%s) user(%s)", channel, user)

	logrus.Infof("reacting to (%v) reactions for channel(%s)", len(reactions), channel)

	// respond to possibly multiple reactions
	for _, r := range reactions {

		// do not react if response is empty
		if len(r.Response) == 0 {
			continue
		}

		// check database for existing response
		exists, err := db.CheckResponse(channel, user, ts)
		if err != nil {
			err = errors.Wrap(err, "could not check for existing response")
			return err
		}

		// do not react if response already exists
		if exists {
			logrus.Infof("skipping, reaction response exists for channel(%s) user(%s) ts(%s)", channel, user, ts)
			continue
		}

		// create default msg options
		options := []slack.MsgOption{
			slack.MsgOptionText(r.Response, false),
			slack.MsgOptionPostMessageParameters(
				slack.PostMessageParameters{
					LinkNames: 1, UnfurlMedia: true,
				}),
			slack.MsgOptionEnableLinkUnfurl(),
		}

		// add timestamp if applicable
		if ts != "none" {
			options = append(options, slack.MsgOptionTS(ts))
		}

		// post the reaction
		logrus.Infof("posting reaction for channel(%s) user(%s) ts(%s)", channel, user, ts)

		_, mts, err := api.PostMessage(channel, options...)
		if err != nil {
			err = errors.Wrap(err, "could not post response")
			return err
		}

		err = db.StoreResponse(channel, user, ts)
		if err != nil {
			err = errors.Wrap(err, "response posted, but could not remember the response")
			return err
		}

		logrus.Infof("reaction posted for channel(%s) user(%s) ts(%s) at msg_ts(%s)", channel, user, ts, mts)
	}
	return nil
}
