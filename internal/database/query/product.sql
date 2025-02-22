-- name: CreateProduct :one
INSERT INTO products (name, description, price, stock, image_url, category_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdateProduct :one
UPDATE products
SET name = COALESCE($2, name),
    description = COALESCE($3, description),
    price = COALESCE($4, price),
    stock = COALESCE($5, stock),
    image_url = COALESCE($6, image_url),
    category_id = COALESCE($7, category_id),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM products WHERE id = $1;

-- name: GetProductByID :one
SELECT * FROM products WHERE id = $1;

-- name: ListProducts :many
SELECT * FROM products
LIMIT $1 OFFSET $2;
