-- name: CreateAuthToken :one
INSERT INTO auth_tokens (user_id, token, expires_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ValidateAuthToken :one
SELECT * FROM auth_tokens
WHERE user_id = $1 AND token = $2 AND expires_at > NOW();

-- name: DeleteAuthToken :exec
DELETE FROM auth_tokens WHERE user_id = $1 AND token = $2;

-- name: DeleteAllAuthTokens :exec
DELETE FROM auth_tokens WHERE user_id = $1;