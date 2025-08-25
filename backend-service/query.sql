-- name: CreateUser :one
INSERT INTO users (name, username, email, password)
VALUES ($1, $2, $3, $4)
RETURNING id, username, email;

-- name: GetUserById :one
SELECT id, username, email, password
FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT id, username, email, password
FROM users
WHERE email = $1;

-- name: CreateProfile :exec
INSERT INTO profiles (user_id, exp_needed)
VALUES ($1, $2);

-- name: CreateStatistics :exec
INSERT INTO statistics (user_id)
VALUES ($1);

-- name: CreateLog :one
INSERT INTO logs (user_id, text, is_system, is_marked, is_private)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, text, is_system, is_marked, is_private;

-- name: GetLogs :many
SELECT text, created_at, is_marked, is_system, is_private
FROM logs
WHERE user_id = $1 AND is_marked = $2 AND is_system = $3 AND is_private = $4;
