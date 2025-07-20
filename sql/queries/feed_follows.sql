-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
  INSERT INTO feed_follows(
  created_at,
  updated_at,
  user_id,
  feed_id
) VALUES (
  $1,
  $2,
  $3,
  $4
)
RETURNING *
)

SELECT inserted_feed_follow.*, feeds.name AS feed_name, users.name AS user_name
FROM inserted_feed_follow
JOIN feeds ON inserted_feed_follow.feed_id = feeds.id
JOIN users ON inserted_feed_follow.user_id = users.id;

-- name: GetFeedFollowsForUser :many
SELECT feed_follows.*, feeds.name AS feed_name, users.name AS user_name FROM feed_follows
JOIN feeds ON feed_follows.feed_id = feeds.id
JOIN users ON feed_follows.user_id = users.id
WHERE feed_follows.user_id = $1;

-- name: RemoveFeedFollow :one
DELETE FROM feed_follows WHERE feed_follows.user_id = $1 AND feed_follows.feed_id = $2
RETURNING *;