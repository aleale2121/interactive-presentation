// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: presentation_query.sql

package persistence

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/google/uuid"
)

const createPresentationAndPolls = `-- name: CreatePresentationAndPolls :one
WITH presentations_cte AS (
INSERT INTO presentations
  (id, currentpollindex)
VALUES
  (uuid_generate_v4(), 0)
RETURNING id
),
data_cte AS
(
  SELECT
  to_jsonb($1
::jsonb) AS data
),
polls_cte AS
(
INSERT INTO polls
  (question, presentationid, pollindex)
SELECT
  (arr.elem ->> 'question')
::TEXT AS question,
    pc.id AS presentationid,
    arr.idx AS pollindex
  FROM presentations_cte pc, data_cte dc,
    jsonb_array_elements
(dc.data)
WITH ORDINALITY AS arr
(elem, idx)
  RETURNING id, pollindex
),
options_cte AS
(
INSERT INTO options
  (pollid, optionkey, optionvalue)
SELECT
  pc.id,
  (o ->> 'key')
::TEXT AS optionkey,
(o ->> 'value')::TEXT AS optionvalue
  FROM polls_cte pc, data_cte dc,
    LATERAL jsonb_array_elements
(dc.data -> pc.pollindex - 1 -> 'options') AS o
  RETURNING id,pollid
)

SELECT id
FROM presentations_cte
`

func (q *Queries) CreatePresentationAndPolls(ctx context.Context, polls json.RawMessage) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, createPresentationAndPolls, polls)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const getNextPoll = `-- name: GetNextPoll :one
WITH updated_polls_cte AS(
UPDATE presentations
	SET currentpollindex = currentpollindex + 1
	WHERE presentations.id = $1 AND currentpollindex + 1 <= (
    SELECT count(*)
  FROM polls
  WHERE presentationid = $1
  )
RETURNING id, currentpollindex
)
SELECT
  p.id AS id,
  p.question AS question,
  (
    SELECT jsonb_agg(jsonb_build_object(
      'key', o.optionkey,
      'value', o.optionvalue
    )) AS options
  FROM options o
    INNER JOIN polls p ON o.pollid = p.id
  WHERE p.presentationid = upc.id AND p.pollindex=upc.currentpollindex
  ) AS options
FROM polls p
  INNER JOIN updated_polls_cte upc ON p.presentationid = upc.id AND p.pollindex=upc.currentpollindex
WHERE upc.id IS NOT NULL
`

type GetNextPollRow struct {
	ID       uuid.UUID       `db:"id" json:"id"`
	Question string          `db:"question" json:"question"`
	Options  json.RawMessage `db:"options" json:"options"`
}

func (q *Queries) GetNextPoll(ctx context.Context, id uuid.UUID) (GetNextPollRow, error) {
	row := q.db.QueryRowContext(ctx, getNextPoll, id)
	var i GetNextPollRow
	err := row.Scan(&i.ID, &i.Question, &i.Options)
	return i, err
}

const getPoll = `-- name: GetPoll :one
SELECT
  p.id AS id,
  p.question AS question,
  (
    SELECT jsonb_agg(jsonb_build_object(
      'key', o.optionkey,
      'value', o.optionvalue
    )) AS options
  FROM options o
  WHERE o.pollid = p.id
  ) AS options
FROM polls p
WHERE p.presentationid = $1 and p.pollindex=(SELECT currentpollindex
  FROM presentations
  WHERE id=$1)
`

type GetPollRow struct {
	ID       uuid.UUID       `db:"id" json:"id"`
	Question string          `db:"question" json:"question"`
	Options  json.RawMessage `db:"options" json:"options"`
}

func (q *Queries) GetPoll(ctx context.Context, presentationid uuid.UUID) (GetPollRow, error) {
	row := q.db.QueryRowContext(ctx, getPoll, presentationid)
	var i GetPollRow
	err := row.Scan(&i.ID, &i.Question, &i.Options)
	return i, err
}

const getPresentation = `-- name: GetPresentation :one
SELECT id, currentpollindex
FROM presentations
WHERE id = $1
`

func (q *Queries) GetPresentation(ctx context.Context, id uuid.UUID) (Presentation, error) {
	row := q.db.QueryRowContext(ctx, getPresentation, id)
	var i Presentation
	err := row.Scan(&i.ID, &i.Currentpollindex)
	return i, err
}

const getPresentationAndPoll = `-- name: GetPresentationAndPoll :one
SELECT 
  pr.id as presentation_id,
  pr.currentpollindex,
  pl.id as poll_id,
  pl.question as question,
  pl.pollindex
FROM
  presentations as pr JOIN polls as pl ON pr.id=pl.presentationid
WHERE 
  pr.id=$1 AND pl.id=$2
`

type GetPresentationAndPollParams struct {
	PresentationID uuid.UUID `db:"presentation_id" json:"presentation_id"`
	PollID         uuid.UUID `db:"poll_id" json:"poll_id"`
}

type GetPresentationAndPollRow struct {
	PresentationID   uuid.UUID     `db:"presentation_id" json:"presentation_id"`
	Currentpollindex sql.NullInt32 `db:"currentpollindex" json:"currentpollindex"`
	PollID           uuid.UUID     `db:"poll_id" json:"poll_id"`
	Question         string        `db:"question" json:"question"`
	Pollindex        int32         `db:"pollindex" json:"pollindex"`
}

func (q *Queries) GetPresentationAndPoll(ctx context.Context, arg GetPresentationAndPollParams) (GetPresentationAndPollRow, error) {
	row := q.db.QueryRowContext(ctx, getPresentationAndPoll, arg.PresentationID, arg.PollID)
	var i GetPresentationAndPollRow
	err := row.Scan(
		&i.PresentationID,
		&i.Currentpollindex,
		&i.PollID,
		&i.Question,
		&i.Pollindex,
	)
	return i, err
}

const getPreviousPoll = `-- name: GetPreviousPoll :one
WITH updated_polls_cte AS(
UPDATE presentations
	SET currentpollindex = currentpollindex - 1
	WHERE presentations.id = $1 AND currentpollindex - 1 >=0
RETURNING id, currentpollindex
)
SELECT
  p.id AS id,
  p.question AS question,
  (
    SELECT jsonb_agg(jsonb_build_object(
      'key', o.optionkey,
      'value', o.optionvalue
    )) AS options
  FROM options o
    INNER JOIN polls p ON o.pollid = p.id
  WHERE p.presentationid = upc.id AND p.pollindex=upc.currentpollindex
  ) AS options
FROM polls p
  INNER JOIN updated_polls_cte upc ON p.presentationid = upc.id AND p.pollindex=upc.currentpollindex
WHERE upc.id IS NOT NULL
`

type GetPreviousPollRow struct {
	ID       uuid.UUID       `db:"id" json:"id"`
	Question string          `db:"question" json:"question"`
	Options  json.RawMessage `db:"options" json:"options"`
}

func (q *Queries) GetPreviousPoll(ctx context.Context, id uuid.UUID) (GetPreviousPollRow, error) {
	row := q.db.QueryRowContext(ctx, getPreviousPoll, id)
	var i GetPreviousPollRow
	err := row.Scan(&i.ID, &i.Question, &i.Options)
	return i, err
}

const listPolls = `-- name: ListPolls :many
SELECT
  p.id AS id,
  p.question,
  jsonb_agg(jsonb_build_object(
        'key', o.optionKey,
        'value', o.optionValue
    )) AS options
FROM
  polls p
  JOIN
  options o ON p.id = o.pollID
WHERE
    p.presentationID = $1
GROUP BY
    p.id, p.question, p.pollindex, p.createdAt
ORDER BY
    p.pollindex
`

type ListPollsRow struct {
	ID       uuid.UUID       `db:"id" json:"id"`
	Question string          `db:"question" json:"question"`
	Options  json.RawMessage `db:"options" json:"options"`
}

func (q *Queries) ListPolls(ctx context.Context, presentationid uuid.UUID) ([]ListPollsRow, error) {
	rows, err := q.db.QueryContext(ctx, listPolls, presentationid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListPollsRow{}
	for rows.Next() {
		var i ListPollsRow
		if err := rows.Scan(&i.ID, &i.Question, &i.Options); err != nil {
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
