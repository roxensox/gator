-- name: CreateFeedFollow :one
SELECT * FROM users JOIN (SELECT * FROM feed_follows JOIN feeds ON feed_id = feeds.id)
