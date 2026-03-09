-- name: CreateOAuthState :one
INSERT INTO oauth_states(state, user_id, expires_at)
VALUES(
    $1,
    $2,
    $3
)RETURNING *;

-- name: GetUserFromState :one
SELECT * FROM oauth_states WHERE state = $1;
