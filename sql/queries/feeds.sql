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

-- name: GetAllFeeds :many
SELECT name, url, user_id FROM feeds;

-- name: GetFeedByUrl :one
SELECT id, name, url, user_id FROM feeds WHERE url = $1;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET updated_at = $1, last_fetched_at = $1
WHERE id = $2;

-- name: GetNextFeedToFetch :one
SELECT * FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST;