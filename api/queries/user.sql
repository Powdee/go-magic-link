-- Inserts a new user or updates an existing one with the new magic link and username details
-- name: UpsertUserWithToken :exec
INSERT INTO users (email, username, magic_token, token_expiration)
VALUES ($1, $2, $3, $4)
ON CONFLICT (email) DO UPDATE SET
    username = COALESCE(EXCLUDED.username, users.username),  -- Only update username if EXCLUDED.username is not null
    magic_token = EXCLUDED.magic_token,
    token_expiration = EXCLUDED.token_expiration;

-- Retrieves a user by their magic token if the token has not expired
-- name: GetUserByToken :one
SELECT id, email, username, magic_token, token_expiration FROM users
WHERE magic_token = $1 AND token_expiration > NOW();
