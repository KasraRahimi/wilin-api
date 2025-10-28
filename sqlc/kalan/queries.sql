-- name: ReadKalan :many
SELECT * FROM kalan ORDER BY id;

-- name: ReadKalanByEntry :one
SELECT * FROM kalan WHERE entry = ? LIMIT 1;

-- name: CreateKalan :execresult
INSERT INTO kalan (entry, pos, gloss, notes) VALUES (?, ?, ?, ?);