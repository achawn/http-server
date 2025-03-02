-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
	gen_random_uuid(),
	now(),
	now(),
	$1,
	$2
)
RETURNING *;

-- name: GetUsers :many
SELECT * FROM users
ORDER BY created_at;

-- name: GetUser :one
SELECT * FROM users
WHERE email = $1;

-- name: UpdateUser :one
UPDATE users
SET email = $2,
hashed_password = $3
WHERE id = $1
RETURNING *;
