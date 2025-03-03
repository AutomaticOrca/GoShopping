-- name: GetOrderItemsByOrderID :many
SELECT product_id, quantity, price
FROM order_items
WHERE order_id = $1;

-- name: CreateOrderItem :exec
INSERT INTO order_items (order_id, product_id, quantity, price, subtotal)
VALUES ($1, $2, $3, $4, $5);