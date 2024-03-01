-- name: ResetUsers :exec
DELETE FROM users;

-- name: ResetFeeds :exec
DELETE FROM feeds;

-- name: ResetFeedFollows :exec
DELETE FROM feed_follows;
