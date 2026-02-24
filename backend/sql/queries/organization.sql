-- name: CreateOrganization :one
INSERT INTO organization (name, organization_type)
VALUES (
    $1,
    $2
) RETURNING *;


-- name: GetOrganizationFromID :one
SELECT * FROM organization WHERE organization_id = $1;
