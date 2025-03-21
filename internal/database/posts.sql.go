// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: posts.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createPost = `-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES (
	$1,
	$2,
	$3,
	$4,
	$5,
	$6,
	$7,
	$8
	)
	RETURNING id, created_at, updated_at, title, url, description, published_at, feed_id
`

type CreatePostParams struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Title       sql.NullString
	Url         sql.NullString
	Description sql.NullString
	PublishedAt time.Time
	FeedID      uuid.NullUUID
}

func (q *Queries) CreatePost(ctx context.Context, arg CreatePostParams) (Post, error) {
	row := q.db.QueryRowContext(ctx, createPost,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Title,
		arg.Url,
		arg.Description,
		arg.PublishedAt,
		arg.FeedID,
	)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Title,
		&i.Url,
		&i.Description,
		&i.PublishedAt,
		&i.FeedID,
	)
	return i, err
}

const getPostsForUser = `-- name: GetPostsForUser :one
SELECT p.id, p.created_at, p.updated_at, p.title, p.url, p.description, p.published_at, p.feed_id 
FROM posts p
LEFT JOIN feed_follows ff ON ff.feed_id = p.feed_id
WHERE ff.user_id = $1
ORDER BY p.published_at DESC
LIMIT $2
`

type GetPostsForUserParams struct {
	UserID uuid.NullUUID
	Limit  int32
}

func (q *Queries) GetPostsForUser(ctx context.Context, arg GetPostsForUserParams) (Post, error) {
	row := q.db.QueryRowContext(ctx, getPostsForUser, arg.UserID, arg.Limit)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Title,
		&i.Url,
		&i.Description,
		&i.PublishedAt,
		&i.FeedID,
	)
	return i, err
}
