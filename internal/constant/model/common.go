package model

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Poll struct {
	ID       uuid.UUID       `db:"id" json:"id" binding:"required"`
	Question string          `db:"question" json:"question" binding:"required"`
	Options  json.RawMessage `db:"options"`
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

type PresentationPoll struct {
	Question string   `json:"question" binding:"required"`
	Options  []Option `json:"options" binding:"required,dive"`
}
