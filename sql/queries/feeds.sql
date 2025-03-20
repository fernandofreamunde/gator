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
WHERE url = $1
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

-- name: MarkFeedFetched :exec
UPDATE feeds
SET 
	last_fetched_at = NOW(),
	updated_at = NOW()
WHERE id = $1;

-- name: GetNextFeedToFetch :one
SELECT id, url
FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;

