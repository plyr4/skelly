package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

// InvalidCommand returns a formatted error for improper flag usage
// with a CLI command
func InvalidCommand(f string) error {
	return fmt.Errorf("invalid command: Flag '--%s' is not set or is empty", f)
}

// InvalidFlagValue returns a formatted error for improper flag usage
func InvalidFlagValue(v string, f string) error {
	return fmt.Errorf("invalid value '%s' for flag '--%s'", v, f)
}

// ParsePayload takes request body and parses payload parameter into a usable string
func ParsePayload(b []byte) (string, error) {

	// convert the bytes to a string
	body := string(b)

	// unescape the body
	body, err := url.QueryUnescape(body)
	if err != nil {
		err = errors.Wrap(err, "could not query unescape")
		return body, err
	}

	// parse out the payload parameter
	body = strings.Replace(body, "payload=", "", 1)
	return body, nil
}

// GetThreadTimestamp uses conversations.history to retrieve the timestamp for either the message
// or the thread parent message, if one exists
func GetThreadTimestamp(bToken, channel, ts string) (string, error) {

	logrus.Infof("getting parent timestamp for ts(%s)", ts)

	// configure conversation.history fetch parameters
	params := &slack.GetConversationRepliesParameters{
		ChannelID: channel,
		Timestamp: ts,
	}

	// create a new slack api client
	api := slack.New(bToken)

	// fetch the thread parent, if it exists
	replies, _, _, err := api.GetConversationReplies(params)
	if err != nil {
		err = errors.Wrap(err, "could not fetch conversation history")
		return "", err
	}

	// if the message belongs to a thread, the parent will be the
	// first message in the conversation's history
	if len(replies) > 0 {
		parent := replies[0]
		if len(parent.ThreadTimestamp) > 0 {
			// set the reaction timestamp to the thread's parent message
			ts = parent.ThreadTimestamp
		}
	}

	logrus.Infof("using parent timestamp ts(%s)", ts)

	return ts, nil
}

// SendError responds to a user interaction with an error
func SendError(e, channel, user string) error {

	logrus.Infof("responding with error(%s)", e)

	// build ephemeral response
	text := slack.MsgOptionText(e, false)

	// create default msg options
	options := []slack.MsgOption{
		text,
	}

	// create an api client
	bToken := os.Getenv("SKELLY_BOT_TOKEN")

	api := slack.New(bToken)

	// post message
	_, err := api.PostEphemeral(channel, user, options...)
	if err != nil {
		err = errors.Wrap(err, "could not post response")
		return err
	}
	return nil
}

// Respond takes response text and posts it to the response url
func Respond(responseURL string, response interface{}) error {

	// encode the response
	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(response)

	// respond to the user via the response url
	resp, err := http.Post(responseURL, "application/json; charset=utf-8", buffer)
	if err != nil {
		err = errors.Wrap(err, "could not post response")
	}

	defer resp.Body.Close()

	// read the ressponse body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.Wrap(err, "could not read response body")
		return err
	}

	// check for slack ok
	if string(body) != "ok" {
		return errors.New("slack response not ok")
	}

	return nil
}

// RespondOK responds to a request with Content-Type: application/json and 200 OK
// meant for acknowledging incoming Slack requests
func RespondOK(c *gin.Context) {

	// set to json for slack events api
	c.Header("Content-Type", "application/json")

	// respond with 200
	c.AbortWithStatus(http.StatusOK)
}

// Unique takes string slice and removes unique elements
func Unique(s []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range s {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
