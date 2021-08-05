package skelly

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/slack-go/slack"
)

// modal builds the default view modal for managing a reaction
func modal(callback, header, metadata, response string) slack.ModalViewRequest {

	// header section
	headerText := slack.NewTextBlockObject("mrkdwn", header+" A reaction will trigger a response once a day for all users that type in a channel.", false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	// response input
	responseText := slack.NewTextBlockObject("plain_text", "Response", false, false)
	responsePlaceholder := slack.NewTextBlockObject("plain_text", "Enter a response for when users type in this channel", false, false)

	responseElement := slack.NewPlainTextInputBlockElement(responsePlaceholder, "response")
	responseElement.Multiline = true
	responseElement.InitialValue = response

	responseInput := slack.NewInputBlock("Response", responseText, responseElement)

	// build message from blocks
	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			headerSection,
			responseInput,
		},
	}

	// configure modal view request
	request := slack.ModalViewRequest{
		Type:            slack.ViewType("modal"),
		Title:           slack.NewTextBlockObject("plain_text", "Skelly", false, false),
		Close:           slack.NewTextBlockObject("plain_text", "Close", false, false),
		Submit:          slack.NewTextBlockObject("plain_text", "Submit", false, false),
		Blocks:          blocks,
		PrivateMetadata: metadata,
		CallbackID:      callback,
	}

	return request
}

// deleteModal builds the view modal for deleting a reaction
func deleteModal(callback, metadata string) slack.ModalViewRequest {

	// header section
	headerText := slack.NewTextBlockObject("mrkdwn", "Delete a reaction.", false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	// build message from blocks
	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			headerSection,
		},
	}

	// configure modal view request
	request := slack.ModalViewRequest{
		Type:            slack.ViewType("modal"),
		Title:           slack.NewTextBlockObject("plain_text", "Skelly", false, false),
		Close:           slack.NewTextBlockObject("plain_text", "Close", false, false),
		Submit:          slack.NewTextBlockObject("plain_text", "Delete", false, false),
		Blocks:          blocks,
		PrivateMetadata: metadata,
		CallbackID:      callback,
	}

	return request
}

// parseViewResponse takes view and extracts response
func parseViewResponse(view *slack.View) (string, error) {

	// check for valid response state
	_, ok := view.State.Values["Response"]
	if !ok {
		err := errors.New("no Response")
		return "", err
	}

	_, ok = view.State.Values["Response"]["response"]
	if !ok {
		err := errors.New("no Response.response")
		return "", err
	}

	// extract response view state value
	response := view.State.Values["Response"]["response"].Value
	if len(response) == 0 {
		err := errors.New("no Response.response value")
		return "", err
	}

	return response, nil
}

// parseViewMetadata takes view and extracts args from metadata
func parseViewMetadata(view *slack.View) (string, error) {

	// split metadata by delimiter :
	metadata := strings.Split(view.PrivateMetadata, " ")
	if len(metadata) < 4 {
		err := fmt.Errorf("invalid view submission metadata(%v)", metadata)
		return "", err
	}

	// extract args

	channel := metadata[1]
	if len(channel) == 0 {
		return "", errors.New("bad metadata no channel")
	}

	return channel, nil
}
