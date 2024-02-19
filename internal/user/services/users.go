package services

import (
	"context"
	"fmt"
	"time"

	"github.com/glu/shopvui/idl/pb"
	"github.com/glu/shopvui/internal/userm/models"
	"github.com/glu/shopvui/internal/userm/repositories"
	"github.com/glu/shopvui/transform"
	"github.com/glu/shopvui/utils"
	"github.com/glu/shopvui/utils/authenticate"
	"github.com/google/uuid"
)

type userService struct {
	DB             models.DBTX
	UserRepository interface {
		GetByEmail(ctx context.Context, db models.DBTX, email string) (*models.User, error)
		CreateUserV2(ctx context.Context, db models.DBTX, user *models.User) error
		Create(ctx context.Context, db models.DBTX, user *models.User) error
		GetUserByID(ctx context.Context, db models.DBTX, userID string) (*models.User, error)
	}
	authenticator authenticate.Authenticator
	pb.UnimplementedUserServiceServer
}

func NewUserService(db models.DBTX, authenticator authenticate.Authenticator) pb.UserServiceServer {
	return &userService{
		DB:             db,
		UserRepository: new(repositories.UserRepository),
		authenticator:  authenticator,
	}
}

// gRPC service
func (u *userService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {

	user, err := u.UserRepository.GetByEmail(ctx, u.DB, req.Email)
	if err != nil {
		return nil, err
	}
	err = utils.CheckPassword(req.Password, user.Password)
	if err != nil {
		return nil, err
	}

	token, _, err := u.authenticator.CreateToken(user.UserID, 24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("unable to generate token: %v", err)
	}
	return &pb.LoginResponse{
		UserId:      user.UserID,
		AccessToken: token,
	}, nil
}

func validateRegisterRequest(req *pb.RegisterRequest) error {
	if req.User.Email == "" {
		return fmt.Errorf("userService: email is required")
	}
	if req.User.Password == "" {
		return fmt.Errorf("userService: password is required")
	}

	return nil
}

func (u *userService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	if err := validateRegisterRequest(req); err != nil {
		return nil, err
	}
	data := transform.PbToUserPtr(req.User)
	data.UserID = uuid.NewString()
	hashedPwd, err := utils.HashPassword(data.Password)
	if err != nil {
		return nil, fmt.Errorf("cant hash your password")
	}

	data.Password = hashedPwd
	if err := u.UserRepository.CreateUserV2(ctx, u.DB, data); err != nil {
		return nil, err
	}

	token, _, err := u.authenticator.CreateToken(data.UserID, 24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("unable to generate token: %v", err)
	}

	return &pb.RegisterResponse{
		UserId:      data.UserID,
		AccessToken: token,
	}, nil
}

func (u *userService) GetUserByID(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.GetUserByIDResponse, error) {

	user, err := u.UserRepository.GetUserByID(ctx, u.DB, req.UserId)
	if err != nil {
		return nil, fmt.Errorf("userService: %w", err)
	}

	return &pb.GetUserByIDResponse{
		User: transform.UserToPbPtr(user),
	}, nil
}

// // func validateRegisterRequest(req *pb.RegisterRequest) error {
// // 	if req.Email == "" {
// // 		return fmt.Errorf("email is required")
// // 	}
// // 	if req.Password == "" {
// // 		return fmt.Errorf("password is required")
// // 	}
// // 	return nil
// // }

// // func toGrpcUser(req *pb.RegisterRequest) (*entities.User, error) {
// // 	hashPassword, err := utils.HashPassword(req.Password)
// // 	if err != nil {
// // 		return &entities.User{}, err
// // 	}
// // 	return &entities.User{
// // 		ID:        database.Text(uuid.NewString()),
// // 		Email:     database.Text(req.Email),
// // 		FirstName: database.Text(req.FirstName),
// // 		LastName:  database.Text(req.LastName),
// // 		Password:  database.Text(hashPassword),
// // 	}, nil
// // }
// // func (u *UserService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {

// // 	if err := validateRegisterRequest(req); err != nil {
// // 		return nil, err
// // 	}
// // 	_, err := u.UserRepo.GetUser(ctx, u.DB, database.Text(req.GetEmail()))
// // 	if err == nil {
// // 		return nil, fmt.Errorf("user already exists")
// // 	}

// // 	ur, err := toGrpcUser(req)
// // 	if err != nil {
// // 		return nil, err
// // 	}
// // 	user, err := u.UserRepo.CreateUser(ctx, u.DB, ur)
// // 	if err != nil {
// // 		return nil, fmt.Errorf("can't create user: %v", err)
// // 	}

// // 	authenticateMaker, err := authenticate.NewPasetoMaker(u.config.authenticateSymmetricKey)
// // 	if err != nil {
// // 		return nil, fmt.Errorf("cannot create authenticate maker: %w", err)
// // 	}

// // 	accessauthenticate, _, err := authenticateMaker.Createauthenticate(
// // 		user.ID.String,
// // 		u.config.AccessauthenticateDuration,
// // 	)
// // 	if err != nil {
// // 		return nil, fmt.Errorf("cannot create access authenticate: %w", err)
// // 	}

// // 	return &pb.RegisterResponse{
// // 		User: &pb.User{
// // 			Id:         user.ID.String,
// // 			Email:      user.Email.String,
// // 			FirstName:  user.FirstName.String,
// // 			LastName:   user.LastName.String,
// // 			CreateDate: timestamppb.New(user.InsertedAt.Time),
// // 		},
// // 		Accessauthenticate: accessauthenticate,
// // 	}, nil
// // }

// func (u *UserService) GetMe(ctx context.Context, req *pb.GetMeRequest) (*pb.GetMeResponse, error) {
// 	pay, err := u.authorizeUser(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	user, err := u.UserRepo.GetUserByID(ctx, u.DB, database.Text(pay.UserID))
// 	if err != nil {
// 		return nil, fmt.Errorf("GetUser: %v", err)
// 	}

// 	return &pb.GetMeResponse{
// 		User: &pb.User{
// 			Id:         user.ID.String,
// 			Email:      user.Email.String,
// 			FirstName:  user.FirstName.String,
// 			LastName:   user.LastName.String,
// 			CreateDate: timestamppb.New(user.InsertedAt.Time),
// 		},
// 	}, nil
// }

// // RestFul api service
// func (s *Server) loginUser(ctx *gin.Context) {
// 	var req models.LoginUserRequest
// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, "cannot bin JSON")
// 		return
// 	}
// 	email := database.Text(req.Email)
// 	user, err := s.UserRepo.GetUser(ctx, s.db, email)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			ctx.JSON(http.StatusNotFound, "User not found")
// 			return
// 		}
// 		ctx.JSON(http.StatusInternalServerError, "User not found")
// 		return
// 	}
// 	if !user.Active.Bool {
// 		ctx.JSON(http.StatusInternalServerError, "user war detected")
// 	}

// 	err = utils.CheckPassword(req.Password, user.Password.String)
// 	if err != nil {
// 		ctx.JSON(http.StatusUnauthorized, err)
// 		return
// 	}

// 	accessauthenticate, _, err := s.authenticateMaker.Createauthenticate(
// 		user.ID.String,
// 		s.config.AccessauthenticateDuration,
// 	)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, "error")
// 		return
// 	}

// 	// refreshauthenticate, refreshPayload, err := server.authenticateMaker.Createauthenticate(
// 	// 	user.Username,
// 	// 	server.config.RefreshauthenticateDuration,
// 	// )
// 	// if err != nil {
// 	// 	ctx.JSON(http.StatusInternalServerError, "error")
// 	// 	return
// 	// }

// // 	rsp := models.LoginUserResponse{
// // 		Accessauthenticate: accessauthenticate,
// // 		//AccessauthenticateExpiresAt:  accessPayload.ExpiredAt,
// // 		//Refreshauthenticate:          refreshauthenticate,
// // 		//RefreshauthenticateExpiresAt: refreshPayload.ExpiredAt,
// // 		User: models.UserResponse{
// // 			Email:     user.Email.String,
// // 			FirstName: user.FirstName.String,
// // 			LastName:  user.LastName.String,
// // 			CreatedAt: user.InsertedAt.Time,
// // 			UpdatedAt: user.InsertedAt.Time,
// // 		},
// // 	}
// // 	ctx.JSON(http.StatusOK, rsp)
// // }

// // func toUser(req models.CreateUserRequest) (*entities.User, error) {
// // 	hashPassword, err := utils.HashPassword(req.Password)
// // 	if err != nil {
// // 		return &entities.User{}, err
// // 	}
// // 	return &entities.User{
// // 		ID:        database.Text(uuid.NewString()),
// // 		Email:     database.Text(req.Email),
// // 		FirstName: database.Text(req.FirstName),
// // 		LastName:  database.Text(req.LastName),
// // 		Password:  database.Text(hashPassword),
// // 	}, nil
// // }

// // func (s *Server) register(ctx *gin.Context) {
// // 	var req models.CreateUserRequest
// // 	if err := ctx.ShouldBindJSON(&req); err != nil {
// // 		ctx.JSON(http.StatusBadRequest, err)
// // 		return
// // 	}
// // 	email := database.Text(req.Email)
// // 	_, err := s.UserRepo.GetUser(ctx, s.db, email)
// // 	if err == nil {
// // 		ctx.JSON(http.StatusInternalServerError, "User already exists")
// // 		return
// // 	}
// // 	u, err := toUser(req)
// // 	if err != nil {
// // 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// // 		return
// // 	}

// // 	user, err := s.UserRepo.CreateUser(ctx, s.db, u)
// // 	if err != nil {
// // 		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("can't create user: %v", err))
// // 		return
// // 	}
// // 	rsp := models.UserResponse{
// // 		ID:        user.ID.String,
// // 		FirstName: user.FirstName.String,
// // 		LastName:  user.LastName.String,
// // 		Email:     user.Email.String,
// // 		CreatedAt: user.InsertedAt.Time,
// // 		UpdatedAt: user.UpdatedAt.Time,
// // 	}
// // 	ctx.JSON(http.StatusOK, rsp)
// // }

// // func toRole(req models.AddRoleRequest) *entities.Role {
// // 	return &entities.Role{
// // 		ID:   database.Text(uuid.NewString()),
// // 		Name: database.Text(req.Name),
// // 	}
// // }

// // func (s *Server) addRole(ctx *gin.Context) {
// // 	var req models.AddRoleRequest
// // 	if err := ctx.ShouldBindJSON(&req); err != nil {
// // 		ctx.JSON(http.StatusBadRequest, "can't bind JSON")
// // 		return
// // 	}
// // 	roles := toRole(req)
// // 	err := s.UserRepo.AddRoles(ctx, s.db, roles)
// // 	if err != nil {
// // 		fmt.Println(err.Error())
// // 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// // 		return
// // 	}
// // 	ctx.JSON(http.StatusOK, "success")
// // }

// // // need to update for scaling
// // func validateRole(role string) error {
// // 	if role == constants.RoleAdmin {
// // 		return fmt.Errorf("permission denied")
// // 	}
// // 	if role != constants.RoleUser && role != constants.RoleGuest {
// // 		return fmt.Errorf("invalid role")
// // 	}
// // 	return nil
// // }

// // func toUpdateRole(req models.UpdateRoleRequest, roleID string) *entities.UserRole {
// // 	return &entities.UserRole{
// // 		UserID: database.Text(req.UserID),
// // 		RoleID: database.Text(roleID),
// // 	}
// // }

// // func (s *Server) updateRole(ctx *gin.Context) {
// // 	var req models.UpdateRoleRequest
// // 	if err := ctx.ShouldBindJSON(&req); err != nil {
// // 		ctx.JSON(http.StatusBadRequest, "can't bind JSON")
// // 		return
// // 	}
// // 	if err := validateRole(req.RoleName); err != nil {
// // 		ctx.JSON(http.StatusBadRequest, err)
// // 		return
// // 	}
// // 	role, err := s.UserRepo.GetRole(ctx, s.db, database.Text(req.RoleName))
// // 	if err != nil {
// // 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// // 		return
// // 	}
// // 	user := toUpdateRole(req, role.ID.String)
// // 	_, err = s.UserRepo.UpdateRole(ctx, s.db, user)
// // 	if err != nil {
// // 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// // 		return
// // 	}
// // 	ctx.JSON(http.StatusOK, "success")
// // }
