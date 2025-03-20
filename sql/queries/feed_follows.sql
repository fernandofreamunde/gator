-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
VALUES (
	$1,
	$2,
	$3,
	$4,
	$5
	)
	RETURNING *
)
SELECT
	inserted_feed_follow.*,
	f.name AS feed_name,
	u.name AS user_name
FROM inserted_feed_follow
INNER JOIN users u ON u.id = inserted_feed_follow.user_id
INNER JOIN feeds f ON f.id = inserted_feed_follow.feed_id
;

-- name: GetUserFeedFollows :many
SELECT ff.*, u.*, f.*
FROM feed_follows ff
INNER JOIN users u ON u.id = ff.user_id
INNER JOIN feeds f ON f.id = ff.feed_id
WHERE ff.user_id = $1;

-- name: DeleteUserFeedFollow :exec
DELETE FROM feed_follows
WHERE id = $1;

-- name: GetUserFeedFollowByUrl :one
SELECT ff.*
FROM feed_follows ff
INNER JOIN users u ON u.id = ff.user_id
INNER JOIN feeds f ON f.id = ff.feed_id
WHERE ff.user_id = $1
AND f.url = $2;

