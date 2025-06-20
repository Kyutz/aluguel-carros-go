package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./aluguel_carros.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	usuario := "admin"
	senhaHash := "$2a$10$UXIpKQbqqAW5slGQHlaHOOI7hXTwuNXEY7gWslfn2zgfmh6sdkqCq"
	papel := "admin"

	_, err = db.Exec("INSERT INTO usuarios (usuario, senha_hash, papel) VALUES (?, ?, ?)", usuario, senhaHash, papel)
	if err != nil {
		log.Fatal("Erro inserindo usuário admin:", err)
	}

	log.Println("Usuário admin criado com sucesso!")
}
