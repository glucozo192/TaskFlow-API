package postgres

import (
	"context"
	"fmt"

	"github.com/glu-project/internal/taks/constants"
	"github.com/glu-project/internal/taks/models"
	"github.com/glu-project/utils"
	dbutil "github.com/glu-project/utils/database"

	"github.com/jackc/pgx/v5"
)

type TaskRepository struct{}

// Create inserts a new task and returns the generated ID.
func (r *TaskRepository) Create(ctx context.Context, db models.DBTX, task *models.Task) (string, error) {
	m, err := utils.FieldsByDBTag(task)
	if err != nil {
		return "", err
	}
	id, err := dbutil.InsertRowReturning(ctx, db, constants.TBN_Tasks, m, "id", pgx.RowTo[string])
	if err != nil {
		return "", fmt.Errorf("dbutil.InsertRowReturning: %w", err)
	}
	return id, nil
}

// GetByID retrieves a single task by its primary key.
func (r *TaskRepository) GetByID(ctx context.Context, db models.DBTX, id string) (*models.Task, error) {
	const query = `SELECT * FROM tasks WHERE id = $1 AND deleted_at IS NULL`
	task, err := dbutil.SelectRow(ctx, db, query, []any{id}, pgx.RowToAddrOfStructByName[models.Task])
	if err != nil {
		return nil, fmt.Errorf("dbutil.SelectRow: %w", err)
	}
	return task, nil
}

// GetList returns a paginated list of tasks belonging to a specific user.
func (r *TaskRepository) GetList(ctx context.Context, db models.DBTX, userID string, args models.Paging) ([]*models.Task, error) {
	query := fmt.Sprintf(`
	SELECT *
	FROM tasks
	WHERE user_id = $1 AND deleted_at IS NULL
	ORDER BY %s %s
	LIMIT $2 OFFSET $3
	`, args.GetOrderBy(), args.GetOrderType())

	tasks, err := dbutil.Select(ctx, db, query, []any{userID, args.GetLimit(), args.GetOffset()}, pgx.RowToAddrOfStructByName[models.Task])
	if err != nil {
		return nil, fmt.Errorf("dbutil.Select: %w", err)
	}
	return tasks, nil
}

// GetTotal counts all non-deleted tasks for a given user.
func (r *TaskRepository) GetTotal(ctx context.Context, db models.DBTX, userID string) (int32, error) {
	const query = `SELECT count(*) FROM tasks WHERE user_id = $1 AND deleted_at IS NULL`
	total, err := dbutil.SelectRow(ctx, db, query, []any{userID}, pgx.RowTo[int32])
	if err != nil {
		return 0, fmt.Errorf("dbutil.SelectRow: %w", err)
	}
	return total, nil
}

// Update patches the mutable fields of a task.
func (r *TaskRepository) Update(ctx context.Context, db models.DBTX, task *models.Task) error {
	const query = `
	UPDATE tasks SET
		title       = coalesce($2, title),
		description = coalesce($3, description),
		status      = coalesce($4, status),
		updated_at  = NOW()
	WHERE id = $1 AND deleted_at IS NULL
	`
	if _, err := db.Exec(ctx, query,
		task.ID,
		task.Title,
		task.Description,
		task.Status,
	); err != nil {
		return fmt.Errorf("db.Exec Update: %w", err)
	}
	return nil
}

// Delete soft-deletes a task by setting its deleted_at timestamp.
func (r *TaskRepository) Delete(ctx context.Context, db models.DBTX, id string) error {
	const query = `UPDATE tasks SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	if _, err := db.Exec(ctx, query, id); err != nil {
		return fmt.Errorf("db.Exec Delete: %w", err)
	}
	return nil
}
