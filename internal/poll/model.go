package poll

import "time"

// Poll represents a poll/question.
type Poll struct {
	ID        int64      `json:"id"`
	Question  string     `json:"question"`
	Options   []Option   `json:"options,omitempty"`
	UserID    int64      `json:"user_id"`
	CreatedAt time.Time  `json:"created_at"`
}

// Option represents an option/answer for a poll.
type Option struct {
	ID      int64  `json:"id"`
	PollID  int64  `json:"poll_id"`
	Text    string `json:"text"`
	Votes   int64  `json:"votes,omitempty"` 
}

// Vote represents a vote cast by a user for a particular option in a poll.
type Vote struct {
	ID        int64     `json:"id"`
	PollID    int64     `json:"poll_id"`
	OptionID  int64     `json:"option_id"`
	UserID    int64     `json:"user_id"` 
	CreatedAt time.Time `json:"created_at"`
}