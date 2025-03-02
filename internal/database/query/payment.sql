-- name: CancelPayment :exec
UPDATE payments
SET status = 'cancelled', updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND status = 'pending';

-- name: TimeoutPayment :exec
UPDATE payments
SET status = 'cancelled', updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND status = 'pending' AND created_at < NOW() - INTERVAL '15 minutes';

-- name: CreatePayment :one
INSERT INTO payments (order_id, user_id, amount, status, transaction_id)
VALUES ($1, $2, $3, 'pending', $4)
RETURNING id, created_at;

-- name: ConfirmPayment :exec
UPDATE payments
SET status = 'success', updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND status = 'pending';
