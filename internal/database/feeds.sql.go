// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: feeds.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createFeed = `-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
	$1,
	$2,
	$3,
	$4,
	$5,
	$6
	)
	RETURNING id, created_at, updated_at, name, url, user_id, last_fetched_at
`

type CreateFeedParams struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      sql.NullString
	Url       sql.NullString
	UserID    uuid.NullUUID
}

func (q *Queries) CreateFeed(ctx context.Context, arg CreateFeedParams) (Feed, error) {
	row := q.db.QueryRowContext(ctx, createFeed,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Name,
		arg.Url,
		arg.UserID,
	)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Url,
		&i.UserID,
		&i.LastFetchedAt,
	)
	return i, err
}

const deleteFeed = `-- name: DeleteFeed :exec
DELETE FROM feeds WHERE id = $1
`

func (q *Queries) DeleteFeed(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteFeed, id)
	return err
}

const getFeed = `-- name: GetFeed :one
SELECT id, created_at, updated_at, name, url, user_id, last_fetched_at FROM feeds
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetFeed(ctx context.Context, id uuid.UUID) (Feed, error) {
	row := q.db.QueryRowContext(ctx, getFeed, id)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Url,
		&i.UserID,
		&i.LastFetchedAt,
	)
	return i, err
}

const getFeedByName = `-- name: GetFeedByName :one
SELECT id, created_at, updated_at, name, url, user_id, last_fetched_at FROM feeds
WHERE name = $1
LIMIT 1
`

func (q *Queries) GetFeedByName(ctx context.Context, name sql.NullString) (Feed, error) {
	row := q.db.QueryRowContext(ctx, getFeedByName, name)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Url,
		&i.UserID,
		&i.LastFetchedAt,
	)
	return i, err
}

const getFeedByUrl = `-- name: GetFeedByUrl :one
SELECT id, created_at, updated_at, name, url, user_id, last_fetched_at FROM feeds
WHERE url = $1
LIMIT 1
`

func (q *Queries) GetFeedByUrl(ctx context.Context, url sql.NullString) (Feed, error) {
	row := q.db.QueryRowContext(ctx, getFeedByUrl, url)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Url,
		&i.UserID,
		&i.LastFetchedAt,
	)
	return i, err
}

const getFeeds = `-- name: GetFeeds :many
SELECT f.id, f.created_at, f.updated_at, f.name, f.url, f.user_id, f.last_fetched_at, u.name as username
FROM feeds f
LEFT JOIN users u ON u.id = f.user_id
`

type GetFeedsRow struct {
	ID            uuid.UUID
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Name          sql.NullString
	Url           sql.NullString
	UserID        uuid.NullUUID
	LastFetchedAt sql.NullTime
	Username      sql.NullString
}

func (q *Queries) GetFeeds(ctx context.Context) ([]GetFeedsRow, error) {
	rows, err := q.db.QueryContext(ctx, getFeeds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetFeedsRow
	for rows.Next() {
		var i GetFeedsRow
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Name,
			&i.Url,
			&i.UserID,
			&i.LastFetchedAt,
			&i.Username,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getNextFeedToFetch = `-- name: GetNextFeedToFetch :one
SELECT id, url
FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1
`

type GetNextFeedToFetchRow struct {
	ID  uuid.UUID
	Url sql.NullString
}

func (q *Queries) GetNextFeedToFetch(ctx context.Context) (GetNextFeedToFetchRow, error) {
	row := q.db.QueryRowContext(ctx, getNextFeedToFetch)
	var i GetNextFeedToFetchRow
	err := row.Scan(&i.ID, &i.Url)
	return i, err
}

const getUserFeeds = `-- name: GetUserFeeds :many
SELECT f.id, f.created_at, f.updated_at, f.name, f.url, f.user_id, f.last_fetched_at, u.id, u.created_at, u.updated_at, u.name
FROM feeds f
LEFT JOIN users u ON u.id = f.user_id
WHERE f.user_id = $1
`

type GetUserFeedsRow struct {
	ID            uuid.UUID
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Name          sql.NullString
	Url           sql.NullString
	UserID        uuid.NullUUID
	LastFetchedAt sql.NullTime
	ID_2          uuid.NullUUID
	CreatedAt_2   sql.NullTime
	UpdatedAt_2   sql.NullTime
	Name_2        sql.NullString
}

func (q *Queries) GetUserFeeds(ctx context.Context, userID uuid.NullUUID) ([]GetUserFeedsRow, error) {
	rows, err := q.db.QueryContext(ctx, getUserFeeds, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUserFeedsRow
	for rows.Next() {
		var i GetUserFeedsRow
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Name,
			&i.Url,
			&i.UserID,
			&i.LastFetchedAt,
			&i.ID_2,
			&i.CreatedAt_2,
			&i.UpdatedAt_2,
			&i.Name_2,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const markFeedFetched = `-- name: MarkFeedFetched :exec
UPDATE feeds
SET 
	last_fetched_at = NOW(),
	updated_at = NOW()
WHERE id = $1
`

func (q *Queries) MarkFeedFetched(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, markFeedFetched, id)
	return err
}
