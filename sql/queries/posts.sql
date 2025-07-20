-- name: CreatePost :one
INSERT INTO posts(
  created_at,
  updated_at,
  title,
  url,
  description,
  published_at,
  feed_id
) VALUES (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6,
  $7
) RETURNING *;

-- name: GetPostsForUser :many
SELECT posts.*, users.name AS user_name, users.id as user_id FROM posts
JOIN feeds ON feed_id = feeds.id
JOIN users ON feeds.user_id = user_id
WHERE user_id = $1
LIMIT $2;