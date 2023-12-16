-- name: CreateOption :exec
INSERT INTO options
  (pollid, optionkey, optionvalue)
VALUES
  ($1, $2, $3);

-- name: GetOptionByKey :one
SELECT *
FROM options
WHERE optionkey = $1;


-- name: GetPollOptions :many
SELECT optionkey, optionvalue
FROM options
WHERE pollid = $1
ORDER BY optionkey;

