package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/glu-project/internal/user/models"
	"github.com/glu-project/transform"
	"github.com/glu-project/transformhelpers"
	"github.com/glu-project/utils"
	"github.com/glu-project/utils/authenticate"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type authService struct {
	db    models.DBTX
	ecoDB models.DBTX
	pb.UnimplementedAuthServiceServer

	authenticator authenticate.Authenticator
	userRepo      interface {
		Create(ctx context.Context, db models.DBTX, user *models.User) (string, error)
		GetByID(ctx context.Context, db models.DBTX, id string) (*models.User, error)
		GetByUsername(ctx context.Context, db models.DBTX, username string) (*models.User, error)
	}
	roleRepo interface {
		GetByID(ctx context.Context, db models.DBTX, id string) (*models.Role, error)
	}
}

func NewAuthService(db models.DBTX, ecosystemDB models.DBTX, auth authenticate.Authenticator) pb.AuthServiceServer {
	return &authService{
		db:            db,
		ecoDB:         ecosystemDB,
		userRepo:      new(postgres.UserRepository),
		authenticator: auth,
		roleRepo:      new(postgres.RoleRepository),
	}
}

func (s *authService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {

	user, err := s.userRepo.GetByUsername(ctx, s.db, req.Username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Errorf(codes.InvalidArgument, "username do not exist")
		}
		return nil, status.Errorf(codes.Internal, "s.userRepo.GetByUsername: unexpected error: %v", err)
	}

	if err := utils.CheckPassword(req.Password, user.Password.String); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "your password is incorrect %v", err)
	}

	role, err := s.roleRepo.GetByID(ctx, s.db, user.RoleID.String)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Errorf(codes.Internal, "s.roleRepo.GetByID: role_id is not exist")
		}
		return nil, status.Errorf(codes.Internal, "s.userRepo.GetByUsername: unexpected error: %v", err)
	}
	// change token
	tkn, err := s.authenticator.Generate(&authenticate.Payload{
		UserID: user.ID,
		RoleID: role.ID,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to generate token: %v", err)
	}

	// create login_user
	loginUser := &models.LoginUser{
		ID:            uuid.NewString(),
		UserID:        transformhelpers.StringToPgtypeText(user.ID),
		DevicesAccess: transformhelpers.StringToPgtypeText(req.DeviceAccess),
		Email:         user.Email,
		Token:         transformhelpers.StringToPgtypeText(tkn.Token),
	}

	if err := s.loginUserRepo.Create(ctx, s.db, loginUser); err != nil {
		return nil, status.Errorf(codes.Internal, "s.loginUserRepo.Create: unexpected error: %v", err)
	}

	return &pb.LoginResponse{
		UserId: user.ID,
		Token:  tkn.Token,
	}, nil
}

func (s *authService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	user := transform.PbToUserPtr(req.Data)
	user.ID = uuid.NewString()
	hashPassword, err := utils.HashPassword(user.Password.String)
	if err != nil {
		return nil, fmt.Errorf("unable to hash password: %v", err)
	}
	user.Password = pgtype.Text{String: hashPassword, Valid: true}

	id, err := s.userRepo.Create(ctx, s.db, user)
	if err != nil {
		return nil, err
	}
	tkn, err := s.authenticator.Generate(&authenticate.Payload{
		UserID: id,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to generate token: %v", err)
	}

	return &pb.RegisterResponse{
		Id:    id,
		Token: tkn.Token,
	}, nil
}
