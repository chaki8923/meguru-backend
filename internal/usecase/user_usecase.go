package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"meguru-backend/internal/domain/domain_service"
	"meguru-backend/internal/domain/entity"
	"meguru-backend/internal/domain/repository"
	"meguru-backend/internal/middleware"
	dto "meguru-backend/internal/usecase/dto/users"
	"meguru-backend/internal/usecase/query_service"
)

type UserUsecase struct {
	userRepo          repository.UserRepository
	userDomainService *domain_service.UserDomainService
	userQueryService  query_service.UserQueryService
}

// NewUserUsecase creates a new UserUsecase with the provided UserRepository and DomainService.
func NewUserUsecase(
	userRepo repository.UserRepository,
	userDomainService *domain_service.UserDomainService,
	userQueryService query_service.UserQueryService,
) *UserUsecase {
	return &UserUsecase{
		userRepo:          userRepo,
		userDomainService: userDomainService,
		userQueryService:  userQueryService,
	}
}

func (u *UserUsecase) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.CreateUserResponse, string, error) {
	exists, err := u.userDomainService.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, "", err
	}
	if exists {
		return nil, "", errors.New("user with this email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	user, err := entity.NewUser(
		0,
		uuid.New().String(),
		req.Name,
		req.Email,
		string(hashedPassword),
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return nil, "", err
	}

	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, "", err
	}

	token, err := middleware.GenerateJWT(user.UserID.Value())
	if err != nil {
		return nil, "", err
	}

	resp := &dto.CreateUserResponse{
		ID:        user.ID,
		UserID:    user.UserID.String(),
		Name:      user.Name.String(),
		Email:     user.Email.String(),
		CreatedAt: user.CreatedAt,
	}

	return resp, token, nil
}

func (u *UserUsecase) Signin(ctx context.Context, req *dto.SigninRequest) (*dto.SigninResponse, error) {
	user, err := u.userDomainService.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	token, err := middleware.GenerateJWT(user.UserID.Value())
	if err != nil {
		return nil, err
	}

	userResp := &dto.GetUserResponse{
		ID:        user.ID,
		UserID:    user.UserID.String(),
		Email:     user.Email.String(),
		Name:      user.Name.String(),
		CreatedAt: user.CreatedAt,
	}

	return &dto.SigninResponse{
		Token: token,
		User:  userResp,
	}, nil
}

func (u *UserUsecase) GetUserByID(ctx context.Context, userID string) (*dto.GetUserResponse, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	user, err := u.userQueryService.GetUserByUserID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}

	resp := &dto.GetUserResponse{
		ID:        user.ID,
		UserID:    user.UserID.String(),
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
	}

	return resp, nil
}
