-- name: CreateVote :exec
INSERT INTO votes
  (pollid, optionkey, clientid)
VALUES
  ($1, $2, $3);

-- name: GetVotes :many
SELECT
  votes.optionkey as "key",
  votes.clientid as "client_id"
FROM votes
  JOIN polls ON votes.pollid = polls.id
  JOIN presentations ON polls.presentationID = presentations.id
WHERE presentations.id = sqlc.arg(presentation_id)
  AND polls.id = sqlc.arg(poll_id)
ORDER BY votes.optionkey;

