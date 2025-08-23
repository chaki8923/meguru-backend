package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"meguru-backend/internal/domain/entity"
	"meguru-backend/internal/domain/repository"
	"meguru-backend/internal/infrastructure/email"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	userRepo               repository.UserRepository
	passwordResetTokenRepo repository.PasswordResetTokenRepository
	emailService           *email.EmailService
	jwtService             *JWTService
}

type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
}

type CreateUserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type LoginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginUserResponse struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Name  string    `json:"name"`
	Token string    `json:"token"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ForgotPasswordResponse struct {
	Message string `json:"message"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

type ResetPasswordResponse struct {
	Message string `json:"message"`
}

func NewUserUsecase(userRepo repository.UserRepository, passwordResetTokenRepo repository.PasswordResetTokenRepository, emailService *email.EmailService, jwtService *JWTService) *UserUsecase {
	return &UserUsecase{
		userRepo:               userRepo,
		passwordResetTokenRepo: passwordResetTokenRepo,
		emailService:           emailService,
		jwtService:             jwtService,
	}
}

func (u *UserUsecase) CreateUser(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error) {
	existingUser, _ := u.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Name:         req.Name,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return &CreateUserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (u *UserUsecase) LoginUser(ctx context.Context, req *LoginUserRequest) (*LoginUserResponse, error) {
	user, err := u.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}
	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	token, err := u.jwtService.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &LoginUserResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
		Token: token,
	}, nil
}

// ForgotPassword パスワードリセットトークンを生成してメール送信
func (u *UserUsecase) ForgotPassword(ctx context.Context, req *ForgotPasswordRequest) (*ForgotPasswordResponse, error) {
	// ユーザーを取得
	user, err := u.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		// セキュリティ上、ユーザーが存在しない場合でも成功レスポンスを返す
		return &ForgotPasswordResponse{Message: "パスワードリセット用のメールを送信しました"}, nil
	}
	if user == nil {
		// セキュリティ上、ユーザーが存在しない場合でも成功レスポンスを返す
		return &ForgotPasswordResponse{Message: "パスワードリセット用のメールを送信しました"}, nil
	}

	// 既存の未使用トークンを削除
	if err := u.passwordResetTokenRepo.DeleteByUserID(ctx, user.ID); err != nil {
		return nil, errors.New("パスワードリセット申請に失敗しました")
	}

	token, err := generateSecureToken(32)
	if err != nil {
		return nil, errors.New("パスワードリセット申請に失敗しました")
	}

	resetToken := &entity.PasswordResetToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		Used:      false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := u.passwordResetTokenRepo.Create(ctx, resetToken); err != nil {
		return nil, errors.New("パスワードリセット申請に失敗しました")
	}

	if u.emailService != nil {
		if err := u.sendPasswordResetEmail(ctx, user, token); err != nil {
			// メール送信失敗はログに記録するのみ
			println("パスワードリセットメール送信に失敗しました:", err.Error())
		}
	}

	return &ForgotPasswordResponse{Message: "パスワードリセット用のメールを送信しました"}, nil
}

// ResetPassword トークンを使用してパスワードをリセット
func (u *UserUsecase) ResetPassword(ctx context.Context, req *ResetPasswordRequest) (*ResetPasswordResponse, error) {
	resetToken, err := u.passwordResetTokenRepo.FindByToken(ctx, req.Token)
	if err != nil {
		return nil, errors.New("無効なトークンです")
	}
	if resetToken == nil {
		return nil, errors.New("無効なトークンです")
	}

	// トークンの有効期限を確認
	if time.Now().After(resetToken.ExpiresAt) {
		return nil, errors.New("トークンの有効期限が切れています")
	}

	// 新しいパスワードをハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("パスワードリセットに失敗しました")
	}

	// ユーザーのパスワードを更新
	user := &entity.User{
		ID:           resetToken.UserID,
		PasswordHash: string(hashedPassword),
		UpdatedAt:    time.Now(),
	}

	if err := u.userRepo.UpdatePassword(ctx, user.ID, user.PasswordHash); err != nil {
		return nil, errors.New("パスワードリセットに失敗しました")
	}

	// トークンを使用済みにマーク
	if err := u.passwordResetTokenRepo.UpdateUsed(ctx, resetToken.ID, true); err != nil {
		// パスワード更新は成功しているので、ログに記録するのみ
		println("トークンの使用済みマークに失敗しました:", err.Error())
	}

	return &ResetPasswordResponse{Message: "パスワードが正常に変更されました"}, nil
}

func generateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// sendPasswordResetEmail パスワードリセットメールを送信
func (u *UserUsecase) sendPasswordResetEmail(ctx context.Context, user *entity.User, token string) error {
	resetURL := "http://localhost:3000/auth/reset-password?token=" + token

	subject := "【meguru】パスワードリセットのご依頼"
	body := `
<html>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
	<div style="max-width: 600px; margin: 0 auto; padding: 20px;">
		<div style="background: linear-gradient(135deg, #F1B300, #563124); padding: 30px; border-radius: 10px; text-align: center; margin-bottom: 30px;">
			<h1 style="color: white; margin: 0; font-size: 24px;">パスワードリセット</h1>
			<p style="color: #fff3f0; margin: 10px 0 0 0;">パスワードのリセットを承りました</p>
		</div>
		
		<div style="background: #f8f9fa; padding: 25px; border-radius: 10px; margin-bottom: 25px;">
			<h2 style="color: #563124; margin-top: 0;">パスワードリセットについて</h2>
			<p>` + user.Name + ` 様</p>
			<p>パスワードリセットのご依頼を承りました。<br>
			以下のボタンをクリックして、新しいパスワードを設定してください。</p>
			<p style="color: #dc3545; font-size: 14px;">
				<strong>※このリンクは24時間有効です。</strong><br>
				<strong>※心当たりがない場合は、このメールを無視してください。</strong>
			</p>
		</div>
		
		<div style="text-align: center; margin: 30px 0;">
			<a href="` + resetURL + `" 
			   style="background: linear-gradient(135deg, #F1B300, #563124); 
			          color: white; 
			          padding: 15px 30px; 
			          text-decoration: none; 
			          border-radius: 25px; 
			          font-weight: bold; 
			          display: inline-block;">
				パスワードをリセットする
			</a>
		</div>
		
		<div style="border-top: 1px solid #dee2e6; padding-top: 20px; text-align: center; color: #6c757d; font-size: 12px;">
			<p>このメールは自動送信されています。ご返信いただいてもお答えできません。</p>
			<p>© 2024 meguru. All rights reserved.</p>
		</div>
	</div>
</body>
</html>`

	return u.emailService.SendEmail(user.Email, subject, body)
}
