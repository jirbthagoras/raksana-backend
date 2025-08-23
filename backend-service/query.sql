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
