package skelly

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/davidvader/skelly/util"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

// HandleInteraction takes request body and interaction callback and executes the appropriate interaction
func HandleInteraction(c *gin.Context, body string) error {

	// parse the main interaction callback
	callback, err := parseInteraction(body)
	if err != nil {
		err = errors.Wrap(err, "could not parse interaction")
		return err
	}

	// execute interaction
	switch callback.Type {

	case slack.InteractionTypeViewSubmission:

		// execute async to allow http connection to close
		go func() {

			// handle the view submission
			err := handleViewSubmission(&callback.View, callback.User.ID, callback.ResponseURL)
			if err != nil {
				err = errors.Wrap(err, "could not handle submission")
				logrus.Error(err)
				return
			}
		}()

		// acknowledge the submission
		util.RespondOK(c)

		return nil

	default:
		err := fmt.Errorf("unsupported interaction type(%s)", callback.Type)
		return err
	}
}

// parseInteraction takes request body and parses it into an interaction callback
func parseInteraction(body string) (*slack.InteractionCallback, error) {

	// unmarshal into struct
	var cb slack.InteractionCallback
	if err := json.Unmarshal([]byte(body), &cb); err != nil {
		err = errors.Wrap(err, "could not unmarshal callback payload")
		return nil, err
	}
	return &cb, nil
}

// handleViewSubmission takes slack view, extracts callback id, and executes a view submission
func handleViewSubmission(view *slack.View, user, responseURL string) error {

	// extract callbackID
	callbackID := strings.Split(view.CallbackID, ":")
	if len(callbackID) == 0 {
		err := fmt.Errorf("invalid callback id(%s)", callbackID)
		return err
	}

	// execute submission action
	switch callbackID[0] {
	case addSubCommand:

		// handle view submission for /skelly add
		err := handleAddSubmission(view, user, responseURL)
		if err != nil {
			err = errors.Wrap(err, "could not handle add submission")
			return err
		}

		return nil

	case updateSubCommand:

		// handle view submission for /skelly update
		err := handleUpdateSubmission(view, user, responseURL)
		if err != nil {
			err = errors.Wrap(err, "could not handle update submission")
			return err
		}

		return nil

	case deleteSubCommand:

		// handle view submission for /skelly delete
		err := handleDeleteSubmission(view, user, responseURL)
		if err != nil {
			err = errors.Wrap(err, "could not handle delete submission")
			return err
		}

		return nil

	default:
		err := fmt.Errorf("unsupported submission action(%s)", callbackID[0])
		return err
	}
}
