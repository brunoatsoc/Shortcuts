package main

import (
	"intellisys-backend/api"
	"intellisys-backend/internal/db"
	"log"
)

// Função main que vai iniciar a aplicação
func main() {
	db, err := db.NewConnectionDB() // Cria conexão com o banco de dados

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	log.Println("Conexão com o banco de dados estabelecida com sucesso!")

	server := api.NewServerAPI(db) // Cria o servidor para a API

	if err := server.RunServer(); err != nil {
		log.Fatal(err)
	}
}
