package models

// CreateVoteRequest represents the request model for creating a vote.
type CreateVoteRequest struct {
	PollID   string `json:"poll_id" binding:"required"`
	Key      string `json:"key" binding:"required"`
	ClientId string `json:"client_id" binding:"required"`
}

// VoteResponse represents the response model for a vote.
type VoteResponse struct {
	Key      string `json:"key"`
	ClientId string `json:"client_id"`
}

// CreatePresentationRequest represents the request model for creating a presentation.
type CreatePresentionRequest struct {
	Polls []Poll `json:"polls"`
}

type Poll struct {
	Question string   `json:"question" binding:"required"`
	Options  []Option `json:"options" binding:"required,dive"`
}

type Option struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value" binding:"required"`
}
