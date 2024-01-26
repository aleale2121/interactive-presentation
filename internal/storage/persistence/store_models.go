package persistence

import "github.com/google/uuid"

type VoteParams struct {
	PresentationID uuid.UUID
	Pollid         uuid.UUID `db:"pollid"`
	Optionkey      string    `db:"optionkey"`
	Clientid       string    `db:"clientid"`
}
