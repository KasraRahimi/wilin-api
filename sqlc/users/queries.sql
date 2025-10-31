-- name: CreateUser :execresult
INSERT INTO
    users (
        email,
        username,
        password,
        role
    )
VALUES (?, ?, ?, ?);

-- name: ReadUserByUsername :one
SELECT
    id,
    email,
    username,
    password,
    role
FROM users
WHERE
    username = ?
LIMIT 1;

-- name: ReadUserByEmail :one
SELECT
    id,
    email,
    username,
    password,
    role
FROM users
WHERE
    email = ?
LIMIT 1;