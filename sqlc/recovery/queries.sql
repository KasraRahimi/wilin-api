-- name: Create :execresult
INSERT INTO recoveries (id, user_id, expired_at) VALUES (?, ?, ?);

-- name: ReadByID :one
SELECT * FROM recoveries WHERE id = ? LIMIT 1;

-- name: ReadByUserID :many
SELECT * FROM recoveries where user_id = ?;

-- name: DeleteByID :execresult
DELETE FROM recoveries WHERE id = ?;