package domain_service

import (
	"context"
	"meguru-backend/internal/domain/entity"
	"meguru-backend/internal/domain/repository"
)

type UserDomainService struct {
	userRepo repository.UserRepository
}

func NewUserDomainService(userRepo repository.UserRepository) *UserDomainService {
	return &UserDomainService{
		userRepo: userRepo,
	}
}

func (s *UserDomainService) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	return s.userRepo.GetByEmail(ctx, email)
}

func (s *UserDomainService) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	user, err := s.FindByEmail(ctx, email)
	if err != nil {
		return false, err
	}
	return user != nil, nil
}
