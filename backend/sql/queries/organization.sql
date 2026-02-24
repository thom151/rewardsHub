-- name: CreateOrganization :one
INSERT INTO organization (name, organization_type)
VALUES (
    $1,
    $2
) RETURNING *;
