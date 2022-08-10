-- name: CreateUser :one
INSERT INTO users (username, email, full_name, password)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUser :one
SELECT *
FROM users
WHERE username = sqlc.arg(username)
LIMIT 1;