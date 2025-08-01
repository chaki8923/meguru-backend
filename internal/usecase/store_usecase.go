package usecase

import (
	"context"
	"time"

	"meguru-backend/internal/domain/entity"
	"meguru-backend/internal/domain/repository"
	"github.com/google/uuid"
)

type StoreUsecase struct {
	storeRepo repository.StoreRepository
}

type CreateStoreRequest struct {
	Name        string `json:"name" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Zipcode     string `json:"zipcode" binding:"required"`
	Prefecture  string `json:"prefecture" binding:"required"`
	City        string `json:"city" binding:"required"`
	Street      string `json:"street" binding:"required"`
}

type CreateStoreResponse struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	PhoneNumber string   `json:"phone_number"`
	Zipcode    string    `json:"zipcode"`
	Prefecture string    `json:"prefecture"`
	City       string    `json:"city"`
	Street     string    `json:"street"`
	CreatedAt  time.Time `json:"created_at"`
}

type UpdateStoreRequest struct {
	Name       string `json:"name" binding:"required"`
	Prefecture string `json:"prefecture" binding:"required"`
	City       string `json:"city" binding:"required"`
	Street     string `json:"street" binding:"required"`
}

type UpdateStoreResponse struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Prefecture string    `json:"prefecture"`
	City       string    `json:"city"`
	Street     string    `json:"street"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func NewStoreUsecase(storeRepo repository.StoreRepository) *StoreUsecase {
	return &StoreUsecase{
		storeRepo: storeRepo,
	}
}

func (u *StoreUsecase) CreateStore(ctx context.Context, req *CreateStoreRequest) (*CreateStoreResponse, error) {
	store := &entity.Store{
		ID:          uuid.New(),
		Name:        req.Name,
		Email:       req.Email,
		Password:    req.Password, // Note: In a real application, hash the password
		PhoneNumber: req.PhoneNumber,
		Zipcode:     req.Zipcode,
		Prefecture:  req.Prefecture,
		City:        req.City,
		Street:      req.Street,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := u.storeRepo.Create(ctx, store); err != nil {
		return nil, err
	}

	return &CreateStoreResponse{
		ID:          store.ID,
		Name:        store.Name,
		Email:       store.Email,
		PhoneNumber: store.PhoneNumber,
		Zipcode:     store.Zipcode,
		Prefecture:  store.Prefecture,
		City:        store.City,
		Street:      store.Street,
		CreatedAt:   store.CreatedAt,
	}, nil
}

func (u *StoreUsecase) UpdateStore(ctx context.Context, id uuid.UUID, req *UpdateStoreRequest) (*UpdateStoreResponse, error) {
	store, err := u.storeRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	store.Name = req.Name
	store.Prefecture = req.Prefecture
	store.City = req.City
	store.Street = req.Street
	store.UpdatedAt = time.Now()

	if err := u.storeRepo.Update(ctx, store); err != nil {
		return nil, err
	}

	return &UpdateStoreResponse{
		ID:         store.ID,
		Name:       store.Name,
		Prefecture: store.Prefecture,
		City:       store.City,
		Street:     store.Street,
		UpdatedAt:  store.UpdatedAt,
	}, nil
}

func (u *StoreUsecase) GetStore(ctx context.Context, id uuid.UUID) (*entity.Store, error) {
	return u.storeRepo.FindByID(ctx, id)
}

func (u *StoreUsecase) GetAllStores(ctx context.Context) ([]*entity.Store, error) {
	return u.storeRepo.FindAll(ctx)
}