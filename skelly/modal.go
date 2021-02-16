package skelly

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/davidvader/skelly/emojis"
	"github.com/pkg/errors"
	"github.com/slack-go/slack"
)

// modal builds the default view modal for managing a reaction
func modal(callback, header, metadata, emoji, usergroup, response string) slack.ModalViewRequest {

	// header section
	headerText := slack.NewTextBlockObject("mrkdwn", header+" A reaction will trigger a response when an emoji is added to a message by a member of the specified user group.", false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	// user args
	ugText := usergroup
	if usergroup == "none" {
		ugText += " (all users)"
	}

	t := []*slack.TextBlockObject{
		slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Emoji*: %s", emoji), false, false),
		slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Usergroup*: %s", ugText), false, false),
	}
	userInputSection := slack.NewSectionBlock(nil, t, nil)

	// response input
	responseText := slack.NewTextBlockObject("plain_text", "Response", false, false)
	responsePlaceholder := slack.NewTextBlockObject("plain_text", "Enter a response to the emoji", false, false)

	responseElement := slack.NewPlainTextInputBlockElement(responsePlaceholder, "response")
	responseElement.Multiline = true
	responseElement.InitialValue = response

	responseInput := slack.NewInputBlock("Response", responseText, responseElement)

	// build message from blocks
	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			headerSection,
			userInputSection,
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
func deleteModal(callback, metadata, emoji, usergroup string) slack.ModalViewRequest {

	// header section
	headerText := slack.NewTextBlockObject("mrkdwn", "Delete a reaction.", false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	// user args
	ugText := usergroup
	if usergroup == "none" {
		ugText += " (all users)"
	}

	e := slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Emoji*: %s", emoji), false, false)
	ug := slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Usergroup*: %s", ugText), false, false)

	emojiSection := slack.NewSectionBlock(nil, []*slack.TextBlockObject{e}, nil)
	usergroupSection := slack.NewSectionBlock(nil, []*slack.TextBlockObject{ug}, nil)

	// build message from blocks
	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			headerSection,
			emojiSection,
			usergroupSection,
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
func parseViewMetadata(view *slack.View) (string, string, string, error) {

	// split metadata by delimiter :
	metadata := strings.Split(view.PrivateMetadata, " ")
	if len(metadata) < 4 {
		err := fmt.Errorf("invalid view submission metadata(%v)", metadata)
		return "", "", "", err
	}

	// extract args
	emoji := metadata[1]
	if len(emoji) == 0 {
		return "", "", "", errors.New("bad metadata no emoji")
	}

	usergroup := metadata[2]
	if len(usergroup) == 0 {
		return "", "", "", errors.New("bad metadata no usergroup")
	}

	channel := metadata[3]
	if len(channel) == 0 {
		return "", "", "", errors.New("bad metadata no channel")
	}

	return channel, emoji, usergroup, nil
}

// parseUserGroup takes usergroup in format <!subteam^ID> and parses information
func parseUserGroup(usergroup string) (string, string, error) {

	// check for none
	if usergroup == "none" {
		return "none", "", nil
	}

	// compile usergroup regexp
	r, err := regexp.Compile(`\<\!subteam\^(.*)\|\@(.*)\>`)
	if err != nil {
		return "", "", errors.Wrap(err, "could not compile regexp")
	}

	// find substrings
	match := r.FindStringSubmatch(usergroup)

	// if regexp did not find enough substrings
	if len(match) < 3 {
		return "", "", fmt.Errorf("could not match regex for usergroup(%s)", usergroup)
	}

	id := match[1]

	handle := match[2]

	return id, handle, nil
}

// parseEmoji takes emoji in format :smile: and parses ID (smile)
func parseEmoji(emoji string) (string, error) {

	// compile emoji regexp
	r, err := regexp.Compile(`:(.*?):`)
	if err != nil {
		return "", errors.Wrap(err, "could not compile regexp")
	}

	// find substrings
	match := r.FindStringSubmatch(emoji)

	// if regexp did not find enough substrings
	if len(match) < 2 {
		return "", fmt.Errorf("could not match regex for emoji(%s)", emoji)
	}

	id := match[1]

	// retrieve the emoji's core shortname
	shortname, err := emojis.GetShortname(id)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("could not get shortname(%s)", id))
	}

	return shortname, nil
}
