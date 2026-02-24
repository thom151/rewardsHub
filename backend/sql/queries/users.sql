-- name: CreateUser :one
INSERT INTO users (email, first_name, last_name, phone)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;


-- name: DeleteUser :one
DELETE FROM users WHERE user_id = $1 RETURNING *;



-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserByID :one
SELECT * FROM users WHERE user_id = $1;
