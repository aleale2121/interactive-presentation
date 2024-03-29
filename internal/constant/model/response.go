package model

import (
	"encoding/json"
)

// PresentionResponseDTO represents the request model for creating a presentation.
type PresentionResponseDTO struct {
	CurrentPollIndex int32                      `json:"current_poll_index"`
	Polls            []PresentationResponsePoll `json:"polls"`
}

// VoteResponse represents the response model for a vote.
type VoteResponse struct {
	Key      string `json:"key"`
	ClientId string `json:"client_id"`
}

type PresentationResponsePoll struct {
	Question string          `json:"question" binding:"required"`
	Options  json.RawMessage `json:"options" binding:"required,dive"`
}

type CurrentPoll struct {
	ID string `json:"id"`
	Poll
}
