-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: ResetUsers :exec
DELETE FROM users;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: UpdateUserWithID :one
UPDATE users
SET email = $1, hashed_password = $2, updated_at = $3
WHERE id = $4
RETURNING *;

-- name: UpgradeUser :one
UPDATE users
SET is_chirpy_red = $1, updated_at = NOW()
WHERE id = $2
RETURNING *;