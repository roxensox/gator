-- name: CreatePost :one
INSERT INTO posts (
	id,
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
	$2,
	$3,
	$4,
	$5,
	$6,
	$7
) RETURNING *;

-- name: GetPostsForUser :many
SELECT posts.*, feed_follows.u_id
FROM feed_follows 
INNER JOIN posts
ON posts.feed_id = feed_follows.f_id
WHERE feed_follows.u_id = $1
ORDER BY published_at DESC
LIMIT $2;
