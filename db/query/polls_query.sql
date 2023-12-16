-- name: CreatePoll :one
INSERT INTO polls
  (presentationid, question,pollindex)
VALUES
  ($1, $2, $3)
RETURNING *;

-- name: GetPoll :one
SELECT *
FROM polls
WHERE id = $1;

-- name: GetPollByPID :one
SELECT *
FROM polls
WHERE id = $1 and presentationid = $2;


-- name: GetPresentationPolls :many
SELECT id, question
FROM polls
WHERE presentationid = $1
ORDER BY createdat ASC;