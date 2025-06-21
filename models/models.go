package models

import (
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Usuario struct {
	ID           int    `db:"id_usuario"`
	Username     string `db:"username"`
	PasswordHash string `db:"password_hash"`
	Papel        string `db:"papel"`
}

type Cliente struct {
	ID                  int    `db:"id_cliente" json:"id"`
	Nome                string `db:"nome" json:"nome"`
	Email               string `db:"email" json:"email"`
	Telefone            string `db:"telefone" json:"telefone"`
	Endereco            string `db:"endereco" json:"endereco"`
	DocumentoIdentidade string `db:"documento_identidade" json:"documento_identidade"`
	Username            string `db:"username" json:"username"`
}

type Carro struct {
	ID              int     `db:"id_carro" json:"id"`
	Modelo          string  `db:"modelo" json:"modelo"`
	Marca           string  `db:"marca" json:"marca"`
	Ano             int     `db:"ano" json:"ano"`
	Placa           string  `db:"placa" json:"placa"`
	Cor             string  `db:"cor" json:"cor"`
	Disponibilidade bool    `db:"disponibilidade" json:"disponibilidade"`
	ValorDiaria     float64 `db:"valor_diaria" json:"valor_diaria"`
}

type Locacao struct {
	ID         int       `db:"id_locacao"`
	IDCliente  int       `db:"id_cliente"`
	IDCarro    int       `db:"id_carro"`
	DataInicio time.Time `db:"data_inicio"`
	DataFim    time.Time `db:"data_fim"`
	ValorTotal float64   `db:"valor_total"`
	Status     string    `db:"status"`
}

// Você calcularia ValorTotal no código Go antes de salvar a Locacao

type Pagamento struct {
	ID              int       `db:"id_pagamento"`
	IDLocacao       int       `db:"id_locacao"`
	DataPagamento   time.Time `db:"data_pagamento"`
	ValorPago       float64   `db:"valor_pago"`
	FormaPagamento  string    `db:"forma_pagamento"`
	StatusPagamento string    `db:"status_pagamento"`
}

// --- Cliente ---
// Função para gerar o hash da senha
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func GetAllClientes(db *sql.DB) ([]Cliente, error) {
	rows, err := db.Query("SELECT id_cliente, nome, email, telefone, endereco, documento_identidade, username FROM clientes") // Adicione 'username' aqui
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clientes []Cliente
	for rows.Next() {
		var c Cliente
		err := rows.Scan(&c.ID, &c.Nome, &c.Email, &c.Telefone, &c.Endereco, &c.DocumentoIdentidade, &c.Username) // Adicione &c.Username
		if err != nil {
			return nil, err
		}
		clientes = append(clientes, c)
	}
	return clientes, nil
}

func GetClienteByID(db *sql.DB, id int) (Cliente, error) {
	var c Cliente
	err := db.QueryRow("SELECT id_cliente, nome, email, telefone, endereco, documento_identidade, username FROM clientes WHERE id_cliente = ?", id). // Adicione 'username' aqui
																				Scan(&c.ID, &c.Nome, &c.Email, &c.Telefone, &c.Endereco, &c.DocumentoIdentidade, &c.Username) // Adicione &c.Username
	return c, err
}

// Em models.go
func CreateCliente(db *sql.DB, c Cliente, senha string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	res, err := tx.Exec(`INSERT INTO clientes (nome, email, telefone, endereco, documento_identidade, username)
		VALUES (?, ?, ?, ?, ?, ?)`,
		c.Nome, c.Email, c.Telefone, c.Endereco, c.DocumentoIdentidade, c.Username)
	if err != nil {
		tx.Rollback()
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}
	c.ID = int(id)

	senhaHash, err := HashPassword(senha)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`INSERT INTO usuarios (usuario, senha_hash, papel) VALUES (?, ?, ?)`,
		c.Username, senhaHash, "cliente")
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
func UpdateCliente(db *sql.DB, c Cliente) error {
	_, err := db.Exec(`UPDATE clientes SET nome=?, email=?, telefone=?, endereco=?, documento_identidade=?, username=?
		WHERE id_cliente=?`,
		c.Nome, c.Email, c.Telefone, c.Endereco, c.DocumentoIdentidade, c.Username, c.ID)
	return err
}

func DeleteCliente(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM clientes WHERE id_cliente = ?", id)
	return err
}

// --- Carro ---

func GetAllCarros(db *sql.DB) ([]Carro, error) {
	rows, err := db.Query("SELECT id_carro, modelo, marca, ano, placa, cor, disponibilidade, valor_diaria FROM carros")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var carros []Carro
	for rows.Next() {
		var c Carro
		err := rows.Scan(&c.ID, &c.Modelo, &c.Marca, &c.Ano, &c.Placa, &c.Cor, &c.Disponibilidade, &c.ValorDiaria)
		if err != nil {
			return nil, err
		}
		carros = append(carros, c)
	}
	return carros, nil
}

func GetCarroByID(db *sql.DB, id int) (Carro, error) {
	var c Carro
	err := db.QueryRow("SELECT id_carro, modelo, marca, ano, placa, cor, disponibilidade, valor_diaria FROM carros WHERE id_carro = ?", id).
		Scan(&c.ID, &c.Modelo, &c.Marca, &c.Ano, &c.Placa, &c.Cor, &c.Disponibilidade, &c.ValorDiaria)
	return c, err
}

func CreateCarro(db *sql.DB, c Carro) error {
	_, err := db.Exec(`INSERT INTO carros (modelo, marca, ano, placa, cor, disponibilidade, valor_diaria)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		c.Modelo, c.Marca, c.Ano, c.Placa, c.Cor, c.Disponibilidade, c.ValorDiaria)
	return err
}

func UpdateCarro(db *sql.DB, c Carro) error {
	_, err := db.Exec(`UPDATE carros SET modelo=?, marca=?, ano=?, placa=?, cor=?, disponibilidade=?, valor_diaria=?
		WHERE id_carro=?`,
		c.Modelo, c.Marca, c.Ano, c.Placa, c.Cor, c.Disponibilidade, c.ValorDiaria, c.ID)
	return err
}

func DeleteCarro(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM carros WHERE id_carro = ?", id)
	return err
}

// --- Locacao ---

func GetAllLocacoes(db *sql.DB) ([]Locacao, error) {
	rows, err := db.Query("SELECT id_locacao, id_cliente, id_carro, data_inicio, data_fim, valor_total, status FROM locacoes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locacoes []Locacao
	for rows.Next() {
		var l Locacao
		err := rows.Scan(&l.ID, &l.IDCliente, &l.IDCarro, &l.DataInicio, &l.DataFim, &l.ValorTotal, &l.Status)
		if err != nil {
			return nil, err
		}
		locacoes = append(locacoes, l)
	}
	return locacoes, nil
}

func GetLocacaoByID(db *sql.DB, id int) (Locacao, error) {
	var l Locacao
	err := db.QueryRow("SELECT id_locacao, id_cliente, id_carro, data_inicio, data_fim, valor_total, status FROM locacoes WHERE id_locacao = ?", id).
		Scan(&l.ID, &l.IDCliente, &l.IDCarro, &l.DataInicio, &l.DataFim, &l.ValorTotal, &l.Status)
	return l, err
}

func CreateLocacao(db *sql.DB, l Locacao) error {
	_, err := db.Exec(`INSERT INTO locacoes (id_cliente, id_carro, data_inicio, data_fim, valor_total, status)
		VALUES (?, ?, ?, ?, ?, ?)`,
		l.IDCliente, l.IDCarro, l.DataInicio, l.DataFim, l.ValorTotal, l.Status)
	return err
}

func UpdateLocacao(db *sql.DB, l Locacao) error {
	_, err := db.Exec(`UPDATE locacoes SET id_cliente=?, id_carro=?, data_inicio=?, data_fim=?, valor_total=?, status=?
		WHERE id_locacao=?`,
		l.IDCliente, l.IDCarro, l.DataInicio, l.DataFim, l.ValorTotal, l.Status, l.ID)
	return err
}

func DeleteLocacao(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM locacoes WHERE id_locacao = ?", id)
	return err
}

// --- Pagamento ---

func GetAllPagamentos(db *sql.DB) ([]Pagamento, error) {
	rows, err := db.Query("SELECT id_pagamento, id_locacao, data_pagamento, valor_pago, forma_pagamento, status_pagamento FROM pagamentos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pagamentos []Pagamento
	for rows.Next() {
		var p Pagamento
		err := rows.Scan(&p.ID, &p.IDLocacao, &p.DataPagamento, &p.ValorPago, &p.FormaPagamento, &p.StatusPagamento)
		if err != nil {
			return nil, err
		}
		pagamentos = append(pagamentos, p)
	}
	return pagamentos, nil
}

func GetPagamentoByID(db *sql.DB, id int) (Pagamento, error) {
	var p Pagamento
	err := db.QueryRow("SELECT id_pagamento, id_locacao, data_pagamento, valor_pago, forma_pagamento, status_pagamento FROM pagamentos WHERE id_pagamento = ?", id).
		Scan(&p.ID, &p.IDLocacao, &p.DataPagamento, &p.ValorPago, &p.FormaPagamento, &p.StatusPagamento)
	return p, err
}

func CreatePagamento(db *sql.DB, p Pagamento) error {
	_, err := db.Exec(`INSERT INTO pagamentos (id_locacao, data_pagamento, valor_pago, forma_pagamento, status_pagamento)
		VALUES (?, ?, ?, ?, ?)`,
		p.IDLocacao, p.DataPagamento, p.ValorPago, p.FormaPagamento, p.StatusPagamento)
	return err
}

func UpdatePagamento(db *sql.DB, p Pagamento) error {
	_, err := db.Exec(`UPDATE pagamentos SET id_locacao=?, data_pagamento=?, valor_pago=?, forma_pagamento=?, status_pagamento=?
		WHERE id_pagamento=?`,
		p.IDLocacao, p.DataPagamento, p.ValorPago, p.FormaPagamento, p.StatusPagamento, p.ID)
	return err
}

func DeletePagamento(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM pagamentos WHERE id_pagamento = ?", id)
	return err
}
