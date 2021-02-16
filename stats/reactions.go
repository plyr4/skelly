package stats

import (
	"github.com/davidvader/skelly/db"
	"github.com/davidvader/skelly/types"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

// GetChannelStats retrieve stats for a channel using data from the db
func GetChannelStats(channel string) (*types.ChannelStats, error) {

	logrus.Infof("getting channel stats for channel(%s)", channel)

	// retrieve reactions from the database
	reactions, err := db.GetChannelReactions(channel)
	if err != nil {
		return nil, err
	}

	// build channel stats
	stats := types.ChannelStats{
		TotalRules: len(*reactions),
	}

	return &stats, nil
}

// GetSkellyStats retrieve stats for all of Skelly
func GetSkellyStats(bToken string) (*types.SkellyStats, error) {

	logrus.Infof("getting skelly stats")

	// retrieve channels from the database
	channels, err := db.GetChannels()
	if err != nil {
		return nil, err
	}

	logrus.Infof("total channels: (%v)", len(*channels))

	// create an api client
	api := slack.New(bToken)

	channelMap := map[string]*slack.Channel{}
	for channelID := range *channels {

		logrus.Tracef("retrieving channel info for id(%v)", channelID)

		// get channel name from the slack api
		convo, err := api.GetConversationInfo(channelID, false)

		// skip direct message groups and bad channels
		if err != nil {
			continue
		}
		channelMap[channelID] = convo
	}

	logrus.Tracef("retrieved channels: %v", channelMap)

	// build skelly stats
	stats := types.SkellyStats{
		TotalChannels: len(*channels),
		Channels:      channelMap,
	}

	return &stats, nil
}
