package models

// PresentionResponse represents the request model for creating a presentation.
type PresentionResponse struct {
	CurrentPollIndex int32  `json:"current_poll_index"`
	Polls            []Poll `json:"polls"`
}

// VoteResponse represents the response model for a vote.
type VoteResponse struct {
	Key      string `json:"key"`
	ClientId string `json:"client_id"`
}
