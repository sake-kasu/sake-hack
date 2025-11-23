.PHONY: help build run dev clean test test-unit test-integration cover lint deps api-validate api-generate api-bundle api-gendoc db-migrate-up db-migrate-down db-migrate-create sqlc-generate

# 変数定義
BINARY_NAME=sake-hack-server
MAIN_PATH=./cmd/server
BUILD_DIR=./bin
POSTGRES_DSN=postgresql://postgres:sakehacksakehack@localhost:5432/sake_hack_app?sslmode=disable
MYSQL_DSN=mysql://root:sakehacksakehack@localhost:3306/sake_hack_pts

help: ## このヘルプメッセージを表示
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# ビルド・実行
build: ## アプリケーションをビルド
	@echo "🔨 $(BINARY_NAME)をビルドしています..."
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

run: ## アプリケーションを実行
	@echo "🚀 $(BINARY_NAME)を実行しています..."
	@go run $(MAIN_PATH)/main.go

dev: ## ホットリロードで開発サーバーを起動(Air使用)
	@echo "🔥 開発サーバーを起動しています(ホットリロード有効)..."
	@air

clean: ## ビルド成果物をクリーンアップ
	@echo "🧹 クリーンアップ中..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html

# テスト・品質
test: ## 全てのテストを実行
	@echo "🧪 テストを実行しています..."
	@go test -v -race ./...

test-unit: ## 単体テストのみ実行
	@echo "⚡ 単体テストを実行しています..."
	@go test -v -short ./...

test-integration: ## 統合テストのみ実行
	@echo "🔗 統合テストを実行しています..."
	@go test -v -run Integration ./...

cover: ## カバレッジ測定付きでテストを実行(自動生成コード除外)
	@echo "📊 カバレッジを計測しています..."
	@go test -v -race -coverprofile=coverage.out -covermode=atomic $$(go list ./... | grep -v "/generated$$")
	@echo ""
	@echo "📈 全体カバレッジ:"
	@go tool cover -func=coverage.out | grep total | awk '{print "   Total Coverage: " $$3}'
	@echo ""
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✅ カバレッジレポートを生成しました: coverage.html"

lint: ## リンターを実行
	@echo "🔍 リンターを実行しています..."
	@golangci-lint run --timeout=5m ./...

gosec-scan: ## セキュリティスキャナーを実行
	@echo "🔒 セキュリティスキャンを実行しています..."
	@gosec -include=G109,G115 ./...

# 依存関係
deps: ## 依存関係を整理
	@echo "📦 依存関係を整理しています..."
	@go mod tidy
	@go mod download

# API開発(OpenAPI仕様から自動生成)
api-validate: ## OpenAPI仕様を検証
	@echo "✅ OpenAPI仕様を検証しています..."
	@npx @redocly/cli lint api/openapi.yaml --config api/redocly.yaml

api-generate: ## OpenAPI仕様からコードを自動生成
	@echo "🤖 OpenAPI仕様からコードを生成しています..."
	@echo "📦 Step 1: OpenAPI仕様をバンドルしています..."
	@npx @redocly/cli bundle api/openapi.yaml -o api/openapi.bundled.yaml
	@echo "⚙️  Step 2: Goコードを生成しています..."
	@mkdir -p api/generated
	@go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest \
		-config api/oapi-codegen.yaml api/openapi.bundled.yaml

api-bundle: ## OpenAPI仕様をバンドル
	@echo "📦 OpenAPI仕様をバンドルしています..."
	@npx @redocly/cli bundle api/openapi.yaml -o api/openapi.bundled.yaml

api-gendoc: ## APIドキュメントを生成
	@echo "📚 APIドキュメントを生成しています..."
	@npx @redocly/cli build-docs api/openapi.yaml -o api/docs/index.html

# データベースマイグレーション
db-migrate-up: ## データベースマイグレーションを実行(up)
	@echo "⬆️  マイグレーションを実行しています(up)..."
	@migrate -path db/migrations -database "$(POSTGRES_DSN)" up

db-migrate-down: ## データベースマイグレーションをロールバック(down)
	@echo "⬇️  マイグレーションをロールバックしています(down)..."
	@migrate -path db/migrations -database "$(POSTGRES_DSN)" down

db-migrate-create: ## 新規マイグレーションを作成(使用例: make db-migrate-create NAME=create_users)
	@echo "✨ マイグレーションを作成しています: $(NAME)"
	@migrate create -ext sql -dir db/migrations -seq $(NAME)

# sqlc
sqlc-generate: ## SQLからGoコードを生成
	@echo "🔧 SQLからGoコードを生成しています..."
	@sqlc generate
