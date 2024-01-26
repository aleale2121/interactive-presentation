package model

import "github.com/google/uuid"

// CreateVoteRequest represents the request model for creating a vote.
type CreateVoteRequestDTO struct {
	PollID   uuid.UUID `json:"poll_id" binding:"required"`
	Key      string    `json:"key" binding:"required"`
	ClientId string    `json:"client_id" binding:"required"`
}

// CreatePresentationRequest represents the request model for creating a presentation.
type CreatePresentionRequestDTO struct {
	Polls []PresentationPoll `json:"polls"`
}
