# Terraform Configuration for Meguru Backend

このディレクトリには、MeguruバックエンドをAWS App Runnerにデプロイするためのterraform設定が含まれています。

## 🚀 簡単デプロイ手順

### 1. 環境変数ファイルの自動生成
```bash
cd /path/to/meguru-backend/terraform
./generate-tfvars.sh
```

このスクリプトが以下を実行します：
- `../.env`ファイルから設定を読み取り
- `terraform.tfvars`ファイルを自動生成
- 必要な環境変数をすべて設定

### 2. Terraformデプロイ
```bash
terraform init
terraform plan
terraform apply
```

## 🔧 手動設定（従来の方法）

自動生成を使わない場合：

1. `terraform.tfvars.example`をコピー
```bash
cp terraform.tfvars.example terraform.tfvars
```

2. `terraform.tfvars`を編集して実際の値を設定

## 📋 主要リソース

- **App Runner Service**: アプリケーション実行環境
- **Aurora PostgreSQL Serverless v2**: データベース
- **VPC + Private Subnets**: ネットワーク（追加コストなし）
- **Security Groups**: セキュリティ設定
- **ECR Access IAM Role**: コンテナレジストリアクセス

## ⚠️ 重要事項

- **コスト最適化**: NAT Gatewayは使用していません（コスト削減）
- **セキュリティ**: terraform.tfvarsはGitignoreされています
- **環境変数**: .envファイルから自動的に変換されます

## 🛠️ トラブルシューティング

### Container exit code: 255
アプリケーションの起動に失敗している場合：
1. 環境変数が正しく設定されているか確認
2. App Runnerのログを確認
3. ECRイメージが正しくビルドされているか確認

### VPC Connector警告
「Public subnet ids detected」警告は解決済みです（Private subnet使用）。

## 📁 ファイル構成

- `main.tf`: メインの設定ファイル
- `terraform.tfvars.example`: 設定例
- `generate-tfvars.sh`: 自動生成スクリプト
- `.gitignore`: 機密ファイルの除外設定


## 削除
- terraform destroy -auto-approve