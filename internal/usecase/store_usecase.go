package usecase

import (
	"context"
	"time"

	"github.com/chaki-ry/meguru-backend/internal/domain/entity"
	"github.com/chaki-ry/meguru-backend/internal/domain/repository"
	"github.com/google/uuid"
)

type StoreUsecase struct {
	storeRepo repository.StoreRepository
}

type CreateStoreRequest struct {
	Name    string `json:"name" binding:"required"`
	Address string `json:"address" binding:"required"`
}

type CreateStoreResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
}

type UpdateStoreRequest struct {
	Name    string `json:"name" binding:"required"`
	Address string `json:"address" binding:"required"`
}

type UpdateStoreResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewStoreUsecase(storeRepo repository.StoreRepository) *StoreUsecase {
	return &StoreUsecase{
		storeRepo: storeRepo,
	}
}

func (u *StoreUsecase) CreateStore(ctx context.Context, req *CreateStoreRequest) (*CreateStoreResponse, error) {
	store := &entity.Store{
		ID:        uuid.New(),
		Name:      req.Name,
		Address:   req.Address,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := u.storeRepo.Create(ctx, store); err != nil {
		return nil, err
	}

	return &CreateStoreResponse{
		ID:        store.ID,
		Name:      store.Name,
		Address:   store.Address,
		CreatedAt: store.CreatedAt,
	}, nil
}

func (u *StoreUsecase) UpdateStore(ctx context.Context, id uuid.UUID, req *UpdateStoreRequest) (*UpdateStoreResponse, error) {
	store, err := u.storeRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	store.Name = req.Name
	store.Address = req.Address
	store.UpdatedAt = time.Now()

	if err := u.storeRepo.Update(ctx, store); err != nil {
		return nil, err
	}

	return &UpdateStoreResponse{
		ID:        store.ID,
		Name:      store.Name,
		Address:   store.Address,
		UpdatedAt: store.UpdatedAt,
	}, nil
}
