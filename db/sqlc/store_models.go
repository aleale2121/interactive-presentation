package db


type OptionRow struct {
	Optionkey   string `json:"optionkey" db:"optionkey"`
	Optionvalue string `json:"optionvalue" db:"optionvalue"`
}

type CurrPollIndexResult struct {
	ID       string
	Question string
	Options  []OptionRow `json:"options"`
}

type VoteParams struct {
	PresentationID string
	Pollid         string `db:"pollid"`
	Optionkey      string `db:"optionkey"`
	Clientid       string `db:"clientid"`
}
