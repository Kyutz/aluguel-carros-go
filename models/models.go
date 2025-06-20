package models

import "time"

type Cliente struct {
	ID                  int    `db:"id_cliente"`
	UserID              int    `db:"user_id"` // Relacionamento com usuário (planejado)
	Nome                string `db:"nome"`
	Email               string `db:"email"`
	Telefone            string `db:"telefone"`
	Endereco            string `db:"endereco"`
	DocumentoIdentidade string `db:"documento_identidade"`
}

type Carro struct {
	ID              int     `db:"id_carro"`
	Modelo          string  `db:"modelo"`
	Marca           string  `db:"marca"`
	Ano             int     `db:"ano"`
	Placa           string  `db:"placa"`
	Cor             string  `db:"cor"`
	Disponibilidade bool    `db:"disponibilidade"`
	ValorDiaria     float64 `db:"valor_diaria"`
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
