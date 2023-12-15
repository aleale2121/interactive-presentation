package models

import "database/sql"

// CreatePresentationRequest represents the request model for creating a presentation.
type CreatePresentationRequest struct {
	UserID        string         `json:"user_id" binding:"required"`
	Title         string         `json:"title" binding:"required"`
	Description   string         `json:"description"`
	CurrentPollID sql.NullString `json:"current_poll_id"`
}

// CreatePollRequest represents the request model for creating a poll.
type CreatePollRequest struct {
	PresentationID string `json:"presentation_id" binding:"required"`
	Question       string `json:"question" binding:"required"`
}

// CreateOptionRequest represents the request model for creating an option.
type CreateOptionRequest struct {
	PollID      string `json:"poll_id" binding:"required"`
	OptionKey   string `json:"option_key" binding:"required"`
	OptionValue string `json:"option_value" binding:"required"`
}

// CreateUserRequest represents the request model for creating a user.
type CreateUserRequest struct {
	Name string `json:"name" binding:"required"`
}

// CreateVoteRequest represents the request model for creating a vote.
type CreateVoteRequest struct {
	PollID   string `json:"poll_id" binding:"required"`
	Key      string `json:"key" binding:"required"`
	ClientId string `json:"client_id" binding:"required"`
}

// PresentationResponse represents the response model for a presentation.
type PresentationResponse struct {
	ID            string         `json:"id"`
	UserID        string         `json:"user_id"`
	Title         string         `json:"title"`
	CurrentPollID sql.NullString `json:"current_poll_id"`
	Description   sql.NullString `json:"description"`
}

// PollResponse represents the response model for a poll.
type PollResponse struct {
	ID       string           `json:"poll_id"`
	Question string           `json:"question"`
	Options  []OptionResponse `json:"options"`
}

// OptionResponse represents the response model for an option.
type OptionResponse struct {
	OptionKey   string `json:"option_key"`
	OptionValue string `json:"option_value"`
}

// UserResponse represents the response model for a user.
type UserResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// VoteResponse represents the response model for a vote.
type VoteResponse struct {
	Key      string `json:"key"`
	ClientId string `json:"client_id"`
}

type Poll struct {
	Question string   `json:"question" binding:"required"`
	Options  []Option `json:"options" binding:"required,dive"`
}

type Option struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value" binding:"required"`
}

type Presention struct {
	Polls []Poll `json:"polls"`
}