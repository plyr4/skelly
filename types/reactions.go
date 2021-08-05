package types

// Reaction is the struct representation for skelly reactions
type Reaction struct {
	Channel   string `json:"channel"`
	CreatedBy string `json:"user"`
	Response  string `json:"response"`
}

// Response is the struct represtation for a stored response
type Response struct {
	Channel   string `json:"channel"`
	User      string `json:"user"`
	Timestamp string `json:"timestamp"`
}
