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
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (url)
DO UPDATE
SET updated_at = $3,
  title = $4,
  description = $6
RETURNING *;

-- name: GetPostsByUser :many
SELECT * FROM posts AS p
  INNER JOIN feeds AS f
    ON p.feed_id = f.id
  INNER JOIN feed_follows AS ff
    ON f.id = ff.feed_id
WHERE ff.user_id = $1
ORDER BY published_at DESC
LIMIT $2;
