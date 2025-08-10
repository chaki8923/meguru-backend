package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"
	"strings"

	"meguru-backend/internal/domain/entity"
	"meguru-backend/internal/domain/repository"
	"meguru-backend/internal/infrastructure/email"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type StoreUsecase struct {
	storeRepo           repository.StoreRepository
	storeTokenRepo      repository.StoreEmailVerificationTokenRepository
	emailService        *email.EmailService
	enableEmailVerification bool // メール認証機能のON/OFF切り替え
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

func NewStoreUsecase(storeRepo repository.StoreRepository, storeTokenRepo repository.StoreEmailVerificationTokenRepository, emailService *email.EmailService, enableEmailVerification bool) *StoreUsecase {
	return &StoreUsecase{
		storeRepo:                storeRepo,
		storeTokenRepo:          storeTokenRepo,
		emailService:            emailService,
		enableEmailVerification: enableEmailVerification,
	}
}

func (u *StoreUsecase) CreateStore(ctx context.Context, req *CreateStoreRequest) (*CreateStoreResponse, error) {
	// パスワードをハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	store := &entity.Store{
		ID:           uuid.New(),
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		PhoneNumber:  req.PhoneNumber,
		Zipcode:      req.Zipcode,
		Prefecture:   req.Prefecture,
		City:         req.City,
		Street:       req.Street,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
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
		PasswordHash:	 string(hashedPassword),
		CreatedAt:	 time.Now(),
		UpdatedAt:	 time.Now(),
	}

	if err := u.storeRepo.Create(ctx, store); err != nil {
		return nil, err
	}

	// メール認証機能が有効な場合、認証トークンを生成・送信
	if u.enableEmailVerification && u.emailService != nil {
		if err := u.sendVerificationEmail(ctx, store, req.Email); err != nil {
			// メール送信失敗してもログに記録するのみ、APIエラーは返さない
			println("メール認証送信に失敗しました:", err.Error())
		}
	} else if u.emailService != nil {
		// 従来のメール送信（認証機能無効時）
		subject := "【meguru】店舗登録が完了しました"
		body := `
<html>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
	<div style="max-width: 600px; margin: 0 auto; padding: 20px;">
		<div style="background: linear-gradient(135deg, #ff7849, #ff6b35); padding: 30px; border-radius: 10px; text-align: center; margin-bottom: 30px;">
			<h1 style="color: white; margin: 0; font-size: 24px;">meguruへようこそ！</h1>
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
			<a href="http://localhost:3000/verify-email" 
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
			<p style="margin-top: 15px;">© 2025 meguru. All rights reserved.</p>
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

	// レスポンスメッセージを動的に設定
	var message string
	if u.enableEmailVerification {
		message = "店舗登録が完了しました。メールアドレスの認証が必要です。送信されたメール内のリンクをクリックして認証を完了してください。"
	} else {
		message = "店舗登録が完了しました。確認メールをお送りしましたので、メール内のリンクからサービスにアクセスしてください。"
	}

	response := &ShopRegisterResponse{
		Success: true,
		Message: message,
	}
	response.Data.Token = token

	return response, nil
}

// sendVerificationEmail メール認証用のメールを送信
func (u *StoreUsecase) sendVerificationEmail(ctx context.Context, store *entity.Store, email string) error {
	// 認証トークンを生成
	token := generateVerificationToken()
	expiresAt := time.Now().Add(24 * time.Hour) // 24時間有効

	// トークンをデータベースに保存
	if err := u.storeTokenRepo.Create(ctx, store.ID, token, expiresAt); err != nil {
		return fmt.Errorf("トークン作成エラー: %w", err)
	}

	// 認証URLを生成
	verifyURL := fmt.Sprintf("http://localhost:3000/verify-email?token=%s", token)

	// メール本文を作成
	subject := "【meguru】メールアドレスの認証が必要です"
	body := fmt.Sprintf(`
<html>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
	<div style="max-width: 600px; margin: 0 auto; padding: 20px;">
		<div style="background: linear-gradient(135deg, #ff7849, #ff6b35); padding: 30px; border-radius: 10px; text-align: center; margin-bottom: 30px;">
			<h1 style="color: white; margin: 0; font-size: 24px;">meguruへようこそ！</h1>
			<p style="color: #fff3f0; margin: 10px 0 0 0;">メールアドレスの認証をお願いします</p>
		</div>
		
		<div style="background: #f8f9fa; padding: 25px; border-radius: 10px; margin-bottom: 25px;">
			<h2 style="color: #ff6b35; margin-top: 0;">認証手順</h2>
			<ol style="padding-left: 20px;">
				<li style="margin-bottom: 10px;">下記の「メールアドレスを認証する」ボタンをクリック</li>
				<li style="margin-bottom: 10px;">認証完了後、ログインしてサービスをご利用ください</li>
			</ol>
		</div>
		
		<div style="text-align: center; margin: 30px 0;">
			<a href="%s" 
			   style="background: linear-gradient(135deg, #ff7849, #ff6b35); 
			          color: white; 
			          padding: 15px 30px; 
			          text-decoration: none; 
			          border-radius: 25px; 
			          font-weight: bold; 
			          display: inline-block;
			          box-shadow: 0 4px 15px rgba(255, 107, 53, 0.3);">
				メールアドレスを認証する
			</a>
		</div>
		
		<div style="background: #fff3cd; border: 1px solid #ffeaa7; padding: 15px; border-radius: 8px; margin-top: 30px;">
			<p style="margin: 0; font-size: 14px; color: #856404;">
				<strong>認証について：</strong><br>
				• この認証リンクは24時間有効です<br>
				• 認証完了後、サービスをご利用いただけます<br>
				• ログイン時のメールアドレス: %s
			</p>
		</div>
		
		<div style="margin-top: 30px; padding-top: 20px; border-top: 1px solid #eee; text-align: center; color: #666; font-size: 12px;">
			<p>このメールに心当たりがない場合は、このメールを無視してください。</p>
			<p style="margin-top: 15px;">© 2025 meguru. All rights reserved.</p>
		</div>
	</div>
</body>
</html>`, verifyURL, email)

	// メールを送信
	return u.emailService.SendEmail(email, subject, body)
}

// VerifyEmail メールアドレスの認証を処理
func (u *StoreUsecase) VerifyEmail(ctx context.Context, token string) error {
	// トークンを検索
	storeToken, err := u.storeTokenRepo.FindByToken(ctx, token)
	if err != nil {
		return fmt.Errorf("トークン検索エラー: %w", err)
	}
	if storeToken == nil {
		return errors.New("無効なトークンです")
	}

	// トークンの有効期限をチェック
	if storeToken.IsExpired() {
		// 期限切れトークンを削除
		u.storeTokenRepo.Delete(ctx, storeToken.ID)
		return errors.New("トークンの有効期限が切れています。再度登録をお試しください")
	}

	// 店舗を取得
	store, err := u.storeRepo.FindByID(ctx, storeToken.StoreID)
	if err != nil {
		return fmt.Errorf("店舗取得エラー: %w", err)
	}
	if store == nil {
		return errors.New("店舗が見つかりません")
	}

	// メール認証を完了
	now := time.Now()
	store.EmailVerifiedAt = &now
	if err := u.storeRepo.Update(ctx, store); err != nil {
		return fmt.Errorf("店舗更新エラー: %w", err)
	}

	// 使用済みトークンを削除
	if err := u.storeTokenRepo.Delete(ctx, storeToken.ID); err != nil {
		// ログのみ（認証は成功しているため）
		println("トークン削除エラー:", err.Error())
	}

	return nil
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
	if err := bcrypt.CompareHashAndPassword([]byte(store.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	// メール認証機能が有効な場合、認証状態をチェック
	if u.enableEmailVerification && !store.IsEmailVerified() {
		return nil, errors.New("メールアドレスの認証が完了していません。メールをご確認ください")
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