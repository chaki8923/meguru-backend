package usecase

import (
	"context"
	"time"

	"meguru-backend/internal/domain/entity"
	"meguru-backend/internal/domain/repository"
	"meguru-backend/internal/infrastructure/email"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type StoreUsecase struct {
	storeRepo repository.StoreRepository
	emailService *email.EmailService
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

// 店舗登録用の新しいリクエスト構造体（フロントエンドから送信される最小限の情報）
type ShopRegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// 店舗登録用のレスポンス構造体
type ShopRegisterResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		Token string `json:"token"`
	} `json:"data"`
}

func NewStoreUsecase(storeRepo repository.StoreRepository, emailService *email.EmailService) *StoreUsecase {
	return &StoreUsecase{
		storeRepo: storeRepo,
		emailService: emailService,
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

// 店舗登録用の新しいメソッド
func (u *StoreUsecase) RegisterShop(ctx context.Context, req *ShopRegisterRequest) (*ShopRegisterResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 最小限の情報で店舗エンティティを作成
	store := &entity.Store{
		ID:		 uuid.New(),
		Email:		 req.Email,
		Password:	 string(hashedPassword),
		CreatedAt:	 time.Now(),
		UpdatedAt:	 time.Now(),
	}

	if err := u.storeRepo.Create(ctx, store); err != nil {
		return nil, err
	}

	// メールを送信（エラーが発生してもAPIのレスポンスには影響させない）
	if u.emailService != nil {
		subject := "店舗登録が完了しました"
		body := "店舗登録ありがとうございます。以下のリンクからサービスにアクセスしてください。\n <a href=\"https://meguru.com/login\">サービスにアクセス</a>"
		if err := u.emailService.SendEmail(req.Email, subject, body); err != nil {
			// メール送信に失敗してもログに記録するのみ、APIエラーは返さない
			// TODO: 本来はloggerを使用してログ出力を行う
			println("メール送信に失敗しました:", err.Error())
		}
	}

	// 簡易的なトークン生成（実際の実装では JWT などを使用）
	token := "temp_token_" + store.ID.String()

	response := &ShopRegisterResponse{
		Success: true,
		Message: "店舗登録が完了しました。確認メールを送信しました。",
	}
	response.Data.Token = token

	return response, nil
}