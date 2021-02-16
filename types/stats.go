package types

import "github.com/slack-go/slack"

// ChannelStats is the struct representation for skelly channel statistics.
type ChannelStats struct {
	TotalRules int
}

// SkellyStats is the struct representation for skelly statistics.
type SkellyStats struct {
	TotalChannels int
	Channels      map[string]*slack.Channel
}
