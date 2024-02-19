package entities

import "github.com/jackc/pgtype"

type UserRole struct {
	UserID    pgtype.Text      `db:"user_id"`
	RoleID    pgtype.Text      `db:"role_id"`
	CreatedAt pgtype.Timestamp `db:"created_at"`
	UpdatedAt pgtype.Timestamp `db:"updated_at"`
}

func (u *UserRole) TableName() string {
	return "user_roles"
}
