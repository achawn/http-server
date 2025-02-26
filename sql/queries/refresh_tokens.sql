-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens(token, created_at, updated_at, user_id, expires_at)
VALUES (
	$1,
	NOW(),
	NOW(),
	$2,
	$3
)
RETURNING *;

-- name: GetUserFromToken :one
SELECT u.* FROM refresh_tokens rt
INNER JOIN users u ON
rt.user_id = u.id
WHERE token = $1
AND revoked_at IS NULL
AND expires_at > NOW();

-- name: RevokeToken :one
UPDATE refresh_tokens
SET revoked_at = NOW(),
updated_at = NOW()
WHERE token = $1
RETURNING *;
