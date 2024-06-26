// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: user.sql

package db

import (
	"context"
	"database/sql"
)

const getUserByToken = `-- name: GetUserByToken :one
SELECT id, email, username, magic_token, token_expiration FROM users
WHERE magic_token = $1 AND token_expiration > NOW()
`

// Retrieves a user by their magic token if the token has not expired
func (q *Queries) GetUserByToken(ctx context.Context, magicToken sql.NullString) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByToken, magicToken)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Username,
		&i.MagicToken,
		&i.TokenExpiration,
	)
	return i, err
}

const upsertUserWithToken = `-- name: UpsertUserWithToken :exec
INSERT INTO users (email, username, magic_token, token_expiration)
VALUES ($1, $2, $3, $4)
ON CONFLICT (email) DO UPDATE SET
    username = COALESCE(EXCLUDED.username, users.username),  -- Only update username if EXCLUDED.username is not null
    magic_token = EXCLUDED.magic_token,
    token_expiration = EXCLUDED.token_expiration
`

type UpsertUserWithTokenParams struct {
	Email           string         `json:"email"`
	Username        sql.NullString `json:"username"`
	MagicToken      sql.NullString `json:"magic_token"`
	TokenExpiration sql.NullTime   `json:"token_expiration"`
}

// Inserts a new user or updates an existing one with the new magic link and username details
func (q *Queries) UpsertUserWithToken(ctx context.Context, arg UpsertUserWithTokenParams) error {
	_, err := q.db.ExecContext(ctx, upsertUserWithToken,
		arg.Email,
		arg.Username,
		arg.MagicToken,
		arg.TokenExpiration,
	)
	return err
}
