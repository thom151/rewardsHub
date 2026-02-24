-- name: CreateProperty :one
INSERT INTO property(organization_id, created_by_user_id, address_line1, address_line2, city, state_region, postal_code, listing_url)
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

