package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Conecta ao banco SQLite (cria se não existir)
	db, err := sql.Open("sqlite3", "./aluguel_carros.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Testa conexão
	err = db.Ping()
	if err != nil {
		log.Fatal("Erro ao conectar no banco:", err)
	}

	// Ativa suporte a chaves estrangeiras (SQLite não ativa por padrão)
	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		log.Fatal("Erro ao ativar foreign keys:", err)
	}

	// Cria tabela Clientes
	createClientes := `
	CREATE TABLE IF NOT EXISTS clientes (
		id_cliente INTEGER PRIMARY KEY AUTOINCREMENT,
		nome TEXT NOT NULL,
		email TEXT,
		telefone TEXT,
		endereco TEXT,
		documento_identidade TEXT
	);`
	_, err = db.Exec(createClientes)
	if err != nil {
		log.Fatal("Erro criando tabela clientes:", err)
	}

	// Cria tabela Carros
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
	_, err = db.Exec(createCarros)
	if err != nil {
		log.Fatal("Erro criando tabela carros:", err)
	}

	// Cria tabela Locações
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
	_, err = db.Exec(createLocacoes)
	if err != nil {
		log.Fatal("Erro criando tabela locacoes:", err)
	}

	// Cria tabela Pagamentos
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
	_, err = db.Exec(createPagamentos)
	if err != nil {
		log.Fatal("Erro criando tabela pagamentos:", err)
	}

	log.Println("Tabelas criadas com sucesso!")
}
