-- name: CreateOrder :one
INSERT INTO orders (user_id, total_price, status)
VALUES ($1, $2, 'pending')
RETURNING id, created_at;

-- name: UpdateOrder :exec
UPDATE orders
SET total_price = $1, status = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $3 AND status NOT IN ('shipped', 'completed', 'cancelled');

-- name: CancelOrder :exec
UPDATE orders
SET status = 'cancelled', updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND status = 'pending';

-- name: SettleOrder :exec
UPDATE orders
SET status = 'paid', updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND status = 'pending';
