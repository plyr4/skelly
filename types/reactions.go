package types

// Reaction is the struct representation for skelly reactions
type Reaction struct {
	Channel       string `json:"channel"`
	Emoji         string `json:"emoji"`
	UserGroup     string `json:"usergroup"`
	UserGroupFull string `json:"usergroup_full"`
	Response      string `json:"response"`
}

// Response is the struct represtation for a stored response
type Response struct {
	Channel   string `json:"channel"`
	Emoji     string `json:"emoji"`
	Timestamp string `json:"timestamp"`
}
