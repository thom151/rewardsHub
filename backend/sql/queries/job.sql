-- name: CreateJob :one
INSERT INTO job (booking_id, assigned_to_user_id)
VALUES (
    $1,
    $2
) RETURNING *;

-- name: InsertEventID :one
UPDATE job
SET google_event_id = $2
WHERE job_id = $1
RETURNING *;
