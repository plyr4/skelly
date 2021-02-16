package db

import "gopkg.in/mgo.v2/bson"

// channelSelector return mgo/bson selector for retrieving reactions by channel
func channelSelector(channel string) bson.M {

	// returns mgo/bson selector containing channel
	return bson.M{
		"channel": channel,
	}
}

// reactionsSelector return mgo/bson selector for retrieving reactions by channel/emoji
func reactionsSelector(channel, emoji string) []bson.DocElem {

	// returns mgo/bson selector containing channel and emoji
	// uses bson.D with channel first for performance
	return bson.D{
		{
			Name:  "channel",
			Value: channel,
		},
		{
			Name:  "emoji",
			Value: emoji,
		},
	}
}

// reactionSelector return mgo/bson selector for retrieving reactions by channel/emoji/usergroup
func reactionSelector(channel, emoji, usergroup string) []bson.DocElem {

	// returns mgo/bson selector containing channel and emoji
	// uses bson.D with channel first for performance
	return bson.D{
		{
			Name:  "channel",
			Value: channel,
		},
		{
			Name:  "emoji",
			Value: emoji,
		},
		{
			Name:  "usergroup",
			Value: usergroup,
		},
	}
}

// responseSelector return mgo/bson selector for retreiving responses by channel/emoji/timestamp
func responseSelector(channel, emoji, timestamp string) []bson.DocElem {

	// return mgo/bson selector containing channel, emoji, and timestamp
	return bson.D{
		{
			Name:  "timestamp",
			Value: timestamp,
		},
		{
			Name:  "channel",
			Value: channel,
		},
		{
			Name:  "emoji",
			Value: emoji,
		},
	}
}
