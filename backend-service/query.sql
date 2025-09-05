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

-- name: CreateProfile :one
INSERT INTO profiles (user_id, exp_needed)
VALUES ($1, $2)
RETURNING *;

-- name: CreateStatistics :exec
INSERT INTO statistics (user_id)
VALUES ($1);

-- name: CreateLog :one
INSERT INTO logs (user_id, text, is_system, is_private)
VALUES ($1, $2, $3, $4)
RETURNING id, text, is_system, is_private;

-- name: GetLogs :many
SELECT text, created_at, is_system, is_private
FROM logs
WHERE user_id = $1 AND is_private = $2
ORDER BY created_at DESC;

-- name: GetUserStatistic :one
SELECT * FROM statistics WHERE user_id = $1;

-- name: UpdateLongestStreak :exec
UPDATE statistics SET longest_streak = $1
WHERE user_id = $2;

-- name: GetUserProfile :one
SELECT p.*
FROM profiles p
JOIN users u ON p.user_id = u.id
WHERE p.user_id = $1;

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

-- name: CountUserActivePackets :one
SELECT COUNT(*) FROM packets 
WHERE user_id = $1 AND completed = false;

-- name: GetUserActivePackets :one
SELECT * FROM packets
WHERE user_id = $1 AND completed = false;

-- name: GetAllPackets :many
SELECT * FROM packets
WHERE user_id = $1;

-- name: GetPacketHabits :many
SELECT 
  *
FROM habits
WHERE packet_id = $1;

-- name: GetPacketUnlockedHabits :many
SELECT * FROM habits
WHERE packet_id = $1 AND locked = false;

-- name: UnlockHabit :exec
UPDATE habits
SET locked = false
WHERE id = $1;

-- name: CountPacketTasks :one
SELECT
  COUNT(*) FILTER (WHERE completed = true) AS completed_task,
  COUNT(*) as assigned_task
FROM tasks
WHERE packet_id = $1 AND user_id = $2;

-- name: CountUserTask :one
SELECT 
  COUNT (*) as assigned_task,
  COUNT (*) FILTER (WHERE completed = true) as completed_task
FROM tasks
WHERE user_id = $1;

-- name: CreateTask :one
INSERT INTO tasks(habit_id, user_id, packet_id, name, description, difficulty)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetTodayTasks :many
SELECT *
FROM tasks
WHERE user_id = $1
  AND DATE(created_at) = CURRENT_DATE
ORDER BY id;

-- name: GetTaskById :one
SELECT * FROM tasks
WHERE id = $1;

-- name: CompleteTask :one
UPDATE tasks
SET completed = true, updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND user_id = $2
  AND DATE(created_at) = CURRENT_DATE
RETURNING *;

-- name: IncreasePacketCompletedTask :exec
UPDATE packets
SET completed_task = completed_task + 1
WHERE id = $1;

-- name: CompletePacket :exec
UPDATE packets
SET completed = true
WHERE id = $1;

-- name: GetUserProfileStatistic :one
SELECT 
    u.id AS user_id,
    u.name,
    u.username,
    u.email,
    p.current_exp,
    p.exp_needed,
    p.level,
    p.points,
    p.profile_key,
    s.challenges,
    s.events,
    s.quests,
    s.treasures,
    s.longest_streak,
    s.tree_grown
FROM users u
JOIN profiles p ON u.id = p.user_id
JOIN statistics s ON u.id = s.user_id
WHERE u.id = $1;

-- name: GetPacketDetail :one
SELECT 
    p.id AS packet_id,
    p.user_id,
    u.name AS user_name,
    u.username,
    p.name AS packet_name,
    p.target,
    p.description,
    p.completed_task,
    p.expected_task,
    p.task_per_day,
    p.completed,
    p.created_at
FROM packets p
JOIN users u ON u.id = p.user_id
WHERE p.id = $1;

-- name: CreateMemory :one
INSERT INTO memories(user_id, file_key, description)
VALUES ($1, $2, $3)
RETURNING id;

-- name: GetMemoryWithParticipation :many
SELECT 
    m.id AS memory_id,
    m.file_key,
    m.description AS memory_description,
    m.created_at AS memory_created_at,
    u.id AS user_id,
    u.name AS user_name,
    CASE 
        WHEN p.id IS NOT NULL THEN TRUE 
        ELSE FALSE 
    END AS is_participation,
    p.challenge_id,
    c.day,
    c.difficulty,
    d.name AS challenge_name,
    d.point_gain
FROM memories m
JOIN users u ON m.user_id = u.id
LEFT JOIN participations p ON m.id = p.memory_id
LEFT JOIN challenges c ON p.challenge_id = c.id
LEFT JOIN details d ON c.detail_id = d.id
WHERE m.user_id = $1
ORDER BY m.created_at ASC;

-- name: DeleteMemory :one
DELETE FROM memories
WHERE id = $1 AND user_id = $2
RETURNING file_key;

-- name: UpdateUserProfile :exec
UPDATE profiles
SET profile_key = $1
WHERE user_id = $2;

-- name: GetLastWeekTasks :many
SELECT *
FROM tasks
WHERE created_at >= NOW() - INTERVAL '7 weeks' AND user_id = $1
ORDER BY created_at DESC;

-- name: GetLatestRecap :one
SELECT *
FROM recaps
WHERE user_id = $1
  AND type = 'weekly'
ORDER BY created_at DESC
LIMIT 1;

-- name: GetWeeklyRecaps :many
SELECT * FROM recaps
WHERE user_id = $1 AND type = 'weekly';

-- name: CreateWeeklyRecap :exec
INSERT INTO recaps(user_id, summary, tips, assigned_task, completed_task, completion_rate, growth_rating, type)
VALUES ($1, $2, $3, $4, $5, $6, $7, 'weekly');

-- name: CreateParticipation :one
INSERT INTO participations(challenge_id, user_id, memory_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetChallengeWithDetail :one
SELECT 
    c.id AS challenge_id,
    c.day,
    c.difficulty,
    d.id AS detail_id,
    d.name,
    d.description,
    d.point_gain,
    d.created_at,
    d.updated_at
FROM challenges c
JOIN details d ON c.detail_id = d.id
WHERE c.id = $1
LIMIT 1;

-- name: IncreaseUserPoints :one
UPDATE profiles
SET points = points + $1
WHERE user_id = $2
RETURNING *;

-- name: IncreaseChallengesFieldByOne :one
UPDATE statistics
SET challenges = challenges + 1
WHERE user_id = $1
RETURNING *;

-- name: IncreaseQuestsFieldByOne :one
UPDATE statistics
SET quests = quests + 1
WHERE user_id = $1
RETURNING *;

-- name: IncreaseEventsFieldByOne :one
UPDATE statistics
SET events = events + 1
WHERE user_id = $1
RETURNING *;

-- name: IncreaseTreasuresFieldByOne :one
UPDATE statistics
SET treasures = treasures + 1
WHERE user_id = $1
RETURNING *;

-- name: CheckParticipation :one
SELECT
COUNT (*) FILTER (WHERE user_id = $1 AND challenge_id = $2)
FROM participations;

-- name: GetTodayChallenge :one
SELECT 
    c.id AS challenge_id,
    c.day,
    c.difficulty,
    d.id AS detail_id,
    d.name,
    d.description,
    d.point_gain,
    d.created_at,
    d.updated_at
FROM challenges c
JOIN details d ON c.detail_id = d.id
ORDER BY d.created_at DESC
LIMIT 1;

-- name: GetAllChallenges :many
SELECT 
c.id AS challenge_id,
c.day,
c.difficulty,
d.id AS detail_id,
d.name,
d.description,
d.point_gain,
d.created_at,
d.updated_at
FROM challenges c
JOIN details d ON c.detail_id = d.id
ORDER BY d.created_at DESC;

-- name: GetMemoriesByChallengeID :many
SELECT 
    m.id AS memory_id,
    m.file_key,
    m.description,
    m.created_at AS memory_created_at,
    u.id AS user_id,
    u.name AS user_name,
    u.username,
    u.email
FROM participations p
JOIN memories m ON p.memory_id = m.id
JOIN users u ON m.user_id = u.id
WHERE p.challenge_id = $1
ORDER BY m.created_at DESC;
