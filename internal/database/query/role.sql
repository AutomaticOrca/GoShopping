-- name: CreateRole :one
INSERT INTO roles (name, description)
VALUES ($1, $2)
RETURNING *;

-- name: ListRoles :many
SELECT id, name, description FROM roles;

-- name: GetRoleByName :one
SELECT id, name, description FROM roles WHERE name = $1;