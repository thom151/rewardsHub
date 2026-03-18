-- name: CreateJob :one
INSERT INTO job (booking_id, assigned_to_user_id, dropbox_folder_path)
VALUES (
    $1,
    $2,
    $3
) RETURNING *;

-- name: InsertEventID :one
UPDATE job
SET google_event_id = $2
WHERE job_id = $1
RETURNING *;


-- name: GetJob :one
SELECT * FROM job WHERE job_id = $1;
