package handlers

import (
	"database/sql"
	"encoding/json"
	"log" // Certifique-se de que 'log' está importado
	"net/http"
	"strconv"
	"time"

	"github.com/Kyutz/aluguel-carros-go/models"
)

// GET /carros/disponiveis - carros disponíveis (cliente)
func CarrosDisponiveisHandler(db *sql.DB) http.HandlerFunc {
	return AuthMiddleware(db, []string{"cliente"}, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet { // Adicionando verificação de método para consistência
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}

		carros, err := models.GetAllCarros(db)
		if err != nil {
			log.Printf("Erro ao buscar todos os carros para disponibilidade: %v", err)
			http.Error(w, "Erro interno ao buscar carros disponíveis", http.StatusInternalServerError)
			return
		}

		var disponiveis []models.Carro
		for _, c := range carros {
			if c.Disponibilidade {
				disponiveis = append(disponiveis, c)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(disponiveis)
	})
}

// POST /aluguel - criar locação (cliente)
func CriarLocacaoHandler(db *sql.DB) http.HandlerFunc {
	return AuthMiddleware(db, []string{"cliente"}, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}

		var l struct {
			IDCarro    int    `json:"id_carro"`
			IDCliente  int    `json:"id_cliente"`  // Reavaliar: idealmente, IDCliente viria da sessão
			DataInicio string `json:"data_inicio"` // formato AAAA-mm-dd
			DataFim    string `json:"data_fim"`
		}
		err := json.NewDecoder(r.Body).Decode(&l)
		if err != nil {
			log.Printf("Erro ao decodificar JSON para criar locação: %v", err)
			http.Error(w, "JSON inválido. Certifique-se de que todos os campos estão corretos.", http.StatusBadRequest)
			return
		}

		inicio, err := time.Parse("2006-01-02", l.DataInicio)
		if err != nil {
			log.Printf("Erro ao fazer parse da data de início '%s': %v", l.DataInicio, err)
			http.Error(w, "Data de início inválida. Use o formato AAAA-MM-DD.", http.StatusBadRequest)
			return
		}
		fim, err := time.Parse("2006-01-02", l.DataFim)
		if err != nil {
			log.Printf("Erro ao fazer parse da data de fim '%s': %v", l.DataFim, err)
			http.Error(w, "Data de fim inválida. Use o formato AAAA-MM-DD.", http.StatusBadRequest)
			return
		}

		if inicio.After(fim) {
			http.Error(w, "A data de início não pode ser depois da data de fim.", http.StatusBadRequest)
			return
		}

		carro, err := models.GetCarroByID(db, l.IDCarro)
		if err != nil {
			log.Printf("Erro ao buscar carro ID %d: %v", l.IDCarro, err)
			http.Error(w, "Carro não encontrado ou erro ao buscar.", http.StatusBadRequest)
			return
		}
		if !carro.Disponibilidade {
			http.Error(w, "Carro atualmente indisponível para locação.", http.StatusConflict) // Status 409 Conflict
			return
		}

		dias := int(fim.Sub(inicio).Hours()/24) + 1
		if dias <= 0 { // Garantir que a duração seja positiva
			http.Error(w, "Duração da locação inválida. Deve ser de pelo menos um dia.", http.StatusBadRequest)
			return
		}
		valor := float64(dias) * carro.ValorDiaria

		locacao := models.Locacao{
			IDCliente:  l.IDCliente, // Lembrete: Idealmente, IDCliente viria da sessão
			IDCarro:    l.IDCarro,
			DataInicio: inicio,
			DataFim:    fim,
			ValorTotal: valor,
			Status:     "pendente",
		}
		err = models.CreateLocacao(db, locacao)
		if err != nil {
			log.Printf("Erro ao criar locação no banco de dados para carro %d, cliente %d: %v", l.IDCarro, l.IDCliente, err)
			http.Error(w, "Erro interno ao registrar locação. Tente novamente.", http.StatusInternalServerError)
			return
		}

		// Atualiza carro para indisponível
		carro.Disponibilidade = false
		err = models.UpdateCarro(db, carro) // Capture o erro do UpdateCarro
		if err != nil {
			log.Printf("Atenção: Erro ao atualizar disponibilidade do carro %d após locação: %v", carro.ID, err)
			// Decide se isso deve ser um erro fatal para o cliente ou apenas um log
			// Por enquanto, não vamos retornar erro HTTP para o cliente se a locação foi criada.
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json") // Garante que a resposta é JSON
		json.NewEncoder(w).Encode(map[string]string{"message": "Locação criada com sucesso!"})
	})
}

// GET /minhas-locacoes?id_cliente=123 - locações do cliente
func MinhasLocacoesHandler(db *sql.DB) http.HandlerFunc {
	return AuthMiddleware(db, []string{"cliente"}, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet { // Adicionando verificação de método
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}

		idStr := r.URL.Query().Get("id_cliente")
		if idStr == "" {
			http.Error(w, "ID do cliente é obrigatório na URL.", http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Printf("Erro ao converter ID do cliente '%s' para inteiro: %v", idStr, err)
			http.Error(w, "ID do cliente inválido. Deve ser um número.", http.StatusBadRequest)
			return
		}

		// Melhoria futura: filtrar direto no banco de dados, não aqui em memória
		todas, err := models.GetAllLocacoes(db)
		if err != nil {
			log.Printf("Erro ao buscar todas as locações para filtrar por cliente %d: %v", id, err)
			http.Error(w, "Erro interno ao buscar suas locações.", http.StatusInternalServerError)
			return
		}
		var minhas []models.Locacao
		for _, l := range todas {
			if l.IDCliente == id {
				minhas = append(minhas, l)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(minhas)
	})
}
