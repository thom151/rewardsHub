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


-- name: GetOrganizationPointsBalance :one
SELECT COALESCE(SUM(points_delta), 0)::BIGINT AS points_balance
FROM point_ledger
WHERE organization_id = $1;


-- name: CompleteJob :one
UPDATE job
SET
status = 'completed',
completed_at = NOW(),
updated_at = NOW()
WHERE job_id = $1
RETURNING *;

