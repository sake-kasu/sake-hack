package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLoad_LocalWithConfigFile はローカル環境でconfig.ymlを読み込むテスト
func TestLoad_LocalWithConfigFile(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, "config")
	err := os.MkdirAll(configDir, 0755)
	require.NoError(t, err)

	// テスト用のconfig.ymlを作成
	configContent := `
server:
  port: 9999
  mode: release
  gracefulShutdownTimeout: 60s

database:
  host: test-db
  port: 5433
  database: test_db
  user: test_user
  password: test_pass
  sslmode: require
  maxOpenConns: 50
  maxIdleConns: 10
  connMaxLifetime: 10m

valkey:
  host: test-valkey
  port: 6380
  password: test_valkey_pass
  database: 1
  poolSize: 20
  minIdleConns: 10
  maxRetries: 5
  dialTimeout: 10s
  readTimeout: 5s
  writeTimeout: 5s

jwt:
  secret: "test-secret"
  expiration: 3600
  cookieSecure: true
  cookieName: "test_token"
  cookiePath: "/test"
  cookieDomain: "test.com"

cors:
  allowedOrigins:
    - "http://test.com"
  allowedMethods:
    - "GET"
    - "POST"
  allowedHeaders:
    - "Content-Type"
  exposedHeaders:
    - "X-Total-Count"
  allowCredentials: false
  maxAge: 7200

logging:
  level: info
  format: console
`
	configPath := filepath.Join(configDir, "config.yaml")
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	// カレントディレクトリを変更
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() {
		err := os.Chdir(originalWd)
		require.NoError(t, err)
	}()
	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	// 環境変数をクリーンアップ
	t.Setenv("ENV", "local")

	// テスト実行
	cfg, err := Load()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// 検証
	assert.Equal(t, 9999, cfg.Server.Port)
	assert.Equal(t, "release", cfg.Server.Mode)
	assert.Equal(t, 60*time.Second, cfg.Server.GracefulShutdownTimeout)

	assert.Equal(t, "test-db", cfg.Database.Host)
	assert.Equal(t, 5433, cfg.Database.Port)
	assert.Equal(t, "test_db", cfg.Database.Database)
	assert.Equal(t, "test_user", cfg.Database.User)
	assert.Equal(t, "test_pass", cfg.Database.Password)
	assert.Equal(t, "require", cfg.Database.SSLMode)
	assert.Equal(t, 50, cfg.Database.MaxOpenConns)
	assert.Equal(t, 10, cfg.Database.MaxIdleConns)
	assert.Equal(t, 10*time.Minute, cfg.Database.ConnMaxLifetime)

	assert.Equal(t, "test-valkey", cfg.Valkey.Host)
	assert.Equal(t, 6380, cfg.Valkey.Port)
	assert.Equal(t, "test_valkey_pass", cfg.Valkey.Password)
	assert.Equal(t, 1, cfg.Valkey.Database)
	assert.Equal(t, 20, cfg.Valkey.PoolSize)
	assert.Equal(t, 10, cfg.Valkey.MinIdleConns)
	assert.Equal(t, 5, cfg.Valkey.MaxRetries)
	assert.Equal(t, 10*time.Second, cfg.Valkey.DialTimeout)
	assert.Equal(t, 5*time.Second, cfg.Valkey.ReadTimeout)
	assert.Equal(t, 5*time.Second, cfg.Valkey.WriteTimeout)

	assert.Equal(t, "test-secret", cfg.JWT.Secret)
	assert.Equal(t, 3600, cfg.JWT.Expiration)
	assert.True(t, cfg.JWT.CookieSecure)
	assert.Equal(t, "test_token", cfg.JWT.CookieName)
	assert.Equal(t, "/test", cfg.JWT.CookiePath)
	assert.Equal(t, "test.com", cfg.JWT.CookieDomain)

	assert.Equal(t, []string{"http://test.com"}, cfg.CORS.AllowedOrigins)
	assert.Equal(t, []string{"GET", "POST"}, cfg.CORS.AllowedMethods)
	assert.Equal(t, []string{"Content-Type"}, cfg.CORS.AllowedHeaders)
	assert.Equal(t, []string{"X-Total-Count"}, cfg.CORS.ExposedHeaders)
	assert.False(t, cfg.CORS.AllowCredentials)
	assert.Equal(t, 7200, cfg.CORS.MaxAge)

	assert.Equal(t, "info", cfg.Logging.Level)
	assert.Equal(t, "console", cfg.Logging.Format)
}

// TestLoad_LocalWithoutConfigFile はローカル環境でconfig.ymlがない場合のテスト
func TestLoad_LocalWithoutConfigFile(t *testing.T) {
	// テスト用の一時ディレクトリを作成(config.ymlなし)
	tmpDir := t.TempDir()

	// カレントディレクトリを変更
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() {
		err := os.Chdir(originalWd)
		require.NoError(t, err)
	}()
	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	// 環境変数をクリーンアップ
	t.Setenv("ENV", "local")

	// テスト実行(ファイルがなくてもエラーにならない)
	cfg, err := Load()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// デフォルト値が設定されていることを確認
	assert.Equal(t, 8080, cfg.Server.Port)
	assert.Equal(t, "debug", cfg.Server.Mode)
	assert.Equal(t, 30*time.Second, cfg.Server.GracefulShutdownTimeout)
}

// TestLoad_Production は本番環境で環境変数を使用するテスト
func TestLoad_Production(t *testing.T) {
	// 環境変数を設定
	t.Setenv("ENV", "production")
	t.Setenv("SERVER_PORT", "3000")
	t.Setenv("SERVER_MODE", "release")
	t.Setenv("DATABASE_HOST", "prod-db")
	t.Setenv("DATABASE_PORT", "5432")
	t.Setenv("VALKEY_HOST", "prod-valkey")

	// テスト実行
	cfg, err := Load()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// 環境変数の値が設定されていることを確認
	assert.Equal(t, 3000, cfg.Server.Port)
	assert.Equal(t, "release", cfg.Server.Mode)
	assert.Equal(t, "prod-db", cfg.Database.Host)
	assert.Equal(t, 5432, cfg.Database.Port)
	assert.Equal(t, "prod-valkey", cfg.Valkey.Host)
}

// TestLoad_EnvironmentVariableOverride は環境変数が設定ファイルを上書きするテスト
func TestLoad_EnvironmentVariableOverride(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, "config")
	err := os.MkdirAll(configDir, 0755)
	require.NoError(t, err)

	// テスト用のconfig.ymlを作成
	configContent := `
server:
  port: 8080
  mode: debug

database:
  host: localhost
  port: 5432
`
	configPath := filepath.Join(configDir, "config.yaml")
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	// カレントディレクトリを変更
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() {
		err := os.Chdir(originalWd)
		require.NoError(t, err)
	}()
	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	// 環境変数を設定(ファイルの値を上書き)
	t.Setenv("ENV", "local")
	t.Setenv("SERVER_PORT", "9000")
	t.Setenv("DATABASE_HOST", "override-db")

	// テスト実行
	cfg, err := Load()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// 環境変数が優先されることを確認
	assert.Equal(t, 9000, cfg.Server.Port)
	assert.Equal(t, "override-db", cfg.Database.Host)
	// ファイルの値はそのまま
	assert.Equal(t, "debug", cfg.Server.Mode)
	assert.Equal(t, 5432, cfg.Database.Port)
}

// TestLoad_DefaultENV はENV環境変数が未設定の場合のテスト
func TestLoad_DefaultENV(t *testing.T) {
	// ENVを明示的にアンセット
	t.Setenv("ENV", "")

	// テスト実行(デフォルトでlocalになる)
	cfg, err := Load()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// デフォルト値が設定されていることを確認
	assert.Equal(t, 8080, cfg.Server.Port)
}

// TestLoad_InvalidYAML は不正なYAMLファイルの場合のテスト
func TestLoad_InvalidYAML(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, "config")
	err := os.MkdirAll(configDir, 0755)
	require.NoError(t, err)

	// 不正なYAMLファイルを作成
	invalidContent := `
server:
  port: "invalid_port"  # ポートは数値であるべき
  gracefulShutdownTimeout: invalid_duration  # time.Durationのパースエラー
`
	configPath := filepath.Join(configDir, "config.yaml")
	err = os.WriteFile(configPath, []byte(invalidContent), 0644)
	require.NoError(t, err)

	// カレントディレクトリを変更
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() {
		err := os.Chdir(originalWd)
		require.NoError(t, err)
	}()
	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	// 環境変数を設定
	t.Setenv("ENV", "local")

	// テスト実行(アンマーシャルエラーが発生)
	cfg, err := Load()
	assert.Error(t, err)
	assert.Nil(t, cfg)
	assert.Contains(t, err.Error(), "設定のアンマーシャルに失敗しました")
}
