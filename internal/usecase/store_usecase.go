package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"meguru-backend/internal/domain/domain_service"
	"meguru-backend/internal/domain/entity"
	repository_interface "meguru-backend/internal/domain/repository"
	"meguru-backend/internal/middleware"
	dto "meguru-backend/internal/usecase/dto/stores"
	"meguru-backend/internal/usecase/query_service"
)

type StoreUsecase struct {
	storeRepo          repository_interface.StoreRepository
	storeDomainService *domain_service.StoreDomainService
	storeQueryService  query_service.StoreQueryService
}

func NewStoreUsecase(
	storeRepo repository_interface.StoreRepository,
	storeDomainService *domain_service.StoreDomainService,
	storeQueryService query_service.StoreQueryService,
) *StoreUsecase {
	return &StoreUsecase{
		storeRepo:          storeRepo,
		storeDomainService: storeDomainService,
		storeQueryService:  storeQueryService,
	}
}

func (u *StoreUsecase) CreateStore(ctx context.Context, req *dto.CreateStoreRequest) (*dto.CreateStoreResponse, error) {
	exists, err := u.storeDomainService.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("store with this email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	store, err := entity.NewStore(
		0, // IDはDB側で生成されるため0でOK
		uuid.New().String(),
		req.Name,
		req.Email,
		string(hashedPassword),
		req.PhoneNumber,
		req.Zipcode,
		req.Prefecture,
		req.City,
		req.Street,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return nil, err
	}

	if err := u.storeRepo.Create(ctx, store); err != nil {
		return nil, err
	}

	token, err := middleware.GenerateJWT(store.StoreID.Value())
	if err != nil {
		return nil, err
	}

	resp := &dto.GetStoreResponse{
		ID:          store.ID,
		StoreID:     store.StoreID.Value().String(),
		Name:        store.Name,
		Email:       store.Email.Value(),
		PhoneNumber: store.PhoneNumber.Value(),
		Zipcode:     store.Zipcode.Value(),
		Prefecture:  store.Prefecture.Value(),
		City:        store.City.Value(),
		Street:      store.Street.Value(),
		CreatedAt:   store.CreatedAt,
	}

	return &dto.CreateStoreResponse{
		Token: token,
		Store: resp,
	}, nil

}

func (u *StoreUsecase) SigninStore(ctx context.Context, req *dto.SigninStoreRequest) (*dto.SigninStoreResponse, error) {
	store, err := u.storeDomainService.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if store == nil {
		return nil, errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(store.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	token, err := middleware.GenerateJWT(store.StoreID.Value())
	if err != nil {
		return nil, err
	}

	storeResp := &dto.GetStoreResponse{
		ID:          store.ID,
		StoreID:     store.StoreID.Value().String(),
		Name:        store.Name,
		Email:       store.Email.Value(),
		PhoneNumber: store.PhoneNumber.Value(),
		Zipcode:     store.Zipcode.Value(),
		Prefecture:  store.Prefecture.Value(),
		City:        store.City.Value(),
		Street:      store.Street.Value(),
		CreatedAt:   store.CreatedAt,
	}

	return &dto.SigninStoreResponse{
		Token: token,
		Store: storeResp,
	}, nil
}

func (u *StoreUsecase) GetStoreByID(ctx context.Context, storeID string) (*dto.GetStoreResponse, error) {
	uuidStoreID, err := uuid.Parse(storeID)
	if err != nil {
		return nil, err
	}

	store, err := u.storeQueryService.GetStoreByID(ctx, uuidStoreID)
	if err != nil {
		return nil, err
	}
	if store == nil {
		return nil, nil
	}

	resp := &dto.GetStoreResponse{
		ID:          store.ID,
		StoreID:     store.StoreID.String(),
		Name:        store.Name,
		Email:       store.Email,
		PhoneNumber: store.PhoneNumber,
		Zipcode:     store.Zipcode,
		Prefecture:  store.Prefecture,
		City:        store.City,
		Street:      store.Street,
		CreatedAt:   store.CreatedAt,
	}

	return resp, nil
}
