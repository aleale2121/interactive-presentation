-- name: CreatePresentation :one
INSERT INTO presentations (currentpollindex)
VALUES ($1)
RETURNING id;


-- name: CreatePresentationAndPolls :one
WITH presentations_cte AS (
  INSERT INTO presentations (id, currentpollindex)
  VALUES (uuid_generate_v4(), 0)
  RETURNING id
),
data_cte AS (
  SELECT
    to_jsonb($1::jsonb) AS data
),
polls_cte AS (
  INSERT INTO polls (question, presentationid, pollindex)
  SELECT
    (arr.elem ->> 'question')::TEXT AS question,
    pc.id AS presentationid,
    arr.idx AS pollindex
  FROM presentations_cte pc, data_cte dc,
    jsonb_array_elements(dc.data) WITH ORDINALITY AS arr(elem, idx)
  RETURNING id, pollindex
)
SELECT presentations_cte.id
FROM presentations_cte;


-- name: GetPresentation :one
SELECT *
FROM presentations
WHERE id = $1;

-- name: CreatePoll :one
INSERT INTO polls (presentationid, question,pollindex)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetPoll :one
SELECT *
FROM polls
WHERE id = $1;

-- name: GetPollByPID :one
SELECT *
FROM polls
WHERE id = $1 and presentationid = $2;

-- name: CreateOption :exec
INSERT INTO options (pollid, optionkey, optionvalue)
VALUES ($1, $2, $3);

-- name: GetOptionByKey :one
SELECT *
FROM options
WHERE optionkey = $1;


-- name: GetPresentationPolls :many
SELECT id,question
FROM polls
WHERE presentationid = $1
ORDER BY createdat ASC;

-- name: GetPollOptions :many
SELECT optionkey, optionvalue
FROM options
WHERE pollid = $1
ORDER BY optionkey;


-- name: CreateVote :exec
INSERT INTO votes (id, pollid, optionkey, clientid)
VALUES ($1, $2, $3, $4);

-- name: GetVote :many
SELECT v.id AS voteID, v.pollid, o.id AS optionid, o.optionkey, o.optionvalue,v.clientid
FROM votes AS v
JOIN options AS o ON v.optionkey = o.optionkey
WHERE v.pollid = $1;

-- name: UpdateCurrPollIndex :one
UPDATE presentations
SET currentpollindex = $1
WHERE id = $2
RETURNING currentpollindex;

-- name: GetPollsCount :one
SELECT COUNT(*) AS polls_count
FROM polls
WHERE presentationid = $1;

-- name: GetPollVotes :many
SELECT
  votes.optionkey,
  votes.clientid
FROM votes
JOIN polls ON votes.pollid = polls.id
JOIN presentations ON polls.presentationID = presentations.id
WHERE presentations.id = $1
AND polls.id = $2
ORDER BY votes.optionkey;
