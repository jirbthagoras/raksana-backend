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

-- name: GetStatisticByUserID :one
SELECT * FROM statistics WHERE user_id = $1;

-- name: UpdateLongestStreak :exec
UPDATE statistics SET longest_streak = $1
WHERE user_id = $2;

-- name: GetProfileByUserId :one
SELECT current_exp, exp_needed, level, points
FROM profiles
WHERE user_id = $1;

-- name: IncreaseExp :one
UPDATE profiles 
SET current_exp = current_exp + @exp_gain::int
WHERE user_id = @user_id::int
RETURNING current_exp, exp_needed, level;

-- name: UpdateLevelAndExpNeeded :one
UPDATE profiles                                  m
SET exp_needed = $1, level = level + 1
WHERE user_id = $2
RETURNING level;

-- name: CreatePacket :one
INSERT INTO packets  (user_id, name, target, description, expected_task, task_per_day)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id;

-- name: CreateHabit :one
INSERT INTO habits (packet_id, name, description, difficulty, locked, weight)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id;

-- name: CountActivePacketsByUserId :one
SELECT COUNT(*) FROM packets 
WHERE user_id = $1 AND completed = false;

-- name: GetActivePacketsByUserId :one
SELECT * FROM packets
WHERE user_id = $1 AND completed = false;

-- name: GetHabitsByPacketId :many
SELECT * FROM habits
WHERE packet_id = $1;

-- name: GetUnlockedHabitsByPacketId :many
SELECT * FROM habits
WHERE packet_id = $1 AND locked = false;

-- name: UnlockHabit :exec
UPDATE habits
SET locked = false
WHERE id = $1;

-- name: CountAssignedTask :one
SELECT
  COUNT(*) FILTER (WHERE completed = true) AS completed_task
FROM tasks
WHERE packet_id = $1 AND user_id = $2;

-- name: CreateTask :one
INSERT INTO tasks(habit_id, user_id, packet_id, name, description, difficulty)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetTodayTasks :many
SELECT *
FROM tasks
WHERE user_id = $1
  AND DATE(created_at) = CURRENT_DATE;
