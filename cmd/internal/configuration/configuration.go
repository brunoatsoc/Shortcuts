package configuration

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Estrutura para armazenar as configurações da aplicação
type Configuration struct {
	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     string
	ServerHost string
	ServerPort string
	JWTSecret  []byte
}

var Config = LoadEnvConfig() // Carrega as configurações do ambiente

// Função para carregar as variáveis de ambiente e retornar a configuração
func LoadEnvConfig() Configuration {
	err := godotenv.Load()

	if err != nil {
		log.Println("The .env file could not be loaded using system environment variables.")
	}

	return Configuration{
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		ServerHost: os.Getenv("SERVER_HOST"),
		ServerPort: os.Getenv("SERVER_PORT"),
		JWTSecret:  []byte(os.Getenv("JWT_SECRET")),
	}
} // Fim LoadEnvConfig
