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
