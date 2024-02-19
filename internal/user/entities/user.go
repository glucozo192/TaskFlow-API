package entities

import "github.com/jackc/pgtype"

type User struct {
	ID         pgtype.Text        `db:"user_id"`
	Email      pgtype.Text        `db:"email"`
	Password   pgtype.Text        `db:"password"`
	LastName   pgtype.Text        `db:"last_name"`
	FirstName  pgtype.Text        `db:"first_name"`
	Active     pgtype.Bool        `db:"active"`
	InsertedAt pgtype.Timestamptz `db:"created_at"`
	UpdatedAt  pgtype.Timestamptz `db:"updated_at"`
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		"id",
		"email",
		"password",
		"last_name",
		"first_name",
		"active",
		"inserted_at",
		"updated_at",
	}
	values = []interface{}{
		&u.ID,
		&u.Email,
		&u.Password,
		&u.LastName,
		&u.FirstName,
		&u.Active,
		&u.InsertedAt,
		&u.UpdatedAt,
	}
	return
}
