-- name: CreateCart :one
INSERT INTO carts (user_id)
VALUES ($1)
RETURNING id, user_id, created_at;

-- name: AddToCart :exec
INSERT INTO cart_items (cart_id, product_id, quantity)
VALUES ($1, $2, $3)
ON CONFLICT (cart_id, product_id)
    DO UPDATE SET quantity = cart_items.quantity + EXCLUDED.quantity;

-- name: GetCartByUser :one
SELECT id, user_id, created_at
FROM carts
WHERE user_id = $1
LIMIT 1;

-- name: GetCart :many
SELECT ci.id, ci.cart_id, ci.product_id, ci.quantity, ci.added_at,
       p.name AS product_name, p.price AS product_price
FROM cart_items ci
         JOIN products p ON ci.product_id = p.id
WHERE ci.cart_id = $1;

-- name: ClearCart :exec
DELETE FROM cart_items WHERE cart_id = $1;
