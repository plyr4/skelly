package router

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/davidvader/skelly/skelly"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

// eventsHandler represents the API handler for handling slack events
func eventsHandler(c *gin.Context) {

	// extract the request from the gin context
	r := c.Request

	// retrieve the slack secrets from the environment
	sSecret := os.Getenv("SKELLY_SIGNING_SECRET")

	// set up a signing secret verifier
	verifier, err := slack.NewSecretsVerifier(r.Header, sSecret)
	if err != nil {
		err = errors.Wrap(err, "could not create signing secret verifier")
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
		return
	}

	// update the request body with the verifier
	r.Body = ioutil.NopCloser(io.TeeReader(r.Body, &verifier))

	// read request body
	b, err := c.GetRawData()
	if err != nil {
		err = errors.Wrap(err, "could not read body from request")
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	// retrieve the slack secrets from the environment
	vToken := os.Getenv("SKELLY_VERIFICATION_TOKEN")

	// set up token comparator for use with slackevents
	tokenComparator := &slackevents.TokenComparator{
		VerificationToken: vToken,
	}

	// verify the token and parse the event from the incoming request
	e, err := slackevents.ParseEvent(json.RawMessage(b),
		slackevents.OptionVerifyToken(tokenComparator))
	if err != nil {
		err = errors.Wrap(err, "could not parse event from request")
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	// retrieve the slack secrets from the environment
	bToken := os.Getenv("SKELLY_BOT_TOKEN")

	// handle the event
	err = skelly.HandleEvent(c, b, &e, bToken)
	if err != nil {
		err = errors.Wrap(err, "could not execute slash command")
		logrus.Error(err)
		return
	}

}
