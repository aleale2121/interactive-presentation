package db

import "database/sql"

type PollIndexParams struct {
	Currentpollindex sql.NullInt64 `db:"currentpollindex"`
	ID               string        `db:"id"`
}

type CurrPollIndexResult struct {
	ID       string
	Question string
	Options  []GetPollOptionsRow
}

type VoteParams struct {
	PresentationID string
	Pollid         string `db:"pollid"`
	Optionkey      string `db:"optionkey"`
	Clientid       string `db:"clientid"`
}
