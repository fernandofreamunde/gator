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

-- name: GetUserFeeds :many
SELECT f.*, u.*
FROM feeds f
LEFT JOIN users u ON u.id = f.user_id
WHERE f.user_id = $1;

-- name: GetFeeds :many
SELECT f.*, u.name as username
FROM feeds f
LEFT JOIN users u ON u.id = f.user_id;

