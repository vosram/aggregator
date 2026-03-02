-- name: AddFeed :one
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

-- name: ListFeeds :many
SELECT feeds.*, users.name AS user_name FROM feeds
INNER JOIN users ON feeds.user_id = users.id;

-- name: GetFeedByUrl :one
SELECT * FROM feeds
WHERE url = $1;

-- name: CreateFeedFollow :one
WITH inserted_feed_follows AS (
  INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
  VALUES ($1, $2, $3, $4, $5)
  RETURNING *
)
SELECT
  inserted_feed_follows.*,
  feeds.name AS feed_name,
  users.name AS user_name
FROM inserted_feed_follows
INNER JOIN feeds ON feeds.id = inserted_feed_follows.feed_id
INNER JOIN users ON users.id = inserted_feed_follows.user_id;

-- name: GetFeedFollowsForUser :many
SELECT
  feed_follows.*,
  feeds.name AS feed_name,
  users.name AS user_name
FROM feed_follows
INNER JOIN feeds ON feeds.id = feed_follows.feed_id
INNER JOIN users ON users.id = feed_follows.user_id
WHERE feed_follows.user_id = $1;