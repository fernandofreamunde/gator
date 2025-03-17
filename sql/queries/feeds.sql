-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
	$1,
	$2,
	$3,
	$4,
	$5,
	$6
	)
	RETURNING *;

-- name: GetFeed :one
SELECT * FROM feeds
WHERE id = $1
LIMIT 1;

-- name: GetFeedByUrl :one
SELECT * FROM feeds
WHERE name = $1
LIMIT 1;

-- name: GetFeedByName :one
SELECT * FROM feeds
WHERE name = $1
LIMIT 1;

-- name: DeleteFeed :exec
DELETE FROM feeds WHERE id = $1;

-- name: GetUserFeedss :many
SELECT * FROM feeds WHERE user_id = $1;

