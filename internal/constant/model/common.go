package model

import (
	"github.com/google/uuid"
)

type Poll struct {
	Question string   `json:"question" binding:"required"`
	Options  []Option `json:"options" binding:"required,dive"`
}

type Option struct {
	Key   string `json:"key" db:"optionKey" binding:"required"`
	Value string `json:"value" db:"optionValue" binding:"required"`
}

type Vote struct {
	PollID         uuid.UUID
	PresentationID uuid.UUID
	Options        []string
	VoteCount      int
}

