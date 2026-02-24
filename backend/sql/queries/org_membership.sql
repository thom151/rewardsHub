-- name: CreateOrgMembership :one
INSERT INTO org_membership (organization_id, user_id, org_role)
VALUES(
    $1,
    $2,
    $3
)
RETURNING *;
