-- name: CreateUser :one
INSERT INTO users (username, email, full_name, password)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUserByUsername :one
SELECT *
FROM users
WHERE LOWER(username) = LOWER(sqlc.arg(username))
LIMIT 1;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE LOWER(email) = LOWER(sqlc.arg(email))
LIMIT 1;

-- name: GetUserByUsernameOrEmail :one
SELECT *
FROM users
WHERE LOWER(username) = LOWER(sqlc.arg(username)) OR LOWER(email) = LOWER(sqlc.arg(email))
LIMIT 1;