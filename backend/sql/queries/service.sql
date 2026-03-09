-- name: CreateService :one
INSERT INTO service(name, code, description, base_price_cents, base_points_rewards)
VALUES(
    $1,
    $2,
    $3,
    $4,
    $5
) RETURNING *;

-- name: GetService :one
SELECT * FROM service WHERE service_id = $1;

-- name: GetAllServices :many
SELECT * FROM service;
