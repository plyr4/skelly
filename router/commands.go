package router

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/davidvader/skelly/skelly"
	"github.com/davidvader/skelly/util"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

// commandsHandler represents the API handler for executing slack slash commands
func commandsHandler(c *gin.Context) {

	// extract the request from the gin context
	r := c.Request

	// retrieve the slack secrets from the environment
	sSecret := os.Getenv("SKELLY_SIGNING_SECRET")

	// set up signing secret verification
	verifier, err := slack.NewSecretsVerifier(r.Header, sSecret)
	if err != nil {
		err = errors.Wrap(err, "could not create signing secret verifier")
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	// update the request body with the verifier
	r.Body = ioutil.NopCloser(io.TeeReader(r.Body, &verifier))

	// parse the slash command from the request
	// this also parses the incoming verification secret
	s, err := slack.SlashCommandParse(r)
	if err != nil {
		err = errors.Wrap(err, "could not parse slash command from request")
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	// verify signing secret
	err = verifier.Ensure()
	if err != nil {
		err = errors.Wrap(err, "could not ensure signing secret")
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	// execute async to allow http connection to close
	go func() {

		// handle the command
		err = skelly.HandleSlashCommand(&s)
		if err != nil {
			err = errors.Wrap(err, "could not execute slash command")
			logrus.Error(err)
			return
		}
	}()

	// acknowledge request
	util.RespondOK(c)
}
