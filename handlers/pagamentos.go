package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Kyutz/aluguel-carros-go/models"
)

// POST /pagamento - realiza pagamento (cliente)
func RealizarPagamentoHandler(db *sql.DB) http.HandlerFunc {
	return AuthMiddleware(db, []string{"cliente"}, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método não permitido", 405)
			return
		}
		var input struct {
			IDLocacao      int     `json:"id_locacao"`
			ValorPago      float64 `json:"valor_pago"`
			FormaPagamento string  `json:"forma_pagamento"`
		}
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			http.Error(w, "JSON inválido", 400)
			return
		}

		pagamento := models.Pagamento{
			IDLocacao:       input.IDLocacao,
			DataPagamento:   time.Now(),
			ValorPago:       input.ValorPago,
			FormaPagamento:  input.FormaPagamento,
			StatusPagamento: "confirmado",
		}

		err = models.CreatePagamento(db, pagamento)
		if err != nil {
			http.Error(w, "Erro ao salvar pagamento", 500)
			return
		}

		// Atualiza status da locação para 'pago'
		locacao, err := models.GetLocacaoByID(db, input.IDLocacao)
		if err == nil {
			locacao.Status = "pago"
			_ = models.UpdateLocacao(db, locacao)
		}

		w.WriteHeader(http.StatusCreated)
	})
}

// GET /pagamentos?id_cliente=123 - listar pagamentos do cliente
func PagamentosClienteHandler(db *sql.DB) http.HandlerFunc {
	return AuthMiddleware(db, []string{"cliente"}, func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Query().Get("id_cliente")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "ID inválido", 400)
			return
		}
		todos, err := models.GetAllPagamentos(db)
		if err != nil {
			http.Error(w, "Erro ao buscar pagamentos", 500)
			return
		}
		var meus []models.Pagamento
		for _, p := range todos {
			// buscar locação para verificar cliente
			locacao, err := models.GetLocacaoByID(db, p.IDLocacao)
			if err == nil && locacao.IDCliente == id {
				meus = append(meus, p)
			}
		}
		json.NewEncoder(w).Encode(meus)
	})
}
