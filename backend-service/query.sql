-- name: CreateUser :one
INSERT INTO users (name, username, email, password)
VALUES ($1, $2, $3, $4)
RETURNING id, username, email;

-- name: GetUserById :one
SELECT id, username, email, password, is_admin
FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT id, username, email, password, is_admin 
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
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: GetPacketHabits :many
SELECT 
  *
FROM habits
WHERE packet_id = $1;

-- name: GetLockedHabits :many
SELECT 
  *
FROM habits
WHERE packet_id = $1 AND locked = true;

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
    u.username AS user_name,
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
ORDER BY m.created_at DESC;

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
WHERE created_at >= NOW() - INTERVAL '1 weeks' AND user_id = $1
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
WHERE user_id = $1 AND type = 'weekly'
ORDER BY created_at DESC;

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
ORDER BY created_at DESC;

-- name: GetChallengeWithDetailById :one
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
ORDER BY created_at DESC;

-- name: IncreaseUserPoints :one
UPDATE profiles
SET points = points + $1
WHERE user_id = $2
RETURNING *;

-- name: DecreaseUserPoints :one
UPDATE profiles
SET points = points - $1
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

-- name: GetTreasureByCodeId :one
SELECT * FROM treasures
WHERE code_id = $1 AND claimed = false;

-- name: DeactivateTreasure :exec
UPDATE treasures
SET claimed = true
WHERE id = $1;

-- name: CreateClaimed :exec
INSERT INTO claimed(user_id, treasure_id)
VALUES ($1, $2);

-- name: GetAllClaimedTreasure :many
SELECT 
  t.id AS id,
  t.name AS name,
  c.created_at AS claimed_at,
  t.point_gain AS point_gain
FROM claimed c
JOIN treasures t ON c.treasure_id = t.id
WHERE c.user_id = $1
ORDER BY c.created_at DESC;

-- name: GetUncompletedQuestByCodeId :one
SELECT 
  q.id AS id,
  q.location AS location,
  q.max_contributors AS max_contributors,
  q.latitude AS latitude,
  q.longitude AS longitude,
  d.name AS name,
  d.description AS description,
  d.point_gain AS point_gain
FROM quests q
JOIN details d ON q.detail_id = d.id
WHERE code_id = $1 AND finished = false;

-- name: GetQuestByCodeId :one
SELECT 
  q.id AS id,
  q.location AS location,
  q.max_contributors AS max_contributors,
  q.latitude AS latitude,
  q.longitude AS longitude,
  d.name AS name,
  d.description AS description,
  d.point_gain AS point_gain,
  q.finished
FROM quests q
JOIN details d ON q.detail_id = d.id
WHERE q.code_id = $1;

-- name: FinsihQuest :exec
UPDATE quests
SET finished = true
WHERE id = $1;

-- name: CountQuestContributors :many 
SELECT
  u.id AS id,
  u.username AS username
FROM contributions c
JOIN users u ON c.user_id = u.id
WHERE quest_id = $1;

-- name: CreateContributions :one
INSERT INTO contributions(user_id, quest_id)
VALUES ($1, $2)
RETURNING *;

-- name: GetUserContributions :many
SELECT 
    c.id               AS contribution_id,
    d.name As name,
    d.description AS description,
    q.latitude,
    q.longitude,
    d.point_gain AS point_gain
FROM contributions c
JOIN quests q ON c.quest_id = q.id
JOIN details d ON q.detail_id = d.id
WHERE c.user_id = $1
ORDER BY c.created_at DESC;

-- name: GetContributionDetails :one
SELECT 
    c.id               AS id,
    c.created_at       AS contribution_date,
    q.code_id,
    q.id AS quest_id,
    q.location,
    q.latitude,
    q.longitude,
    q.max_contributors,
    d.name             AS name,
    d.description      AS description,
    d.point_gain,
    d.created_at       AS created_at
FROM contributions c
JOIN quests q ON c.quest_id = q.id
JOIN details d ON q.detail_id = d.id
WHERE c.id = $1;

-- name: GetUserAttendances :many
SELECT 
    a.id AS attendance_id,
    a.created_at AS attended_at,
    a.contact_number,
    e.id AS event_id,
    e.location,
    e.latitude,
    e.longitude,
    e.contact,
    e.starts_at,
    e.ends_at,
    a.attended,
    e.cover_key,
    d.id AS detail_id,
    d.name AS detail_name,
    d.description AS detail_description,
    d.point_gain
FROM attendances a
JOIN events e ON a.event_id = e.id
JOIN details d ON e.detail_id = d.id
WHERE a.user_id = $1
ORDER BY a.created_at DESC;

-- name: GetUserPendingAttendances :many
SELECT 
    a.id AS attendance_id,
    e.id AS event_id,
    a.created_at AS registered_at,
    a.contact_number,
    e.location,
    e.latitude,
    e.longitude,
    e.contact,
    e.starts_at,
    e.ends_at,
    e.cover_key,
    d.name AS detail_name,
    d.description AS detail_description,
    d.point_gain
FROM attendances a
JOIN events e ON a.event_id = e.id
JOIN details d ON e.detail_id = d.id
WHERE a.user_id = $1 AND a.attended = false;

-- name: Attend :exec
UPDATE attendances
SET attended = true
WHERE id = $1;

-- name: GetUserAttendanceByUserId :one
SELECT 
    a.id AS attendance_id,
    a.created_at AS attended_at,
    a.contact_number,
    e.id AS event_id,
    e.location,
    e.latitude,
    e.longitude,
    e.contact,
    e.starts_at,
    e.ends_at,
    a.attended,
    e.cover_key,
    d.id AS detail_id,
    d.name AS detail_name,
    d.description AS detail_description,
    d.point_gain,
    a.created_at
FROM attendances a
JOIN events e ON a.event_id = e.id
JOIN details d ON e.detail_id = d.id
WHERE a.user_id = $1 AND e.id = $2;

-- name: GetUserAttendance :one
SELECT 
    a.id AS attendance_id,
    a.created_at AS attended_at,
    a.contact_number,
    e.id AS event_id,
    e.location,
    e.latitude,
    e.longitude,
    e.contact,
    e.starts_at,
    e.ends_at,
    a.attended,
    e.cover_key,
    d.id AS detail_id,
    d.name AS detail_name,
    d.description AS detail_description,
    d.point_gain,
    a.created_at
FROM attendances a
JOIN events e ON a.event_id = e.id
JOIN details d ON e.detail_id = d.id
WHERE a.id = $1 AND a.user_id = $2;

-- name: UpdaAttendedAt :exec
UPDATE attendances
SET attended_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: GetUserAttendanceById :one
SELECT 
    a.id AS attendance_id,
    a.created_at AS attended_at,
    a.contact_number,
    e.id AS event_id,
    e.location,
    e.latitude,
    e.longitude,
    e.contact,
    e.starts_at,
    e.ends_at,
    a.attended,
    e.cover_key,
    d.name AS name,
    d.description AS description,
    d.point_gain
FROM attendances a
JOIN events e ON a.event_id = e.id
JOIN details d ON e.detail_id = d.id
WHERE a.id = $1;

-- name: GetAttendanceDetails :one
SELECT
    a.id AS attendance_id,
    a.attended,
    a.created_at AS attended_at,
    a.contact_number,
    e.id AS event_id,
    e.location,
    e.latitude,
    e.longitude,
    e.contact,
    e.starts_at,
    e.ends_at,
    e.cover_key,
    d.id AS detail_id,
    d.name AS detail_name,
    d.description AS detail_description,
    d.point_gain
FROM attendances a
JOIN events e ON a.event_id = e.id
JOIN details d ON e.detail_id = d.id
WHERE a.id = $1;

-- name: CreateAttendance :one
INSERT INTO attendances(user_id, event_id, contact_number)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetEventByCodeId :one
SELECT
  e.id,
  d.name AS name,
  d.description AS description,
  d.point_gain AS point_gain
FROM events e
JOIN details d ON e.detail_id = d.id
WHERE e.code_id = $1;

-- name: GetEventById :one
SELECT 
  e.id,
  d.name AS name,
  d.description AS description
FROM events e
JOIN details d ON e.detail_id = d.id
WHERE e.id = $1;

-- name: GetAllEvents :many
SELECT 
    e.id,
    e.detail_id,
    e.code_id,
    e.location,
    e.latitude,
    e.longitude,
    e.contact,
    e.starts_at,
    e.ends_at,
    e.cover_key,
    d.name AS detail_name,
    d.description AS detail_description,
    d.point_gain,
    d.created_at AS detail_created_at,
    d.updated_at AS detail_updated_at,
    (e.ends_at < NOW()) AS ended
FROM events e
JOIN details d ON e.detail_id = d.id
ORDER BY e.starts_at ASC;

-- name: GetAllUser :many
SELECT
  u.id AS user_id,
  p.level AS level,
  p.profile_key AS profile_key,
  u.username AS username,
  u.email AS email
FROM profiles p
JOIN users u ON p.user_id = u.id
WHERE u.is_admin = false;

-- name: AppendHistry :exec
INSERT INTO histories(user_id, name, category, type, amount)
VALUES($1, $2, $3, $4, $5);

-- name: GetUserHistories :many
SELECT * FROM histories
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: GetLastMonthUserLogs :many
SELECT * FROM logs
WHERE user_id =  $1
AND date_trunc('month', created_at) = date_trunc('month', CURRENT_DATE)
ORDER BY created_at DESC;

-- name: GetLastMonthUserHistories :many
SELECT * FROM histories
WHERE user_id = $1
AND date_trunc('month', created_at) = date_trunc('month', CURRENT_DATE)
ORDER BY created_at DESC;

-- name: CreateMonthlyRecap :one
INSERT INTO recaps(user_id, summary, tips, assigned_task, completed_task, completion_rate, growth_rating, type)
VALUES ($1, $2, $3, $4, $5, $6, $7, 'monthly')
RETURNING id;

-- name: CreateRecapDetails :exec
INSERT INTO recap_details(monthly_recap_id, challenges, events, quests, treasures, longest_streak)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetLatestMonhtlyRecap :one
SELECT 
  *,
  (date_trunc('month', created_at) = date_trunc('month', CURRENT_DATE)) 
           AS is_this_month
FROM recaps
WHERE user_id = $1 AND type = 'monthly'
ORDER BY created_at DESC
LIMIT 1;

-- name: GetAllMonthlyRecapsWithDetails :many
SELECT r.id            AS recap_id,
       r.user_id,
       r.summary,
       r.tips,
       r.assigned_task,
       r.completed_task,
       r.completion_rate,
       r.growth_rating,
       r.type,
       r.created_at    AS recap_created_at,
       d.id            AS detail_id,
       d.challenges,
       d.events,
       d.quests,
       d.treasures,
       d.longest_streak,
       d.created_at    AS detail_created_at
FROM recaps r
LEFT JOIN recap_details d 
       ON d.monthly_recap_id = r.id
WHERE r.user_id = $1
  AND r.type = 'monthly'
ORDER BY r.created_at DESC;

-- name: GetNearestQuestWithinRadius :one
SELECT 
    id,
    clue,
    location,
    latitude,
    longitude,
    earth_distance(
        ll_to_earth($1, $2),
        ll_to_earth(latitude, longitude)
    ) AS distance_meters
FROM quests
WHERE finished = false
  AND earth_distance(
        ll_to_earth($1, $2),
        ll_to_earth(latitude, longitude)
    ) <= $3
ORDER BY distance_meters
LIMIT 1;

-- name: CreateScans :one
INSERT INTO scans(user_id, title, description, image_key)
VALUES($1, $2, $3, $4)
RETURNING *;

-- name: CreateItems :one
INSERT INTO items(scan_id, user_id, name, description, value)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: CreateGreenprint :one
INSERT INTO greenprints(title, item_id, image_key, description, sustainability_score, estimated_time)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: CreateMaterials :one
INSERT INTO materials(greenprint_id,name, description, price, quantity)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetMaterials :many
SELECT * FROM materials
WHERE greenprint_id = $1;

-- name: CreateSteps :one
INSERT INTO steps(greenprint_id, description)
VALUES ($1, $2)
RETURNING *;

-- name: GetSteps :many
SELECT * FROM steps
WHERE greenprint_id = $1;

-- name: CreateTools :one
INSERT INTO tools(greenprint_id, name, description, price)
VALUES($1, $2, $3, $4)
RETURNING *;

-- name: GetTools :many
SELECT * FROM tools
WHERE greenprint_id = $1;

-- name: GetParticipants :one
SELECT 
  COUNT (*) FILTER (WHERE challenge_id = $1) AS participants
FROM participations;

-- name: GetAllUserScans :many
SELECT * FROM scans WHERE user_id = $1;

-- name: GetItemsByScanId :many
SELECT * FROM items WHERE scan_id = $1;

-- name: GetItemsById :one
SELECT * FROM items WHERE id = $1;

-- name: GetGreenprints :one
SELECT * FROM greenprints
WHERE item_id = $1;

-- name: GetGreenprintsById :one
SELECT * FROM greenprints
WHERE id = $1;

-- name: GetAllRegions :many
SELECT * FROM regions
ORDER BY tree_amount DESC;

-- name: IncreaseRegionTreeAmount :exec
UPDATE regions
SET tree_amount = tree_amount + $1
WHERE id = $2;

-- name: GetRegionById :one
SELECT * FROM regions
WHERE id = $1;

-- name: GetParticipationByMemoryId :one
SELECT COUNT(*) as participations FROM participations WHERE memory_id = $1;

-- name: GetContribution :one
SELECT
  COUNT(*) AS is_exist
FROM contributions
WHERE quest_id = $1 AND user_id = $2;

-- name: GetQuestByCodeId
