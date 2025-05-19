package poll

import "time"

// Poll represents a poll/question.
type Poll struct {
	ID        int64     `json:"id"`
	Question  string    `json:"question"`
	Options   []Option  `json:"options,omitempty"`
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

// Option represents an option/answer for a poll.
type Option struct {
	ID     int64  `json:"id"`
	PollID int64  `json:"poll_id"`
	Text   string `json:"text"`
	Votes  int64  `json:"votes,omitempty"`
}

// Vote represents a vote cast by a user for a particular option in a poll.
type Vote struct {
	ID        int64     `json:"id"`
	PollID    int64     `json:"poll_id"`
	OptionID  int64     `json:"option_id"`
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

// CreatePollRequest represents the request payload for creating a new poll
type CreatePollRequest struct {
	Question string   `json:"question" example:"What is your favorite programming language?"`
	Options  []string `json:"options" example:"[\"Go\",\"Python\",\"JavaScript\",\"Java\"]"`
}

// CreatePollResponse represents the response for a successfully created poll
type CreatePollResponse struct {
	Poll Poll `json:"poll"`
}

// VotePollRequest represents the request payload for voting on a poll
type VotePollRequest struct {
	OptionID int64 `json:"option_id" example:"1"`
}

// VotePollResponse represents the response for a successfully recorded vote
type VotePollResponse struct {
	Message   string `json:"message" example:"Vote recorded successfully"`
	PollID    int64  `json:"poll_id" example:"1"`
	OptionID  int64  `json:"option_id" example:"2"`
	Timestamp string `json:"timestamp" example:"2025-05-18T10:30:45Z"`
}

// PollResultsResponse represents the response for poll results
type PollResultsResponse struct {
	PollID     int64     `json:"poll_id" example:"1"`
	Question   string    `json:"question" example:"What is your favorite programming language?"`
	TotalVotes int64     `json:"total_votes" example:"42"`
	CreatedAt  time.Time `json:"created_at"`
	Options    []Option  `json:"options"`
}

// OptionResult represents an option with its vote count and percentage
type OptionResult struct {
	ID         int64   `json:"id" example:"1"`
	Text       string  `json:"text" example:"Go"`
	Votes      int64   `json:"votes" example:"25"`
	Percentage float64 `json:"percentage" example:"59.5"`
}
