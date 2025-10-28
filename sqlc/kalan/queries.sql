-- name: CreateKalan :execresult
INSERT INTO kalan (entry, pos, gloss, notes) VALUES (?, ?, ?, ?);

-- name: ReadKalan :many
SELECT * FROM kalan ORDER BY id;

-- name: ReadKalanByEntry :one
SELECT * FROM kalan WHERE entry = ? LIMIT 1;

-- name: ReadKalanById :one
SELECT * FROM kalan WHERE id = ? LIMIT 1;

-- name: ReadKalanBySearch :many
SELECT *
FROM kalan
WHERE (
        sqlc.arg (isEntry) = True
        AND entry LIKE CONCAT('%', sqlc.arg (search), '%')
    )
    OR (
        sqlc.arg (isPos) = True
        AND pos LIKE CONCAT('%', sqlc.arg (search), '%')
    )
    OR (
        sqlc.arg (isGloss) = True
        AND gloss LIKE CONCAT('%', sqlc.arg (search), '%')
    )
    OR (
        sqlc.arg (isNotes) = True
        AND notes LIKE CONCAT('%', sqlc.arg (search), '%')
    )
ORDER BY
    CASE sqlc.arg (sort)
        WHEN 'entry' THEN entry
        WHEN 'pos' THEN pos
        WHEN 'gloss' THEN gloss
        WHEN 'notes' THEN notes
        ELSE id
    END
LIMIT ?
OFFSET
    ?;