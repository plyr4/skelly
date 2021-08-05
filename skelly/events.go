package skelly

import (
	"encoding/json"
	"net/http"

	"github.com/davidvader/skelly/util"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack/slackevents"
)

const (
	emojiChangedEvent = "emoji_changed"
)

// HandleEvent takes gin context and checks if request is a slack api url challenge
// if required, responds with the provided challenge string
func HandleEvent(c *gin.Context, body []byte, e *slackevents.EventsAPIEvent, bToken string) error {

	// verify the router url with the slack api, if needed
	verification, err := verifyURL(c, body, e.Type)
	if err != nil {
		err = errors.Wrap(err, "could not verify url")
		return err
	}

	// if the request is a url verification, exit
	if verification {
		logrus.Info("url verified")
		return nil
	}

	// acknowledge request
	util.RespondOK(c)

	go func() {

		// handle the inner callback event
		switch e.Type {
		case slackevents.CallbackEvent:

			// extract inner event
			innerEvent := e.InnerEvent

			switch innerEvent.Data.(type) {

			// reaction added event
			// case *slackevents.ReactionAddedEvent:

			// 	logrus.Infof("received reaction added event for event_ts(%s) item_type(%s) item_ts(%s)", ev.EventTimestamp, ev.Item.Type, ev.Item.Timestamp)

			// 	// react to the emoji
			// 	err := React(bToken, ev.Item.Channel, ev.Reaction, ev.User, ev.Item.Timestamp)
			// 	if err != nil {
			// 		err = errors.Wrap(err, "could not react")
			// 		logrus.Error(err)
			// 		return
			// 	}
			// 	return

			// unsupported inner event type
			default:
				logrus.Warn("received unsupported inner event callback type: ", e.Type)
				break
			}

		// unsupported outer event type
		default:
			logrus.Warn("received unsupported outer event type: ", e.Type)
			break
		}
	}()
	return nil
}

// verifyURL takes gin context and request body and verifies the challenge presented by the Slack API
func verifyURL(c *gin.Context, body []byte, eventType string) (bool, error) {

	// check if request is the slack api verifying the events url
	if eventType == slackevents.URLVerification {

		logrus.Info("verifying url")

		// read the slack api's url challenge
		r := new(slackevents.ChallengeResponse)

		err := json.Unmarshal([]byte(body), &r)
		if err != nil {
			err = errors.Wrap(err, "could not unmarshal url challenge")
			return false, err
		}

		// echo the challenge
		c.JSON(http.StatusOK, r.Challenge)
		return true, nil
	}

	return false, nil
}
