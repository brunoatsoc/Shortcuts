package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	FalePath FileSavePath
}

type ServerConfig struct {
	Port string
	Host string
}

type DatabaseConfig struct {
	URL string
}

type JWTConfig struct {
	PrivateKeyPath string
	PublicKeyPath  string
	ExpireTime     int64
}

type FileSavePath struct {
	Path string
}

func LoadConfig() (*Config, error) {
	// Ignora erro se .env não existe (Docker injeta as variáveis via environment)
	_ = godotenv.Load()

	expireTime, err := strconv.ParseInt(os.Getenv("EXPIRATION_TIME_MINUTES"), 10, 64)

	if err != nil {
		return nil, fmt.Errorf("erro ao converter EXPIRATION_TIME_MINUTES para int64: %w", err)
	}

	return &Config{
		Server: ServerConfig{
			Port: os.Getenv("SERVER_PORT"),
			Host: os.Getenv("SERVER_HOST"),
		},
		Database: DatabaseConfig{
			URL: os.Getenv("DB_URL"),
		},
		JWT: JWTConfig{
			PrivateKeyPath: os.Getenv("JWT_PRIVATE_KEY_PATH"),
			PublicKeyPath:  os.Getenv("JWT_PUBLIC_KEY_PATH"),
			ExpireTime:     expireTime,
		},
		FalePath: FileSavePath{
			Path: os.Getenv("FILE_PATH"),
		},
	}, nil
}
