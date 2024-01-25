// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: polls_vote_query.sql

package persistence

import (
	"context"

	"github.com/google/uuid"
)

const createVote = `-- name: CreateVote :exec
INSERT INTO votes
  (pollid, optionkey, clientid)
VALUES
  ($1, $2, $3)
`

type CreateVoteParams struct {
	Pollid    uuid.UUID `db:"pollid" json:"pollid"`
	Optionkey string    `db:"optionkey" json:"optionkey"`
	Clientid  string    `db:"clientid" json:"clientid"`
}

func (q *Queries) CreateVote(ctx context.Context, arg CreateVoteParams) error {
	_, err := q.db.ExecContext(ctx, createVote, arg.Pollid, arg.Optionkey, arg.Clientid)
	return err
}

const getVotes = `-- name: GetVotes :many
SELECT
  votes.optionkey as "key",
  votes.clientid as "client_id"
FROM votes
  JOIN polls ON votes.pollid = polls.id
  JOIN presentations ON polls.presentationID = presentations.id
WHERE presentations.id = $1
  AND polls.id = $2
ORDER BY votes.optionkey
`

type GetVotesParams struct {
	PresentationID uuid.UUID `db:"presentation_id" json:"presentation_id"`
	PollID         uuid.UUID `db:"poll_id" json:"poll_id"`
}

type GetVotesRow struct {
	Key      string `db:"key" json:"key"`
	ClientID string `db:"client_id" json:"client_id"`
}

func (q *Queries) GetVotes(ctx context.Context, arg GetVotesParams) ([]GetVotesRow, error) {
	rows, err := q.db.QueryContext(ctx, getVotes, arg.PresentationID, arg.PollID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetVotesRow{}
	for rows.Next() {
		var i GetVotesRow
		if err := rows.Scan(&i.Key, &i.ClientID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
