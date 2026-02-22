-- name: SetPasswordForUser :one
INSERT INTO auth_identity (user_id, provider, provider_subject, provider_hash)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;



-- name: GetUserPasswordHash :one
SELECT * FROM auth_identity WHERE user_id = $1;
