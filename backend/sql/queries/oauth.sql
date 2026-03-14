-- name: CreateOauthConnection :one
INSERT INTO user_oauth_connections (user_id,
    provider,
    service, 
    provider_email, 
    refresh_token, 
    access_token, 
    expires_at, 
    scopes)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
) RETURNING *;


-- name: GetOauthConnectionFromUser :one
SELECT * FROM user_oauth_connections WHERE 
user_id = $1 AND provider = $2 AND service = $3;


-- name: UpdateOauthConnectionTokens :one
UPDATE user_oauth_connections
SET
    refresh_token = $4,
    access_token = $5,
    expires_at = $6,
    scopes = $7,
    updated_at = NOW()
WHERE user_id = $1
  AND provider = $2
  AND service = $3
RETURNING *;
