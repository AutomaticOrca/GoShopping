-- name: CreateUser :one
INSERT INTO users (email,
                   password_hash)
VALUES ($1, $2)
RETURNING *;

-- name: GetUserByID :one
SELECT (
        id, email, full_name, avatar_url, role
           )
FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1;

-- name: DeleteUser :exec
DELETE
FROM users
WHERE id = $1;
