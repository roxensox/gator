-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
	INSERT INTO feed_follows (
		id,
		created_at,
		updated_at,
		u_id,
		f_id
	) VALUES (
		$1,
		$2,
		$3,
		$4,
		$5
		)
	RETURNING *
) SELECT inserted_feed_follow.*, 
	users.name AS user_name,
	feeds.name AS feed_name,
	feeds.url AS feed_url
FROM inserted_feed_follow
INNER JOIN users ON inserted_feed_follow.u_id = users.id
INNER JOIN feeds ON inserted_feed_follow.f_id = feeds.id;

-- name: GetFeedFollowsForUser :many
SELECT feed_follows.*,
	users.name AS user_name,
	feeds.name AS feed_name,
	feeds.url AS feed_url
FROM feed_follows
INNER JOIN users ON feed_follows.u_id = users.id
INNER JOIN feeds ON feed_follows.f_id = feeds.id
WHERE feed_follows.u_id = $1;

-- name: UnfollowFeed :exec
DELETE FROM feed_follows 
	WHERE u_id = $1
	AND f_id = $2;
