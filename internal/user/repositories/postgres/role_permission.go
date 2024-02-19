package postgres

import (
	"context"

	"github.com/glu-project/internal/user/models"
)

type RolePermissionRepository struct {
}

func (r *RolePermissionRepository) GetListRolePermissions(ctx context.Context, db models.DBTX) ([]*models.GetListPermissionsRow, error) {
	const getListPermissions = `
SELECT
role_id, ARRAY_AGG(path)::text[] as permissions
FROM role_permissions
WHERE deleted_at IS NULL
GROUP BY role_id
`
	rows, err := db.Query(ctx, getListPermissions)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*models.GetListPermissionsRow
	for rows.Next() {
		var i models.GetListPermissionsRow
		if err := rows.Scan(&i.RoleID, &i.Permissions); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
