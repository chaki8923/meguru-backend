package domain_service

import (
	"context"

	"meguru-backend/internal/domain/entity"
	"meguru-backend/internal/domain/repository"
)

type StoreDomainService struct {
	storeRepo repository.StoreRepository
}

func NewStoreService(storeRepo repository.StoreRepository) *StoreDomainService {
	return &StoreDomainService{
		storeRepo: storeRepo,
	}
}

func (s *StoreDomainService) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	store, err := s.storeRepo.FindByEmail(ctx, email)
	if err != nil {
		return false, err
	}
	return store != nil, nil
}

func (s *StoreDomainService) FindByEmail(ctx context.Context, email string) (*entity.Store, error) {
	return s.storeRepo.FindByEmail(ctx, email)
}
