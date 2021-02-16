package skelly

import (
	"fmt"
	"strings"

	"github.com/davidvader/skelly/db"
	"github.com/davidvader/skelly/emojis"
	"github.com/davidvader/skelly/types"
	"github.com/davidvader/skelly/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

// React takes channel and emoji and reacts with the appropriate response based on application configuration.
func React(bToken, channel, emoji, user, ts string) error {

	// retrieve the emoji's core shortname
	shortname, err := emojis.GetShortname(emoji)
	if err != nil {
		err = errors.Wrap(err, "could not get emoji shortname")
		return err
	}

	// set emoji to core shortname
	emoji = shortname

	// retrieve all of the reactions for the channel/emoji
	reactions, err := db.GetReactions(channel, emoji)
	if err != nil {
		err = errors.Wrap(err, "could not get reaction from db")
		return err
	}

	logrus.Infof("retreived (%v) reactions for channel(%s) emoji(%s)", len(reactions), channel, emoji)

	// create an api client
	api := slack.New(bToken)

	// fetch all usergroups containing users
	logrus.Info("retrieving usergroups")

	usergroups, err := api.GetUserGroups(slack.GetUserGroupsOptionIncludeUsers(true))

	_usergroups := []slack.UserGroup{}

	// determine all usergroups this user is a part of
	for _, ug := range usergroups {
		for _, u := range ug.Users {
			if user == u {
				_usergroups = append(_usergroups, ug)
				break
			}
		}
	}

	logrus.Infof("filtered (%v) usergroups for user(%s)", len(_usergroups), user)

	_reactions := []*types.Reaction{}

	// filter the reactions based on user id and reaction usergroup
	logrus.Infof("filtering reactions for channel(%s) emoji(%s) user(%s)", channel, emoji, user)

	for _, r := range reactions {
		if r.UserGroup == "none" {
			_reactions = append(_reactions, r)
			continue
		}

		for _, ug := range _usergroups {
			if strings.ToLower(r.UserGroup) == strings.ToLower(ug.ID) {
				_reactions = append(_reactions, r)
				break
			}
		}
	}

	// use conversations.history to get the correct thread timestamp
	// if ts is empty, post as a new message (debug)
	if ts != "none" {
		ts, err = util.GetThreadTimestamp(bToken, channel, ts)
		if err != nil {
			err = errors.Wrap(err, fmt.Sprintf("could not get thread timestamp for ts(%s)", ts))
			return err
		}
	}

	logrus.Infof("reacting to (%v) reactions for channel(%s) emoji(%s)", len(_reactions), channel, emoji)

	// respond to possibly multiple reactions
	for _, r := range _reactions {

		// do not react if response is empty
		if len(r.Response) == 0 {
			continue
		}

		// check database for existing response
		exists, err := db.CheckResponse(channel, emoji, ts)
		if err != nil {
			err = errors.Wrap(err, "could not check for existing response")
			return err
		}

		// do not react if response already exists
		if exists {
			logrus.Infof("skipping, reaction response exists for channel(%s) emoji(%s) ts(%s)", channel, emoji, ts)
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
		logrus.Infof("posting reaction for channel(%s) emoji(%s) ts(%s)", channel, emoji, ts)

		_, mts, err := api.PostMessage(channel, options...)
		if err != nil {
			err = errors.Wrap(err, "could not post response")
			return err
		}

		err = db.StoreResponse(channel, emoji, ts)
		if err != nil {
			err = errors.Wrap(err, "response posted, but could not remember the response")
			return err
		}

		logrus.Infof("reaction posted for channel(%s) emoji(%s) ts(%s) at msg_ts(%s)", channel, emoji, ts, mts)
	}
	return nil
}
