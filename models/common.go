package models

import "github.com/google/uuid"

type Poll struct {
	Question string   `json:"question" binding:"required"`
	Options  []Option `json:"options" binding:"required,dive"`
}

type Option struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value" binding:"required"`
}

type Vote struct {
	PollID         uuid.UUID
	PresentationID uuid.UUID
	Options        []string
	VoteCount      int
}
