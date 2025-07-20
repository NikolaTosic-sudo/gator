-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name)
VALUES (
  $1,
  $2,
  $3,
  $4
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE name = $1;

-- name: GetUserById :one
SELECT id, name FROM users WHERE id = $1;

-- name: ResetDB :exec
DELETE FROM users;
DELETE FROM feeds;
DELETE FROM feed_follows;

-- name: GetUsers :many
SELECT * FROM users;