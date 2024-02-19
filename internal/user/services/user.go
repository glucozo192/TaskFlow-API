package services

import (
	"context"
	"errors"

	"github.com/glu-project/idl/pb"
	"github.com/glu-project/internal/user/models"
	"github.com/glu-project/internal/user/repositories/postgres"
	"github.com/glu-project/transform"

	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserService struct {
	db models.DBTX
	pb.UnimplementedUserServiceServer

	userRepo interface {
		GetByID(ctx context.Context, db models.DBTX, id string) (*models.User, error)
		Delete(ctx context.Context, db models.DBTX, id string) error
		GetList(ctx context.Context, db models.DBTX, args models.Paging) ([]*models.User, error)
		GetTotalUser(ctx context.Context, db models.DBTX) (int32, error)
		UpdateUser(ctx context.Context, db models.DBTX, user *models.User) error
	}
}

func NewUserService(db models.DBTX) pb.UserServiceServer {
	return &UserService{
		db:       db,
		userRepo: new(postgres.UserRepository),
	}
}
func (s *UserService) GetUserByID(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.GetUserByIDResponse, error) {
	user, err := s.userRepo.GetByID(ctx, s.db, req.GetId())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Errorf(codes.InvalidArgument, "not found any user with ID %s", req.GetId())
		}
		return nil, status.Errorf(codes.Internal, "s.userRepo.GetByID: unexpected error: %v", err)
	}
	return &pb.GetUserByIDResponse{
		Data: transform.UserToPbPtr(user),
	}, nil
}

func (s *UserService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	if err := s.userRepo.Delete(ctx, s.db, req.GetId()); err != nil {
		return nil, status.Errorf(codes.Internal, "s.userRepo.Delete: unexpected error: %v", err)
	}
	return &pb.DeleteUserResponse{}, nil
}

func (s *UserService) GetListUser(ctx context.Context, req *pb.GetListUserRequest) (*pb.GetListUserResponse, error) {
	paging := models.NewPagingWithDefault(req.GetPage(), req.GetPageSize(), req.GetOrderBy(), req.GetOrderType().String(), "")

	users, err := s.userRepo.GetList(ctx, s.db, paging)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.userRepo.GetList: unexpected error: %v", err)
	}
	total, err := s.userRepo.GetTotalUser(ctx, s.db)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.userRepo.GetTotalUser: unexpected error: %v", err)
	}

	return &pb.GetListUserResponse{
		Data:       transform.UserToPbPtrList(users),
		Total:      total,
		TotalPages: paging.CalTotalPages(total),
		Page:       req.GetPage(),
	}, nil
}
