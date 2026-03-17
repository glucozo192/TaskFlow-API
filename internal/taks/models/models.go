package models

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type DBTX interface {
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
}

// Task represents a todo task record in the database.
type Task struct {
	ID          string             `db:"id"          json:"id"`
	Title       pgtype.Text        `db:"title"        json:"title"`
	Description pgtype.Text        `db:"description"  json:"description"`
	Status      pgtype.Text        `db:"status"       json:"status"`
	UserID      pgtype.Text        `db:"user_id"      json:"user_id"`
	CreatedAt   pgtype.Timestamptz `db:"created_at"   json:"created_at"`
	UpdatedAt   pgtype.Timestamptz `db:"updated_at"   json:"updated_at"`
	DeletedAt   pgtype.Timestamptz `db:"deleted_at"   json:"deleted_at"`
}
