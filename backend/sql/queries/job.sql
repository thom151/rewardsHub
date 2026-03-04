-- name: CreateJob :one
INSERT INTO job (booking_id, assigned_to_user_id)
VALUES (
    $1,
    $2
) RETURNING *;
