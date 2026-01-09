package db

import (
	"database/sql"
	"fmt"
	"intellisys-backend/internal/configuration"

	_ "github.com/lib/pq"
)

// Função para criar uma nova conexão com o banco de dados
func NewConnectionDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", connectionStringPostgres())

	if err != nil {
		return nil, fmt.Errorf("falha ao abrir conexão: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("não foi possível conectar ao banco: %v", err)
	}

	return db, nil
} // Fim NewConnectionDB

// Função para montar a string de conexão com o Postgres
func connectionStringPostgres() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", configuration.Config.DBHost, configuration.Config.DBUser, configuration.Config.DBPassword, configuration.Config.DBName, configuration.Config.DBPort)
} // Fim connectionStringPostgres
