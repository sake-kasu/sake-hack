package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config はアプリケーション全体の設定
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Valkey   ValkeyConfig   `mapstructure:"valkey"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	CORS     CORSConfig     `mapstructure:"cors"`
	Logging  LoggingConfig  `mapstructure:"logging"`
}

// ServerConfig はサーバー設定
type ServerConfig struct {
	Port                    int           `mapstructure:"port"`
	Mode                    string        `mapstructure:"mode"`
	GracefulShutdownTimeout time.Duration `mapstructure:"gracefulShutdownTimeout"`
}

// DatabaseConfig はPostgreSQL設定
type DatabaseConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	Database        string        `mapstructure:"database"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	SSLMode         string        `mapstructure:"sslmode"`
	MaxOpenConns    int           `mapstructure:"maxOpenConns"`
	MaxIdleConns    int           `mapstructure:"maxIdleConns"`
	ConnMaxLifetime time.Duration `mapstructure:"connMaxLifetime"`
}

// ValkeyConfig はValkey設定
type ValkeyConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	Password     string        `mapstructure:"password"`
	Database     int           `mapstructure:"database"`
	PoolSize     int           `mapstructure:"poolSize"`
	MinIdleConns int           `mapstructure:"minIdleConns"`
	MaxRetries   int           `mapstructure:"maxRetries"`
	DialTimeout  time.Duration `mapstructure:"dialTimeout"`
	ReadTimeout  time.Duration `mapstructure:"readTimeout"`
	WriteTimeout time.Duration `mapstructure:"writeTimeout"`
}

// JWTConfig はJWT設定
type JWTConfig struct {
	Secret       string `mapstructure:"secret"`
	Expiration   int    `mapstructure:"expiration"`
	CookieSecure bool   `mapstructure:"cookieSecure"`
	CookieName   string `mapstructure:"cookieName"`
	CookiePath   string `mapstructure:"cookiePath"`
	CookieDomain string `mapstructure:"cookieDomain"`
}

// CORSConfig はCORS設定
type CORSConfig struct {
	AllowedOrigins   []string `mapstructure:"allowedOrigins"`
	AllowedMethods   []string `mapstructure:"allowedMethods"`
	AllowedHeaders   []string `mapstructure:"allowedHeaders"`
	ExposedHeaders   []string `mapstructure:"exposedHeaders"`
	AllowCredentials bool     `mapstructure:"allowCredentials"`
	MaxAge           int      `mapstructure:"maxAge"`
}

// LoggingConfig はロギング設定
type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

// Load は設定を読み込む
func Load() (*Config, error) {
	v := viper.New()

	// 環境変数の優先順位を設定
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// デフォルト値を設定
	setDefaults(v)

	// 環境判定
	env := os.Getenv("ENV")
	if env == "" {
		env = "local"
	}

	// ローカル環境の場合はYAMLファイルを読み込む
	if env == "local" {
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath("./config")
		v.AddConfigPath("../config")
		v.AddConfigPath("../../config")

		if err := v.ReadInConfig(); err != nil {
			// ファイルが存在しない場合は警告を出すが、エラーにはしない
			fmt.Printf("⚠️  設定ファイルが見つかりません: %v\n", err)
			fmt.Println("デフォルト値と環境変数を使用します")
		} else {
			fmt.Printf("✅ 設定ファイルを読み込みました: %s\n", v.ConfigFileUsed())
		}
	} else {
		fmt.Printf("✅ 環境: %s (環境変数を使用)\n", env)
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("設定のアンマーシャルに失敗しました: %w", err)
	}

	return &config, nil
}

// setDefaults はデフォルト値を設定する
func setDefaults(v *viper.Viper) {
	// Server
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.mode", "debug")
	v.SetDefault("server.gracefulShutdownTimeout", 30*time.Second)

	// Database (PostgreSQL)
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.database", "sake_hack_app")
	v.SetDefault("database.user", "postgres")
	v.SetDefault("database.password", "sakehacksakehack")
	v.SetDefault("database.sslmode", "disable")
	v.SetDefault("database.maxOpenConns", 25)
	v.SetDefault("database.maxIdleConns", 5)
	v.SetDefault("database.connMaxLifetime", 5*time.Minute)

	// Valkey
	v.SetDefault("valkey.host", "localhost")
	v.SetDefault("valkey.port", 6379)
	v.SetDefault("valkey.password", "sakehacksakehack")
	v.SetDefault("valkey.database", 0)
	v.SetDefault("valkey.poolSize", 10)
	v.SetDefault("valkey.minIdleConns", 5)
	v.SetDefault("valkey.maxRetries", 3)
	v.SetDefault("valkey.dialTimeout", 5*time.Second)
	v.SetDefault("valkey.readTimeout", 3*time.Second)
	v.SetDefault("valkey.writeTimeout", 3*time.Second)

	// JWT
	v.SetDefault("jwt.secret", "your-secret-key-change-me-in-production")
	v.SetDefault("jwt.expiration", 86400)
	v.SetDefault("jwt.cookieSecure", false)
	v.SetDefault("jwt.cookieName", "sake_hack_token")
	v.SetDefault("jwt.cookiePath", "/")
	v.SetDefault("jwt.cookieDomain", "")

	// CORS
	v.SetDefault("cors.allowedOrigins", []string{"http://localhost:3000", "http://localhost:8080"})
	v.SetDefault("cors.allowedMethods", []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"})
	v.SetDefault("cors.allowedHeaders", []string{"Origin", "Content-Type", "Accept", "Authorization"})
	v.SetDefault("cors.exposedHeaders", []string{"Content-Length"})
	v.SetDefault("cors.allowCredentials", true)
	v.SetDefault("cors.maxAge", 43200)

	// Logging
	v.SetDefault("logging.level", "debug")
	v.SetDefault("logging.format", "console")
}
