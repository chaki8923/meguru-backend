package usecase

import (
	"context"
	"errors"
	"time"
	"strings"

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

// 店舗ログイン用のリクエスト構造体
type StoreSignInRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// 店舗ログイン用のレスポンス構造体
type StoreSignInResponse struct {
	Token string `json:"token"`
}

// 店舗プロフィール取得用のレスポンス構造体
type StoreProfileResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Zipcode     string `json:"zipcode"`
	Prefecture  string `json:"prefecture"`
	City        string `json:"city"`
	Street      string `json:"street"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// 店舗情報更新用のリクエスト構造体
type StoreUpdateRequest struct {
	Name        string `json:"name" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	PhoneNumber string `json:"phone_number"`
	Zipcode     string `json:"zipcode"`
	Prefecture  string `json:"prefecture"`
	City        string `json:"city"`
	Street      string `json:"street"`
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
		subject := "【めぐる】店舗登録が完了しました"
		body := `
<html>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
	<div style="max-width: 600px; margin: 0 auto; padding: 20px;">
		<div style="background: linear-gradient(135deg, #ff7849, #ff6b35); padding: 30px; border-radius: 10px; text-align: center; margin-bottom: 30px;">
			<h1 style="color: white; margin: 0; font-size: 24px;">めぐるへようこそ！</h1>
			<p style="color: #fff3f0; margin: 10px 0 0 0;">店舗登録が完了しました</p>
		</div>
		
		<div style="background: #f8f9fa; padding: 25px; border-radius: 10px; margin-bottom: 25px;">
			<h2 style="color: #ff6b35; margin-top: 0;">次のステップ</h2>
			<ol style="padding-left: 20px;">
				<li style="margin-bottom: 10px;">下記のリンクをクリックしてサービスにアクセス</li>
				<li style="margin-bottom: 10px;">ログインして店舗情報を詳細設定</li>
				<li style="margin-bottom: 10px;">チラシ作成や商品登録を開始</li>
			</ol>
		</div>
		
		<div style="text-align: center; margin: 30px 0;">
			<a href="http://localhost:3000/login" 
			   style="background: linear-gradient(135deg, #ff7849, #ff6b35); 
			          color: white; 
			          padding: 15px 30px; 
			          text-decoration: none; 
			          border-radius: 25px; 
			          font-weight: bold; 
			          display: inline-block;
			          box-shadow: 0 4px 15px rgba(255, 107, 53, 0.3);">
				サービスにアクセス
			</a>
		</div>
		
		<div style="background: #fff3cd; border: 1px solid #ffeaa7; padding: 15px; border-radius: 8px; margin-top: 30px;">
			<p style="margin: 0; font-size: 14px; color: #856404;">
				<strong>ログイン情報：</strong><br>
				メールアドレス: ` + req.Email + `<br>
				パスワード: （新規登録時に設定したパスワード）
			</p>
		</div>
		
		<div style="margin-top: 30px; padding-top: 20px; border-top: 1px solid #eee; text-align: center; color: #666; font-size: 12px;">
			<p>このメールに心当たりがない場合は、このメールを無視してください。</p>
			<p style="margin-top: 15px;">© 2025 めぐる. All rights reserved.</p>
		</div>
	</div>
</body>
</html>`
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
		Message: "店舗登録が完了しました。確認メールをお送りしましたので、メール内のリンクからサービスにアクセスしてください。",
	}
	response.Data.Token = token

	return response, nil
}

// 店舗ログイン用のメソッド
func (u *StoreUsecase) SignIn(ctx context.Context, req *StoreSignInRequest) (*StoreSignInResponse, error) {
	// メールアドレスで店舗を検索
	store, err := u.storeRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if store == nil {
		return nil, errors.New("invalid email or password")
	}

	// パスワードを検証
	if err := bcrypt.CompareHashAndPassword([]byte(store.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	// 簡易的なトークン生成（実際の実装では JWT などを使用）
	token := "auth_token_" + store.ID.String()

	return &StoreSignInResponse{
		Token: token,
	}, nil
}

// トークンからストアIDを取得するヘルパー関数
func (u *StoreUsecase) getStoreIDFromToken(token string) (uuid.UUID, error) {
	// 簡易実装: "auth_token_" または "temp_token_" プレフィックスを削除してUUIDを取得
	var uuidStr string
	if strings.HasPrefix(token, "auth_token_") {
		uuidStr = strings.TrimPrefix(token, "auth_token_")
	} else if strings.HasPrefix(token, "temp_token_") {
		uuidStr = strings.TrimPrefix(token, "temp_token_")
	} else {
		return uuid.Nil, errors.New("invalid token format")
	}

	storeID, err := uuid.Parse(uuidStr)
	if err != nil {
		return uuid.Nil, errors.New("invalid token")
	}

	return storeID, nil
}

// 店舗プロフィール取得
func (u *StoreUsecase) GetProfile(ctx context.Context, token string) (*StoreProfileResponse, error) {
	storeID, err := u.getStoreIDFromToken(token)
	if err != nil {
		return nil, err
	}

	store, err := u.storeRepo.FindByID(ctx, storeID)
	if err != nil {
		return nil, err
	}
	if store == nil {
		return nil, errors.New("store not found")
	}

	return &StoreProfileResponse{
		ID:          store.ID.String(),
		Name:        store.Name,
		Email:       store.Email,
		PhoneNumber: store.PhoneNumber,
		Zipcode:     store.Zipcode,
		Prefecture:  store.Prefecture,
		City:        store.City,
		Street:      store.Street,
		CreatedAt:   store.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   store.UpdatedAt.Format(time.RFC3339),
	}, nil
}

// 店舗プロフィール更新
func (u *StoreUsecase) UpdateProfile(ctx context.Context, token string, req *StoreUpdateRequest) (*StoreProfileResponse, error) {
	storeID, err := u.getStoreIDFromToken(token)
	if err != nil {
		return nil, err
	}

	store, err := u.storeRepo.FindByID(ctx, storeID)
	if err != nil {
		return nil, err
	}
	if store == nil {
		return nil, errors.New("store not found")
	}

	// 店舗情報を更新
	store.Name = req.Name
	store.Email = req.Email
	store.PhoneNumber = req.PhoneNumber
	store.Zipcode = req.Zipcode
	store.Prefecture = req.Prefecture
	store.City = req.City
	store.Street = req.Street
	store.UpdatedAt = time.Now()

	if err := u.storeRepo.Update(ctx, store); err != nil {
		return nil, err
	}

	return &StoreProfileResponse{
		ID:          store.ID.String(),
		Name:        store.Name,
		Email:       store.Email,
		PhoneNumber: store.PhoneNumber,
		Zipcode:     store.Zipcode,
		Prefecture:  store.Prefecture,
		City:        store.City,
		Street:      store.Street,
		CreatedAt:   store.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   store.UpdatedAt.Format(time.RFC3339),
	}, nil
}