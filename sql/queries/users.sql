-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, password)
VALUES (
  gen_random_uuid(),
  NOW(),
  NOW(),
  $1,
  $2
)
RETURNING *;


-- name: GetUser :one
SELECT * FROM users WHERE email = $1;

-- name: UpdateUser :one
UPDATE users
SET password = $2, email = $3, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpgradeUserToChirpyRed :one
UPDATE users
SET is_chirpy_red = TRUE
WHERE id = $1
RETURNING *;