package db

type PollOption struct {
	Optionkey   string `json:"optionkey" db:"optionkey"`
	Optionvalue string `json:"optionvalue" db:"optionvalue"`
}

type CurrentPoll struct {
	ID       string       `json:"id"`
	Question string       `json:"question"`
	Options  []PollOption `json:"options"`
}

type VoteParams struct {
	PresentationID string
	Pollid         string `db:"pollid"`
	Optionkey      string `db:"optionkey"`
	Clientid       string `db:"clientid"`
}
