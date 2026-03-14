-- name: CreatePointLedger :one
INSERT INTO point_ledger (organization_id, user_id, event_type, event_id, points_delta, description)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
) RETURNING *;
