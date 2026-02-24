-- name: CreateService :one
INSERT INTO service(name, code, description, base_price, base_points_rewards)
VALUES(
    $1,
    $2,
    $3,
    $4,
    $5
) RETURNING *;
