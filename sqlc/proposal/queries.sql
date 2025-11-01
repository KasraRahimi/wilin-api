-- name: CreateProposal :execresult
INSERT INTO
    proposals (
        user_id,
        entry,
        pos,
        gloss,
        notes
    )
VALUES (?, ?, ?, ?, ?);

-- name: ReadAllProposalsWithUsername :many
SELECT p.id, p.user_id, u.username, p.entry, p.pos, p.gloss, p.notes
FROM proposals p
    JOIN users u ON u.id = p.user_id;

-- name: ReadProposalsByUserIDWithUsername :many
SELECT p.id, p.user_id, u.username, p.entry, p.pos, p.gloss, p.notes
FROM proposals p
    JOIN users u ON u.id = p.user_id
WHERE
    u.id = ?;

-- name: ReadProposalByIDWithUsername :one
SELECT p.id, p.user_id, u.username, p.entry, p.pos, p.gloss, p.notes
FROM proposals p
    JOIN users u ON u.id = p.user_id
WHERE
    p.id = ?
LIMIT 1;

-- name: Update :execresult
UPDATE proposals
SET
    user_id = ?,
    entry = ?,
    pos = ?,
    gloss = ?,
    notes = ?
WHERE
    id = ?;

-- name: Delete :execresult
DELETE FROM proposals WHERE id = ?;