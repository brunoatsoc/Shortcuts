package main

import (
	"intellisys-backend/cmd/api"
	"intellisys-backend/db"
	"log"
)

func main() {
	db, err := db.NewConnectionDB()

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	log.Println("Conexão com o banco de dados estabelecida com sucesso!")

	server := api.NewServerAPI(db)

	if err := server.RunServer(); err != nil {
		log.Fatal(err)
	}
}
