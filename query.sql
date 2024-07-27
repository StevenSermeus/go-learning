-- name: GetUser :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1;

-- name: CreateUser :one
INSERT INTO users (application_id, username, pass, last_login) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: UpdateUser :one
UPDATE users SET application_id = $1, username = $2, pass = $3, last_login = $4 WHERE id = $5 RETURNING *;

-- name: DeleteUser :one
DELETE FROM users WHERE id = $1 RETURNING *;

-- name: GetApplication :one
SELECT * FROM application WHERE id = $1;

-- name: CreateApplication :one
INSERT INTO application (name, pass_phrase, api_key, access_token_duration, refresh_token_duration) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetApplicationByAPIKey :one
SELECT * FROM application WHERE api_key = $1;