-- name: CreateVote :exec
INSERT INTO votes
  (id, pollid, optionkey, clientid)
VALUES
  ($1, $2, $3, $4);

-- name: GetVote :many
SELECT v.id AS voteID, v.pollid, o.id AS optionid, o.optionkey, o.optionvalue, v.clientid
FROM votes AS v
  JOIN options AS o ON v.optionkey = o.optionkey
WHERE v.pollid = $1;
