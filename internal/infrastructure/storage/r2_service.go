package storage

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
)

type R2Service struct {
	client    *s3.S3
	bucketURL string
	publicURL string
}

func NewR2Service() *R2Service {
	endpoint := os.Getenv("R2_ENDPOINT")
	accessKey := os.Getenv("R2_ACCESS_KEY")
	secretKey := os.Getenv("R2_SECRET_KEY")
	bucketURL := os.Getenv("R2_BUCKET_URL")
	publicURL := os.Getenv("R2_PUBLIC_BUCKET_DOMAIN")

	// R2_BUCKET_URLが設定されていない場合、R2_ENDPOINTを使用
	if bucketURL == "" && endpoint != "" {
		bucketURL = endpoint
	}

	log.Printf("R2 Config - Endpoint: %s", endpoint)
	log.Printf("R2 Config - BucketURL: %s", bucketURL)
	log.Printf("R2 Config - PublicURL: %s", publicURL)
	if len(accessKey) > 10 {
		log.Printf("R2 Config - AccessKey: %s...", accessKey[:10])
	} else {
		log.Printf("R2 Config - AccessKey: %s", accessKey)
	}

	if endpoint == "" || accessKey == "" || secretKey == "" {
		log.Fatal("R2 環境変数が設定されていません（R2_ENDPOINT, R2_ACCESS_KEY, R2_SECRET_KEY）")
	}

	if publicURL == "" {
		log.Println("警告: R2_PUBLIC_BUCKET_DOMAIN が設定されていません。内部URLを使用します。")
		publicURL = bucketURL
	}

	// AWS SDK を使用してR2に接続
	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String("auto"),
		Endpoint:         aws.String(endpoint),
		Credentials:      credentials.NewStaticCredentials(accessKey, secretKey, ""),
		S3ForcePathStyle: aws.Bool(true),
	})

	if err != nil {
		log.Fatalf("R2セッションの作成に失敗しました: %v", err)
	}

	return &R2Service{
		client:    s3.New(sess),
		bucketURL: strings.TrimSuffix(bucketURL, "/"),
		publicURL: strings.TrimSuffix(publicURL, "/"),
	}
}

// 商品画像をアップロード
func (r *R2Service) UploadProductImage(fileHeader *multipart.FileHeader, storeID, productID string) (string, error) {
	// ファイルを開く
	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("ファイルを開けませんでした: %w", err)
	}
	defer file.Close()

	// ファイル内容を読み取り
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("ファイルの読み取りに失敗しました: %w", err)
	}

	// ファイル拡張子を取得
	ext := filepath.Ext(fileHeader.Filename)
	if ext == "" {
		ext = ".jpg" // デフォルト拡張子
	}

	// 一意のファイル名を生成
	fileName := fmt.Sprintf("%s_%s_%d%s", productID, uuid.New().String()[:8], time.Now().Unix(), ext)
	
	// R2のキー（パス）を生成: stores/{storeID}/products/{fileName}
	key := fmt.Sprintf("stores/%s/products/%s", storeID, fileName)

	// Content-Typeを推定
	contentType := r.getContentType(ext)

	// R2にアップロード
	bucketName := r.getBucketName()
	log.Printf("Uploading to R2 - Bucket: %s, Key: %s, ContentType: %s", bucketName, key, contentType)
	
	_, err = r.client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(key),
		Body:        bytes.NewReader(fileBytes),
		ContentType: aws.String(contentType),
		ACL:         aws.String("public-read"), // 公開読み取り可能
	})

	if err != nil {
		return "", fmt.Errorf("R2へのアップロードに失敗しました: %w", err)
	}

	// 公開URLを生成（公開ドメインを使用）
	imageURL := fmt.Sprintf("%s/%s", r.publicURL, key)
	
	log.Printf("画像をR2にアップロードしました: %s", imageURL)
	return imageURL, nil
}

// 画像を削除
func (r *R2Service) DeleteProductImage(imageURL string) error {
	// URLからキーを抽出（公開URLから）
	key := strings.TrimPrefix(imageURL, r.publicURL+"/")
	
	if key == imageURL {
		// 公開URLでマッチしない場合、内部URLも試す（後方互換性）
		key = strings.TrimPrefix(imageURL, r.bucketURL+"/")
		if key == imageURL {
			return fmt.Errorf("無効な画像URL: %s", imageURL)
		}
	}

	// R2から削除
	_, err := r.client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(r.getBucketName()),
		Key:    aws.String(key),
	})

	if err != nil {
		return fmt.Errorf("R2からの削除に失敗しました: %w", err)
	}

	log.Printf("画像をR2から削除しました: %s", imageURL)
	return nil
}

// Content-Typeを推定
func (r *R2Service) getContentType(ext string) string {
	switch strings.ToLower(ext) {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	default:
		return "image/jpeg"
	}
}

// バケット名を環境変数から取得
func (r *R2Service) getBucketName() string {
	// CloudFlare R2のURLからバケット名を抽出
	// 例: https://account-id.r2.cloudflarestorage.com/bucket-name → bucket-name
	
	// まず専用の環境変数があるかチェック
	if bucketName := os.Getenv("R2_BUCKET_NAME"); bucketName != "" {
		log.Printf("Using R2_BUCKET_NAME: %s", bucketName)
		return bucketName
	}
	
	// URLからパスの最後の部分を抽出
	if strings.Contains(r.bucketURL, "/") {
		parts := strings.Split(r.bucketURL, "/")
		if len(parts) > 0 {
			bucketName := parts[len(parts)-1]
			if bucketName != "" {
				log.Printf("Extracted bucket name from URL path: %s", bucketName)
				return bucketName
			}
		}
	}
	
	log.Printf("Using fallback bucket name: meguru")
	return "meguru" // フォールバック
} 