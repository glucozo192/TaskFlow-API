package postgres

import (
	"context"
	"fmt"

	"github.com/glu-project/internal/user/constants"
	"github.com/glu-project/internal/user/models"
	"github.com/glu-project/utils"
	dbutil "github.com/glu-project/utils/database"

	"github.com/jackc/pgx/v5"
)

type UserRepository struct{}

func (r *UserRepository) GetByID(ctx context.Context, db models.DBTX, id string) (*models.User, error) {
	const query = `SELECT * FROM users WHERE id = $1`
	user, err := dbutil.SelectRow(ctx, db, query, []any{id}, pgx.RowToAddrOfStructByName[models.User])
	if err != nil {
		return nil, fmt.Errorf("dbutil.SelectRow: %w", err)
	}
	return user, nil
}

func (r *UserRepository) Create(ctx context.Context, db models.DBTX, user *models.User) (string, error) {
	m, err := utils.FieldsByDBTag(user)
	if err != nil {
		return "", err
	}
	userID, err := dbutil.InsertRowReturning(ctx, db, constants.TBN_Users, m, "id", pgx.RowTo[string])
	if err != nil {
		return "", fmt.Errorf("dbutil.InsertRowReturning: %w", err)
	}
	return userID, nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, db models.DBTX, username string) (*models.User, error) {
	const query = `SELECT * FROM users WHERE username = $1`
	user, err := dbutil.SelectRow(ctx, db, query, []any{username}, pgx.RowToAddrOfStructByName[models.User])
	if err != nil {
		return nil, fmt.Errorf("dbutil.SelectRow: %w", err)
	}
	return user, nil
}

func (r *UserRepository) Delete(ctx context.Context, db models.DBTX, id string) error {
	const query = `
	DELETE FROM users WHERE id = $1
	`
	if _, err := db.Exec(ctx, query, id); err != nil {
		return fmt.Errorf("db.Exec: %w", err)
	}
	return nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, db models.DBTX, user *models.User) error {
	const query = `
  UPDATE users SET 
	username = coalesce($2, username),
  email = coalesce($3, email),
  password = coalesce($4, password),
	picture = coalesce($5, picture),
	point = coalesce($6, point),
	is_enable_two_fa = coalesce($7, is_enable_two_fa),
  updated_at = NOW()
  WHERE id = $1
  `
	if _, err := db.Exec(ctx, query,
		user.ID,
		user.Username,
		user.Email,
		user.Email,
		user.Password); err != nil {
		return fmt.Errorf("db.Exec: %w", err)
	}
	return nil
}
func (r *UserRepository) GetList(ctx context.Context, db models.DBTX, args models.Paging) ([]*models.User, error) {
	query := `
	SELECT *
	FROM
		users
	ORDER BY %s %s
	LIMIT $1 OFFSET $2
	`
	users, err := dbutil.Select(ctx, db, fmt.Sprintf(query, args.GetOrderBy(), args.GetOrderType()), []any{args.GetLimit(), args.GetOffset()}, pgx.RowToAddrOfStructByName[models.User])
	if err != nil {
		return nil, fmt.Errorf("dbutil.Select: %w", err)
	}
	return users, nil
}

func (r *UserRepository) GetTotalUser(ctx context.Context, db models.DBTX) (int32, error) {
	query := `
	SELECT
		count(*)
	FROM user
	`
	total, err := dbutil.SelectRow(ctx, db, query, []any{}, pgx.RowTo[int32])
	if err != nil {
		return 0, fmt.Errorf("dbutil.SelectRow: %w", err)
	}

	return total, nil
}
