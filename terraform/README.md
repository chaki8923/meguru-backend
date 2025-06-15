# Meguru Backend - Terraform デプロイメント

このディレクトリには、Go（Gin）アプリケーション + Aurora Serverless PostgreSQL + App RunnerのインフラをデプロイするためのTerraformファイルが含まれています。

## 前提条件

1. **Terraform** がインストールされていること（バージョン 1.0 以上）
2. **AWS CLI** が設定されていること
3. 適切なAWS権限があること
4. **Docker イメージ** が既にECR Public またはプライベートECRにプッシュされていること

## デプロイ手順

### 1. 設定ファイルの準備

```bash
cd terraform
cp terraform.tfvars.example terraform.tfvars
```

`terraform.tfvars` ファイルを編集して、適切な値を設定してください：

```hcl
# 必須項目
image_identifier = "014372010501.dkr.ecr.ap-northeast-1.amazonaws.com/meguru/backend:latest"
database_password = "your-secure-password"

# オプション項目（必要に応じて変更）
project_name = "meguru-backend"
aws_region = "ap-northeast-1"
```

### 2. Terraformの初期化

```bash
terraform init
```

### 3. プランの確認

```bash
terraform plan
```

### 4. デプロイの実行

```bash
terraform apply
```

実行時に `yes` を入力してデプロイを確認してください。

### 5. 出力値の確認

デプロイが完了すると、以下の情報が出力されます：

```
app_runner_service_url = "https://xxxxxxxxxx.ap-northeast-1.awsapprunner.com"
aurora_cluster_endpoint = "meguru-backend-aurora-cluster.cluster-xxxxxxxxxx.ap-northeast-1.rds.amazonaws.com"
```

## 重要な改善点

CloudFormationから移行する際の主な改善点：

1. **ヘルスチェック設定の最適化**
   - Interval: 10秒（CloudFormationでは20秒でも失敗していた）
   - Timeout: 5秒（より短く設定）
   
2. **セキュリティグループの詳細制御**
   - より具体的なポート制御
   - セキュリティグループ間の参照

3. **リソースの依存関係管理**
   - `depends_on` での明示的な依存関係
   - より確実なデプロイ順序

## ヘルスチェックエンドポイントの実装

Goアプリケーションで `/health` エンドポイントを実装してください：

```go
func main() {
    r := gin.Default()
    
    // ヘルスチェックエンドポイント
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "status": "ok",
            "timestamp": time.Now().Unix(),
        })
    })
    
    // その他のルートの設定...
    
    r.Run(":" + os.Getenv("PORT"))
}
```

## リソースの削除

```bash
terraform destroy -auto-approve
```

実行時に `yes` を入力してリソースを削除してください。

## トラブルシューティング

### ヘルスチェックが失敗する場合

1. `/health` エンドポイントが実装されているか確認
2. アプリケーションが正しいポートで起動しているか確認
3. データベース接続が正常に動作しているか確認

### データベース接続エラー

1. セキュリティグループの設定を確認
2. データベースの認証情報が正しいか確認
3. VPC設定が適切か確認

## ファイル構成

```
terraform/
├── main.tf                    # メインのTerraform設定
├── terraform.tfvars.example  # 変数の例
├── terraform.tfvars          # 実際の変数値（gitignoreに追加）
└── README.md                  # このファイル
```

## セキュリティ注意事項

- `terraform.tfvars` ファイルにはパスワードなどの機密情報が含まれるため、バージョン管理に含めないでください
- プロダクション環境では、AWS Secrets Manager または Terraform Cloud の変数機能を使用することを推奨します 