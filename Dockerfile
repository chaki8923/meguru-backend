# ----- ビルドステージ -----
# Goの公式イメージをビルド環境として使用 (x86_64アーキテクチャを明示)
FROM --platform=linux/amd64 golang:1.23-alpine AS builder

WORKDIR /app

# GitをインストールしてからGoモジュールをダウンロード
RUN apk add --no-cache git

# アプリケーションのソースコードを全てコピー
COPY . .

# go.sumを生成して依存関係をダウンロード
RUN go mod tidy && go mod download

# migrate CLIをインストール
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.15.2

# ★★★★★ 修正点 ★★★★★
# App Runnerの実行環境(linux/amd64)向けにアプリケーションをビルドします。
# これにより `exec format error` が解消されます。
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/main ./cmd/server/main.go

# ----- 実行ステージ -----
# 軽量なAlpine Linuxを最終的な実行環境として使用 (x86_64アーキテクチャを明示)
FROM --platform=linux/amd64 alpine:latest

# SSL/TLS通信のためにルート証明書をインストール
RUN apk --no-cache add ca-certificates

WORKDIR /app

# ビルドステージからコンパイル済みの実行ファイルとマイグレーションファイルをコピー
COPY --from=builder /app/main .
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate

# マイグレーションファイルをコピー
COPY ./scripts/db/migrations ./migrations

# ★★★★★ 注意点 ★★★★★
# .env ファイルをイメージにコピーする処理は削除しました。
# データベースのパスワードなどの機密情報は、コンテナイメージに含めず、
# CloudFormationやApp Runnerのコンソールから設定する環境変数で渡すのが安全です。

# 起動スクリプトをコピーして実行権限を付与
COPY ./.entrypoint.sh .
RUN chmod +x .entrypoint.sh && \
    # 改行文字を確実にUnix形式に変換
    sed -i 's/\r$//' .entrypoint.sh

ENTRYPOINT ["./.entrypoint.sh"]

# アプリケーションがリッスンするポート
EXPOSE 8080

# コンテナ起動時に実行するコマンド
CMD ["./main"]