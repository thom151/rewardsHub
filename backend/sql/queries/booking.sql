-- name: CreateBooking :one
INSERT INTO booking (organization_id, requested_by_user_id, property_id, preferred_date, schedule_start, schedule_end, status)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
) RETURNING *;

-- name: GetBooking :one
SELECT * FROM booking WHERE booking_id = $1;

-- name: ConfirmBooking :one
UPDATE booking SET status='confirmed',
    updated_at = NOW()
WHERE booking_id = $1
RETURNING *;
