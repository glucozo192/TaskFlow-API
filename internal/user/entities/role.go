package entities

import "github.com/jackc/pgtype"

type Role struct {
	ID         pgtype.Text        `db:"role_id"`
	Name       pgtype.Text        `db:"name"`
	InsertedAt pgtype.Timestamptz `db:"created_at"`
	UpdatedAt  pgtype.Timestamptz `db:"updated_at"`
}

func (r *Role) TableName() string {
	return "roles"
}
