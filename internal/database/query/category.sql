-- name: CreateCategory :one
INSERT INTO categories (name, parent_id)
VALUES ($1, $2)
RETURNING *;

-- name: GetCategoryByID :one
SELECT * FROM categories WHERE id = $1 LIMIT 1;

-- name: ListCategories :many
SELECT * FROM categories ORDER BY id ASC;

-- name: GetSubcategories :many
SELECT * FROM categories WHERE parent_id = $1 ORDER BY id ASC;

-- name: UpdateCategory :one
UPDATE categories
SET name = COALESCE($2, name),
    parent_id = COALESCE($3, parent_id)
WHERE id = $1
RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM categories WHERE id = $1;
