package db

import "gopkg.in/mgo.v2/bson"

// channelSelector return mgo/bson selector for retrieving reactions by channel
func channelSelector(channel string) bson.M {

	// returns mgo/bson selector containing channel
	return bson.M{
		"channel": channel,
	}
}

// reactionsSelector return mgo/bson selector for retrieving reactions by channel
func reactionsSelector(channel string) []bson.DocElem {

	// returns mgo/bson selector containing channel
	// uses bson.D with channel first for performance
	return bson.D{
		{
			Name:  "channel",
			Value: channel,
		},
	}
}

// reactionSelector return mgo/bson selector for retrieving reactions by channel
func reactionSelector(channel string) []bson.DocElem {

	// returns mgo/bson selector containing channel
	// uses bson.D with channel first for performance
	return bson.D{
		{
			Name:  "channel",
			Value: channel,
		},
	}
}

// responseSelector return mgo/bson selector for retreiving responses by channel/user/timestamp
func responseSelector(channel, user, timestamp string) []bson.DocElem {

	// return mgo/bson selector containing channel, user, and timestamp
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
			Name:  "user",
			Value: user,
		},
	}
}
