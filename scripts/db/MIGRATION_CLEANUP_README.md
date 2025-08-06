# マイグレーションファイル整理について

## 概要
2025年1月に実施したマイグレーションファイルの整理により、20個の分散したマイグレーションファイルを2個の統合ファイルに集約しました。

## 変更前の問題点

### 重複テーブル作成
- **users table**: `000001_create_initial_tables.up.sql` と `20250606120000_create_users_table.up.sql` で重複作成
- **stores table**: `000001_create_initial_tables.up.sql` と `20250607120000_create_stores_table.up.sql` で重複作成  
- **push_subscriptions table**: `000001_create_initial_tables.up.sql` と `20250629120000_create_push_subscriptions_table.up.sql` で重複作成

### 断片的なスキーマ変更
- stores tableに対して5回以上の個別変更
  - address削除 → prefecture/city/street追加
  - email/password/phone_number/zipcode追加
  - 上記カラムのNOT NULL制約削除（2回実行）
  - email_verified_at追加

## 変更後の構造

### 統合されたファイル
- `000001_create_all_tables.up.sql`: すべてのテーブル作成とインデックス
- `000001_create_all_tables.down.sql`: 全テーブル削除のロールバック

### 最終的なテーブル構造

#### users
```sql
id BIGSERIAL PRIMARY KEY,
user_id VARCHAR(255) NOT NULL UNIQUE,
name VARCHAR(255) NOT NULL,
email VARCHAR(255) UNIQUE NOT NULL,
password_hash VARCHAR(255) NOT NULL,
created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
```

#### stores  
```sql
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
name VARCHAR(255),
email VARCHAR(255),
password VARCHAR(255),
phone_number VARCHAR(255),
zipcode VARCHAR(255),
prefecture VARCHAR(255),
city VARCHAR(255),
street VARCHAR(255),
email_verified_at TIMESTAMP WITH TIME ZONE,
created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
```

#### 新規テーブル
- `store_email_verification_tokens`: メール認証トークン管理
- `store_products`: 店舗固有の商品情報

## バックアップ
元のマイグレーションファイルは `migrations_backup/` ディレクトリに保存されています。

## 使用方法

### 新規環境でのマイグレーション実行
```bash
# 統合されたマイグレーションを実行
migrate -path scripts/db/migrations -database "postgres://..." up
```

### ロールバック（必要時）
```bash
# 全テーブル削除
migrate -path scripts/db/migrations -database "postgres://..." down
```

## 注意事項
- 既存の本番環境では、このマイグレーションを実行する前に現在のスキーマとの互換性を確認してください
- バックアップファイルは必要に応じて参照できます
- 統合マイグレーションは最終的なDB構造と同等の状態を作成します 