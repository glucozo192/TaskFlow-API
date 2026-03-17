package services

import (
	"context"
	"errors"

	"github.com/glu-project/idl/pb"
	"github.com/glu-project/internal/taks/models"
	esrepo "github.com/glu-project/internal/taks/repositories/elasticsearch"
	pgrepo "github.com/glu-project/internal/taks/repositories/postgres"
	elastic_client "github.com/glu-project/pkg/elasticsearch_client"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TaskService implements pb.TaskServiceServer and orchestrates
// task CRUD logic across Postgres and Elasticsearch.
type TaskService struct {
	pb.UnimplementedTaskServiceServer

	db       models.DBTX
	esClient *elastic_client.ElasticClient

	taskRepo interface {
		Create(ctx context.Context, db models.DBTX, task *models.Task) (string, error)
		GetByID(ctx context.Context, db models.DBTX, id string) (*models.Task, error)
		GetList(ctx context.Context, db models.DBTX, userID string, args models.Paging) ([]*models.Task, error)
		GetTotal(ctx context.Context, db models.DBTX, userID string) (int32, error)
		Update(ctx context.Context, db models.DBTX, task *models.Task) error
		Delete(ctx context.Context, db models.DBTX, id string) error
	}
	taskESRepo interface {
		IndexTask(ctx context.Context, client *elastic_client.ElasticClient, task *models.Task) error
		SearchTasks(ctx context.Context, client *elastic_client.ElasticClient, q string) ([]*models.Task, error)
	}
}

// NewTaskService constructs a TaskService and returns it as the gRPC server interface.
// Pass nil for esClient to disable Elasticsearch indexing.
func NewTaskService(db models.DBTX, esClient *elastic_client.ElasticClient) pb.TaskServiceServer {
	return &TaskService{
		db:         db,
		esClient:   esClient,
		taskRepo:   new(pgrepo.TaskRepository),
		taskESRepo: new(esrepo.TaskESRepository),
	}
}

// ─── helpers ────────────────────────────────────────────────────────────────

func taskToPb(t *models.Task) *pb.Task {
	if t == nil {
		return nil
	}
	return &pb.Task{
		Id:          t.ID,
		Title:       t.Title.String,
		Description: t.Description.String,
		Status:      pb.TaskStatus(pb.TaskStatus_value[t.Status.String]),
		UserId:      t.UserID.String,
	}
}

func tasksToPb(tasks []*models.Task) []*pb.Task {
	out := make([]*pb.Task, 0, len(tasks))
	for _, t := range tasks {
		out = append(out, taskToPb(t))
	}
	return out
}

// ─── gRPC handlers ──────────────────────────────────────────────────────────

// CreateTask creates a new task for the authenticated user.
func (s *TaskService) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
	task := &models.Task{
		ID:    uuid.NewString(),
		Title: pgtype.Text{String: req.GetTitle(), Valid: req.GetTitle() != ""},
		Description: pgtype.Text{
			String: req.GetDescription(),
			Valid:  req.GetDescription() != "",
		},
		Status: pgtype.Text{String: req.GetStatus().String(), Valid: true},
		UserID: pgtype.Text{String: req.GetUserId(), Valid: req.GetUserId() != ""},
	}

	id, err := s.taskRepo.Create(ctx, s.db, task)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.taskRepo.Create: unexpected error: %v", err)
	}
	task.ID = id

	// Best-effort Elasticsearch indexing — do not fail the request if ES is unavailable.
	if s.esClient != nil {
		if esErr := s.taskESRepo.IndexTask(ctx, s.esClient, task); esErr != nil {
			// Log the error in production; for now we swallow it.
			_ = esErr
		}
	}

	return &pb.CreateTaskResponse{Id: id}, nil
}

// GetTaskByID retrieves a single task by its ID.
func (s *TaskService) GetTaskByID(ctx context.Context, req *pb.GetTaskByIDRequest) (*pb.GetTaskByIDResponse, error) {
	task, err := s.taskRepo.GetByID(ctx, s.db, req.GetId())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "task with ID %s not found", req.GetId())
		}
		return nil, status.Errorf(codes.Internal, "s.taskRepo.GetByID: unexpected error: %v", err)
	}
	return &pb.GetTaskByIDResponse{Data: taskToPb(task)}, nil
}

// GetListTask returns a paginated list of tasks for the given user.
func (s *TaskService) GetListTask(ctx context.Context, req *pb.GetListTaskRequest) (*pb.GetListTaskResponse, error) {
	paging := models.NewPagingWithDefault(
		req.GetPage(),
		req.GetPageSize(),
		req.GetOrderBy(),
		req.GetOrderType().String(),
		"",
	)

	tasks, err := s.taskRepo.GetList(ctx, s.db, req.GetUserId(), paging)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.taskRepo.GetList: unexpected error: %v", err)
	}

	total, err := s.taskRepo.GetTotal(ctx, s.db, req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.taskRepo.GetTotal: unexpected error: %v", err)
	}

	return &pb.GetListTaskResponse{
		Data:       tasksToPb(tasks),
		Total:      total,
		TotalPages: paging.CalTotalPages(total),
		Page:       req.GetPage(),
	}, nil
}

// UpdateTask modifies mutable fields of an existing task.
func (s *TaskService) UpdateTask(ctx context.Context, req *pb.UpdateTaskRequest) (*pb.UpdateTaskResponse, error) {
	task := &models.Task{
		ID:          req.GetId(),
		Title:       pgtype.Text{String: req.GetTitle(), Valid: req.GetTitle() != ""},
		Description: pgtype.Text{String: req.GetDescription(), Valid: req.GetDescription() != ""},
		Status:      pgtype.Text{String: req.GetStatus().String(), Valid: req.GetStatus() != pb.TaskStatus_TaskStatus_NONE},
	}

	if err := s.taskRepo.Update(ctx, s.db, task); err != nil {
		return nil, status.Errorf(codes.Internal, "s.taskRepo.Update: unexpected error: %v", err)
	}

	// Re-index the updated task in Elasticsearch.
	if s.esClient != nil {
		if updated, getErr := s.taskRepo.GetByID(ctx, s.db, req.GetId()); getErr == nil {
			_ = s.taskESRepo.IndexTask(ctx, s.esClient, updated)
		}
	}

	return &pb.UpdateTaskResponse{}, nil
}

// DeleteTask soft-deletes a task by ID.
func (s *TaskService) DeleteTask(ctx context.Context, req *pb.DeleteTaskRequest) (*pb.DeleteTaskResponse, error) {
	if err := s.taskRepo.Delete(ctx, s.db, req.GetId()); err != nil {
		return nil, status.Errorf(codes.Internal, "s.taskRepo.Delete: unexpected error: %v", err)
	}
	return &pb.DeleteTaskResponse{}, nil
}
