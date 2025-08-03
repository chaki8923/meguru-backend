package usecase

import (
	"crypto/rand"
	"encoding/hex"
)

// generateVerificationToken メール認証用のトークンを生成
func generateVerificationToken() string {
	bytes := make([]byte, 32) // 256ビット
	if _, err := rand.Read(bytes); err != nil {
		// フォールバック: より単純な方法
		return generateSimpleToken()
	}
	return hex.EncodeToString(bytes)
}

// generateSimpleToken シンプルなトークン生成（フォールバック用）
func generateSimpleToken() string {
	bytes := make([]byte, 16) // 128ビット
	rand.Read(bytes) // エラーを無視（フォールバック）
	return hex.EncodeToString(bytes)
} 