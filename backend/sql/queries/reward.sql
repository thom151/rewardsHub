-- name: CreateReward :one
INSERT INTO reward(code, name, description, service_id, points_cost)
VALUES(
    $1,
    $2,
    $3,
    $4,
    $5
) RETURNING *;

-- name: GetRewardFromID :one
SELECT * FROM reward WHERE reward_id = $1;

-- name: GetRewardFromServiceID :one
SELECT * FROM reward WHERE service_id = $1;
