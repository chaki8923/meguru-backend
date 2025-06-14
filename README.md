# Meguru プロジェクト

Next.js フロントエンドと Go (Gin) バックエンドを使用したユーザー管理システムです。

## アーキテクチャ

### フロントエンド (meguru-front)
- **フレームワーク**: Next.js 15 with TypeScript
- **スタイリング**: Tailwind CSS
- **ポート**: 3000

### バックエンド (meguru-backend)
- **言語**: Go 1.23
- **フレームワーク**: Gin
- **アーキテクチャ**: オニオンアーキテクチャ
- **データベース**: PostgreSQL
- **ポート**: 8080

## 開発環境セットアップ

### 前提条件
- Docker & Docker Compose
- Git
- curl (API テスト用)

**Optional (ローカル開発用):**
- Node.js (18以上)
- Go (1.23以上)

### バックエンド開発環境のセットアップ

#### 1. リポジトリのクローン
```bash
# プロジェクトをクローン
git clone <repository-url>
cd mguru
```

#### 2. バックエンドディレクトリに移動
```bash
cd meguru-backend
```

#### 3. 環境変数ファイルの作成
```bash
# .envファイルを作成
cat > .env << EOF
DB_HOST=localhost
DB_PORT=5432
DB_USER=meguru_user
DB_PASSWORD=meguru_password
DB_NAME=meguru_db
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
PORT=8080
DATABASE_URL=postgres://meguru_user:meguru_password@localhost:5432/meguru_db?sslmode=disable
EOF
```

#### 4. Docker Composeでサービスを起動
```bash
# データベースとアプリケーションを起動
docker-compose up -d

# ログを確認（オプション）
docker-compose logs -f
```

**注意**: Dockerを使用する場合、Go の依存関係は Dockerfile 内で自動的にダウンロード・ビルドされるため、事前に `go mod download` を実行する必要はありません。

これにより以下が起動されます:
- PostgreSQL データベース (ポート 5432)
- Go バックエンドアプリケーション (ポート 8080)

#### 5. サービスの状態確認
```bash
# コンテナの状態確認
docker-compose ps

# アプリケーションが起動していることを確認
curl http://localhost:8080/api/v1/users/signup \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"test123","name":"Test User"}'
```

#### 6. API テスト

**ユーザー登録のテスト:**
```bash
# 新規ユーザー登録
curl -X POST http://localhost:8080/api/v1/users/signup \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "name": "テストユーザー"
  }'

# 期待される正常レスポンス:
# {
#   "data": {
#     "id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
#     "email": "user@example.com",
#     "name": "テストユーザー",
#     "created_at": "2025-XX-XXTXX:XX:XX.XXXXXXXZ"
#   }
# }
```

**重複登録のテスト:**
```bash
# 同じメールアドレスで再度登録（エラーになることを確認）
curl -X POST http://localhost:8080/api/v1/users/signup \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "different123",
    "name": "別のユーザー"
  }'

# 期待されるエラーレスポンス:
# {"error":"user with this email already exists"}
```

#### 7. トラブルシューティング

**コンテナが起動しない場合:**
```bash
# ログを確認
docker-compose logs

# コンテナを停止して再起動
docker-compose down
docker-compose up -d --build
```

**ポートが使用中の場合:**
```bash
# 使用中のポートを確認
lsof -i :8080
lsof -i :5432

# 必要に応じて該当プロセスを停止してから再試行
```

**データベースをリセットしたい場合:**
```bash
# 全てのコンテナとボリュームを削除
docker-compose down -v

# 再起動
docker-compose up -d
```

#### 8. ローカル開発（Docker を使わない場合）

Docker を使わずにローカルでGoアプリケーションを開発したい場合：

**前提条件**: 
- Go 1.23以上がインストールされていること
- PostgreSQL が別途起動していること

```bash
# Go の依存関係をダウンロード
go mod download

# データベースのみ Docker で起動
docker-compose up -d db

# アプリケーションをローカルで実行
go run cmd/server/main.go
```

**ローカル開発用の環境変数 (.env)**:
```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=meguru_user
DB_PASSWORD=meguru_password
DB_NAME=meguru_db
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
PORT=8080
DATABASE_URL=postgres://meguru_user:meguru_password@localhost:5432/meguru_db?sslmode=disable
```

### フロントエンドのセットアップ

1. フロントエンドディレクトリに移動:
```bash
cd ../meguru-front
```

2. 依存関係のインストールと開発サーバーの起動:
```bash
npm install
npm run dev
```

フロントエンドは http://localhost:3000 で利用可能になります。

## API エンドポイント

### ユーザー登録
- **URL**: `POST /api/v1/users/signup`
- **Content-Type**: `application/json`
- **リクエストボディ**:
```json
{
  "email": "user@example.com",
  "password": "password123",
  "name": "ユーザー名"
}
```
- **レスポンス例**:
```json
{
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "name": "ユーザー名",
    "created_at": "2025-05-30T09:22:47.546Z"
  }
}
```
- **エラーレスポンス例**:
```json
{
  "error": "user with this email already exists"
}
```

## ディレクトリ構造

```
meguru/
├── README.md              # このファイル
├── meguru-front/          # Next.js フロントエンド
│   ├── src/
│   │   └── app/
│   │       ├── page.tsx       # メインページ
│   │       └── register/
│   │           └── page.tsx   # ユーザー登録ページ
│   ├── package.json
│   └── ...
└── meguru-backend/        # Go バックエンド
    ├── cmd/
    │   └── server/
    │       └── main.go    # メインサーバーファイル
    ├── internal/
    │   ├── domain/        # ドメイン層
    │   │   ├── domain_service/   # ドメインサービスのファイル（ドメインのビジネスロジック）
    │   │   ├── entity/           # ドメインモデル（エンティティのファイル）
    │   │   ├── repository/       # リポジトリのインターフェイスのファイル
    │   │   └── value_object/     # バリューオブジェクトのファイル
    │   ├── infrastructure/       # インフラストラクチャ層
    │   │   ├── query_service/    # クエリサービス（参照系のデータアクセス）
    │   │   └── repository/       # リポジトリ実装（登録・更新系のデータアクセス）
    │   ├── middleware/           # JWT認証などAPI共通処理のミドルウェア
    │   ├── presentation/         # プレゼンテーション層（APIハンドラー）
    │   │   └── handler/          # HTTPリクエストハンドラーのファイル
    │   ├── routes/               # ルーティング設定ファイル群
    │   └── usecase/              # ユースケース層
    │       ├── dto/              # データ転送用オブジェクト（DTO）
    │       ├── query_model/      # クエリモデル（読み取り用のデータ構造）
    │       ├── query_service/    # ユースケースから利用するクエリサービス
    │       └──                   # ユースケースの実装ファイル
    ├── pkg/
    │   └── database/             # データベース設定や接続ユーティリティ
    ├── scripts/
    │   └── db/
    │       └── migrations/       # データベースマイグレーションファイル
    ├── docker-compose.yml
    ├── Dockerfile
    ├── go.mod
    ├── go.sum
    ├── .env.example              # 環境変数のサンプルファイル
    └── .entrypoint.sh            # コンテナ起動時のエントリポイントスクリプト
```

## データベーススキーマ

```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
); 
```

## 使用技術

### フロントエンド
- Next.js 15
- TypeScript
- Tailwind CSS
- React Hooks

### バックエンド
- Go 1.23
- Gin Web Framework
- PostgreSQL 15
- bcrypt (パスワードハッシュ化)
- Docker & Docker Compose
- UUID (ユーザーID生成)

### アーキテクチャ
- オニオンアーキテクチャ (Clean Architecture)
- 依存関係逆転の原則
- CORS 対応

## オニオンアーキテクチャの詳細

### アーキテクチャ概要

このプロジェクトでは**オニオンアーキテクチャ**（Clean Architecture）を採用しています。このアーキテクチャの特徴は、外側の層が内側の層に依存し、内側の層は外側の層に依存しないことです。

```
┌─────────────────────────────────┐
│     Infrastructure Layer        │  ← 最外層
│  ┌─────────────────────────┐    │
│  │    Interface Layer      │    │
│  │  ┌─────────────────┐    │    │
│  │  │  Usecase Layer  │    │    │
│  │  │  ┌───────────┐  │    │    │
│  │  │  │  Domain   │  │    │    │  ← 最内層
│  │  │  │   Layer   │  │    │    │
│  │  │  └───────────┘  │    │    │
│  │  └─────────────────┘    │    │
│  └─────────────────────────┘    │
└─────────────────────────────────┘
```

### 各層の役割と責務

#### 1. Domain Layer（ドメイン層）- 最内層
**場所**: `internal/domain/`

**役割**: ビジネスロジックの核となる層で、外部に依存しない純粋なビジネスルールを定義

**構成要素**:
- **Entity** (`entity/user.go`): ビジネスオブジェクトの定義
  ```go
  type User struct {
      ID           uuid.UUID `json:"id"`
      Email        string    `json:"email"`
      PasswordHash string    `json:"-"`
      Name         string    `json:"name"`
      CreatedAt    time.Time `json:"created_at"`
      UpdatedAt    time.Time `json:"updated_at"`
  }
  ```

- **Repository Interface** (`repository/user_repository.go`): データアクセスの抽象化
  ```go
  type UserRepository interface {
      Create(ctx context.Context, user *entity.User) error
      GetByEmail(ctx context.Context, email string) (*entity.User, error)
      GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
  }
  ```

#### 2. Usecase Layer（ユースケース層）
**場所**: `internal/usecase/`

**役割**: アプリケーション固有のビジネスロジックを実装。ドメインエンティティを使用してビジネス要件を満たす

**責務**:
- バリデーション
- ビジネスルールの実行
- Repository インターフェースを通じたデータ操作
- レスポンスの構築

**実装例**:
```go
func (u *UserUsecase) CreateUser(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error) {
    // 1. 重複チェック
    existingUser, _ := u.userRepo.GetByEmail(ctx, req.Email)
    if existingUser != nil {
        return nil, errors.New("user with this email already exists")
    }
    
    // 2. パスワードハッシュ化
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    
    // 3. エンティティ作成
    user := &entity.User{
        ID:           uuid.New(),
        Email:        req.Email,
        PasswordHash: string(hashedPassword),
        // ...
    }
    
    // 4. データ保存
    if err := u.userRepo.Create(ctx, user); err != nil {
        return nil, err
    }
    
    return &CreateUserResponse{...}, nil
}
```

#### 3. Interface Layer（インターフェース層）
**場所**: `internal/interface/`

**役割**: 外部からの入力を受け取り、適切な形式に変換してUsecaseに渡す

**構成要素**:
- **Handler** (`handler/user_handler.go`): HTTPリクエストの処理
  - リクエストのパース
  - バリデーション
  - Usecaseの呼び出し
  - レスポンスの構築

**実装例**:
```go
func (c *UserHandler) CreateUser(ctx *gin.Context) {
    var req usecase.CreateUserRequest
    
    // リクエストバインディング
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // Usecaseの呼び出し
    response, err := c.userUsecase.CreateUser(ctx, &req)
    if err != nil {
        ctx.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    ctx.JSON(201, gin.H{"data": response})
}
```

#### 4. Infrastructure Layer（インフラストラクチャ層）- 最外層
**場所**: `internal/infrastructure/`

**役割**: 外部システム（データベース、API、ファイルシステムなど）との接続を担当

**構成要素**:
- **Database** (`database/user_repository.go`): Repository インターフェースの具体的な実装
- **Router** (`router/router.go`): HTTPルーティングの設定

### CQRS（コマンドクエリ分離）について
本設計では CQRS の考え方も取り入れています。CQRSとは「コマンド（Command：状態変更）」と「クエリ（Query：読み取り）」の処理を分離するパターンです。

#### Command（コマンド）】
- 新規作成・更新・削除など状態変更を伴う処理
- ユースケース層でドメインを操作し、書き込み用リポジトリを利用
- 例：CreateUser、UpdateProfile など

#### Query（クエリ）
- データの取得や表示用の処理で状態を変えない
- 読み取り専用のクエリサービスを使い、パフォーマンスを最適化
- ドメインモデルの形に必ずしも拘らず、ビューに最適化したデータ構造を利用することも多い

#### CQRSのメリット
- 読み書きで関心事を分離し、それぞれを最適化可能
- 読み取りのパフォーマンス改善やキャッシュ導入がしやすい
- ビジネスロジックの複雑化を防ぎ、テストも容易に

#### CQRSのイメージ
```
ユーザーリクエスト
 ├─ コマンド系（状態変更） → ユースケース（Command） → リポジトリ → DB書き込み
 └─ クエリ系（読み取り） → ユースケース（Query） → クエリサービス → DB読み取り
```

### データの流れ

ユーザー登録の場合のデータフローを例に説明します：

```
1. HTTP Request
   ↓
2. Router (Infrastructure)
   ↓
3. Handler (Interface) 
   - JSON パース
   - バリデーション
   ↓
4. Usecase (Application)
   - ビジネスロジック実行
   - Repository インターフェース使用
   ↓
5. Repository Implementation (Infrastructure)
   - SQL実行
   - データベースアクセス
   ↓
6. Database (External)

レスポンスは逆の流れで返される
```

### 依存関係の向き

```
Infrastructure → Interface → Usecase → Domain
      ↑                                   ↑
      └── Repository Implementation ──────┘
                    ↑
              Repository Interface
```

- **Domain**: 他のどの層にも依存しない
- **Usecase**: Domainのみに依存
- **Interface**: Usecase と Domain に依存
- **Infrastructure**: 全ての層に依存可能

### メリット

1. **テスタビリティ**: 各層を独立してテスト可能
2. **保守性**: 変更の影響範囲が限定される
3. **拡張性**: 新機能追加が容易
4. **技術選択の柔軟性**: データベースやフレームワークを変更しやすい
5. **ビジネスロジックの独立性**: フレームワークに依存しないビジネスロジック

### 実際のファイル構成

```
internal/
├── domain/                    # ドメイン層
│   ├── domain_service/        # ドメインサービス（エンティティだけで表現しきれないビジネスロジック）
│   ├── entity/                # エンティティ
│   │   └── user.go
│   ├── repository/            # リポジトリのインターフェース
│   │   └── user_repository.go
│   └── value_object/          # バリューオブジェクト
├── usecase/                   # ユースケース層
│   ├── dto/                   # DTO（入力・出力データ構造）
│   ├── query_model/           # クエリモデル（読み取り専用のデータ構造）
│   ├── query_service/         # クエリサービスのインターフェース
│   └── user_usecase.go
├── presentation/              # インターフェース層
│   └── handler/
│       └── user_handler.go
├── infrastructure/           # インフラストラクチャ層
│   ├── query_service/         # クエリサービス（参照系データアクセス）
│   ├── repository/            # データベース実装（書き込み系）
│   │   └── user_repository.go
│   └── middleware/            # JWT認証などのミドルウェア（共通処理）
└── routes/                   # ルーティングファイル
    └── router.go
```

## データベースマイグレーション
本プロジェクトでは、golang-migrate/migrate を利用してデータベースのスキーマ管理を行っています。

### マイグレーションファイルの管理
マイグレーションファイルは meguru-backend/scripts/db/migrations/ に配置しています。

ファイル名はタイムスタンプ + 処理内容を表す形式で管理しています。
```
例：
20250607120000_create_stores_table.up.sql
20250607120000_create_stores_table.down.sql
```

### マイグレーションの追加手順
1. 新しいマイグレーションファイルを作成（例: 20250607123000_add_column_to_users.up.sql、20250607123000_add_column_to_users.down.sql）

2. .up.sql に追加・変更のSQLを記述

3. .down.sql に取り消しのSQLを記述

4. マイグレーションを実行してスキーマを適用

### マイグレーションの適用（最新バージョンまで）
```
docker compose exec app migrate -path ./scripts/db/migrations -database "${DATABASE_URL}" up
```

### ロールバック
```
docker compose exec app migrate -path ./scripts/db/migrations -database "${DATABASE_URL}" down
```

## 今後の機能

- [ ] ユーザーログイン
- [ ] JWT認証
- [ ] ユーザープロフィール編集
- [ ] パスワードリセット
- [ ] メール認証
- [ ] API レート制限
- [ ] ログ機能 