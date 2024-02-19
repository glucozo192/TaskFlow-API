package repositories

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/glu/shopvui/internal/userm/entities"
	"github.com/glu/shopvui/internal/userm/golibs/database"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgxutil"
	"go.uber.org/multierr"
)

type UserRepo struct {
}

const getUserStmTpl = `
	select %s from %s where email = $1
`

func (r *UserRepo) GetUser(ctx context.Context, db database.Ext, email pgtype.Text) (*entities.User, error) {
	e := &entities.User{}
	fields, values := database.FieldMap(e)
	err := db.QueryRow(ctx, fmt.Sprintf(getUserStmTpl, strings.Join(fields, ","), "users"), email).Scan(values...)
	if err != nil {
		return nil, fmt.Errorf("db.QueryRow.Scan: %w", err)
	}
	return e, nil
}

func (r *UserRepo) GetUserByID(ctx context.Context, db database.Ext, userID pgtype.Text) (*entities.User, error) {
	//e := &entities.User{}
	//fields, values := database.FieldMap(e)
	query := fmt.Sprintf(`select * from %s where user_id = $1`, "users")
	//err := db.QueryRow(ctx, query, userID).Scan(values...)

	// rows, err := dbpool.Query(ctx, "SELECT * FROM accounts")
	// if err != nil {
	// 	exit("failed query", err)
	// }

	user, err := pgxutil.SelectRow(ctx, db, query, []any{userID}, pgx.RowToAddrOfStructByPos[entities.User])
	if err != nil {
		return nil, fmt.Errorf("GetUserByID: db.QueryRow.Scan: %w", err)
	}
	return user, nil
}

func (r *UserRepo) CreateUser(ctx context.Context, db database.Ext, e *entities.User) (*entities.User, error) {
	fieldNames, values := database.FieldMap(e)
	placeHolders := database.GeneratePlaceholders(len(fieldNames))
	query := fmt.Sprintf(`
		INSERT INTO users (%s) VALUES (%s)
		ON CONFLICT ON CONSTRAINT pk_users DO UPDATE SET
			first_name = excluded.first_name,
			last_name = excluded.last_name,
            password = excluded.password,
			updated_at = excluded.updated_at
		`, strings.Join(fieldNames, ", "), placeHolders)

	now := time.Now()
	err := multierr.Combine(
		e.InsertedAt.Set(now),
		e.UpdatedAt.Set(now),
		e.Active.Set(true),
	)
	if err != nil {
		return &entities.User{}, fmt.Errorf("multierr.Combine: %w", err)
	}
	_, err = db.Exec(ctx, query, values...)
	if err != nil {
		return &entities.User{}, fmt.Errorf("db.Exec: %w", err)
	}
	return e, nil
}

func (r *UserRepo) AddRoles(ctx context.Context, db database.Ext, roles *entities.Role) error {
	fields, values := database.FieldMap(roles)
	placeHolders := database.GeneratePlaceholders(len(fields))
	query := fmt.Sprintf(`INSERT INTO roles (%s) VALUES (%s)`, strings.Join(fields, ", "), placeHolders)
	fmt.Println(query)
	now := time.Now()
	err := multierr.Combine(
		roles.InsertedAt.Set(now),
		roles.UpdatedAt.Set(now),
	)
	if err != nil {
		return fmt.Errorf("multierr.Combine: %w", err)
	}
	_, err = db.Exec(ctx, query, values...)
	if err != nil {
		return fmt.Errorf("db.Exec: %w", err)
	}
	return nil
}

func (r *UserRepo) GetRole(ctx context.Context, db database.Ext, roleName pgtype.Text) (*entities.Role, error) {
	e := &entities.Role{}
	fields, values := database.FieldMap(e)
	query := fmt.Sprintf(`SELECT %s FROM %s WHERE name = $1`, strings.Join(fields, ","), e.TableName())
	err := db.QueryRow(ctx, query, roleName).Scan(values...)
	if err != nil {
		return nil, fmt.Errorf("db.QueryRow.Scan: %w", err)
	}
	return e, nil
}

func (r *UserRepo) UpdateRole(ctx context.Context, db database.Ext, e *entities.UserRole) (*entities.UserRole, error) {
	fields, values := database.FieldMap(e)
	placeHolders := database.GeneratePlaceholders(len(fields))
	query := fmt.Sprintf(`
		INSERT INTO %s (%s) VALUES (%s)
	`, e.TableName(), strings.Join(fields, ","), placeHolders)
	now := time.Now()
	err := multierr.Combine(
		e.CreatedAt.Set(now),
		e.UpdatedAt.Set(now),
	)
	if err != nil {
		return &entities.UserRole{}, fmt.Errorf("multierr.Combine: %w", err)
	}
	_, err = db.Exec(ctx, query, values...)
	if err != nil {
		return &entities.UserRole{}, fmt.Errorf("db.Exec: %w", err)
	}
	return e, nil
}

// func (r *UserRepo) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
// 	q := models.New(r.db)
// 	result, err := q.GetUserByUsername(ctx, username)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return result, nil
// }
