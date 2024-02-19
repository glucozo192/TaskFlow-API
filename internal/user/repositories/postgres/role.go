package postgres

import (
	"context"
	"fmt"
	"github.com/glu-project/internal/user/models"
	dbutil "github.com/glu-project/utils/database"

	"github.com/jackc/pgx/v5"
)

type RoleRepository struct {
}

func (r *RoleRepository) GetByID(ctx context.Context, db models.DBTX, id string) (*models.Role, error) {
	query := `SELECT * FROM roles WHERE id = $1`
	role, err := dbutil.SelectRow(ctx, db, query, []any{id}, pgx.RowToAddrOfStructByName[models.Role])
	if err != nil {
		return nil, fmt.Errorf("dbutil.SelectRow: %w", err)
	}
	return role, nil
}
