package db

import (
	"database/sql"
	"fmt"
	"intellisys-backend/internal/configuration"

	_ "github.com/lib/pq"
)

func NewConnectionDB() (*sql.DB, error) {
	str := connectionStringPostgres()
	fmt.Println("Connection String:", str) // Linha de depuração
	db, err := sql.Open("postgres", connectionStringPostgres())

	if err != nil {
		return nil, fmt.Errorf("falha ao abrir conexão: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("não foi possível conectar ao banco: %v", err)
	}

	return db, nil
}

func connectionStringPostgres() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", configuration.Config.DBHost, configuration.Config.DBUser, configuration.Config.DBPassword, configuration.Config.DBName, configuration.Config.DBPort)
}
