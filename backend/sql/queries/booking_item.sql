-- name: CreateBookingItem :one
INSERT INTO booking_item(booking_id, service_id, quantity, unit_price_cents, points_award)
VALUES(
    $1,
    $2,
    $3,
    $4,
    $5
)RETURNING *;

-- name: GetBookingItems :many
SELECT * FROM booking_item WHERE booking_id=$1;
