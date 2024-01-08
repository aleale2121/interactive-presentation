package models

import "github.com/google/uuid"

// CreateVoteRequest represents the request model for creating a vote.
type CreateVoteRequest struct {
	PollID   uuid.UUID `json:"poll_id" binding:"required"`
	Key      string    `json:"key" binding:"required"`
	ClientId string    `json:"client_id" binding:"required"`
}

// CreatePresentationRequest represents the request model for creating a presentation.
type CreatePresentionRequest struct {
	Polls []PresentationPoll `json:"polls"`
}
