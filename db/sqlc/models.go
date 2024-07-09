// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Account struct {
	ID        int64
	Owner     string
	Balance   int64
	Currency  string
	CreatedAt pgtype.Timestamptz
}

type Entry struct {
	ID        int64
	AccountID pgtype.Int8
	Amount    int64
	CreatedAt pgtype.Timestamptz
}

type Session struct {
	ID           pgtype.UUID
	Username     string
	RefreshToken string
	UserAgent    string
	ClientIp     string
	ExpiredAt    pgtype.Timestamptz
	CreatedAt    pgtype.Timestamptz
}

type Transfer struct {
	ID          int64
	FromAccount pgtype.Int8
	ToAccount   pgtype.Int8
	// must >0
	Amount    int64
	CreatedAt pgtype.Timestamptz
}

type User struct {
	Username          string
	Email             string
	Password          string
	FullName          string
	PasswordExpiredAt pgtype.Timestamptz
	CreatedAt         pgtype.Timestamptz
}
