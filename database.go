package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func SetupDatabase() {
	var err error
	db, err = sql.Open("sqlite3", "./aluguel_carros.db")
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Erro ao conectar no banco:", err)
	}

	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		log.Fatal("Erro ao ativar foreign keys:", err)
	}

	createClientes := `
	CREATE TABLE IF NOT EXISTS clientes (
		id_cliente INTEGER PRIMARY KEY AUTOINCREMENT,
		nome TEXT NOT NULL,
		email TEXT,
		telefone TEXT,
		endereco TEXT,
		documento_identidade TEXT
	);`

	createCarros := `
	CREATE TABLE IF NOT EXISTS carros (
		id_carro INTEGER PRIMARY KEY AUTOINCREMENT,
		modelo TEXT NOT NULL,
		marca TEXT,
		ano INTEGER,
		placa TEXT UNIQUE,
		cor TEXT,
		disponibilidade BOOLEAN NOT NULL
	);`

	createLocacoes := `
	CREATE TABLE IF NOT EXISTS locacoes (
		id_locacao INTEGER PRIMARY KEY AUTOINCREMENT,
		id_cliente INTEGER,
		id_carro INTEGER,
		data_inicio DATE,
		data_fim DATE,
		valor_total REAL,
		status TEXT,
		FOREIGN KEY (id_cliente) REFERENCES clientes(id_cliente),
		FOREIGN KEY (id_carro) REFERENCES carros(id_carro)
	);`

	createPagamentos := `
	CREATE TABLE IF NOT EXISTS pagamentos (
		id_pagamento INTEGER PRIMARY KEY AUTOINCREMENT,
		id_locacao INTEGER,
		data_pagamento DATE,
		valor_pago REAL,
		forma_pagamento TEXT,
		status_pagamento TEXT,
		FOREIGN KEY (id_locacao) REFERENCES locacoes(id_locacao)
	);`

	for _, query := range []string{createClientes, createCarros, createLocacoes, createPagamentos} {
		_, err = db.Exec(query)
		if err != nil {
			log.Fatal("Erro criando tabela:", err)
		}
	}

	log.Println("Tabelas criadas com sucesso!")
}
