-- name: CreateRedemption :one
INSERT INTO redemption (
organization_id,
requested_by_user_id,
reward_id,
applied_booking_id,
points_cost,
status
)
VALUES (
$1,
$2,
$3,
$4,
$5,
'pending'
)
RETURNING *;


-- name: GetRedemptionFromID :one
SELECT * FROM redemption WHERE redemption_id = $1;


-- name: ApproveRedemption :one
UPDATE redemption
SET
    status = 'approved',
    approved_at = NOW()
WHERE redemption_id = $1
RETURNING *;

