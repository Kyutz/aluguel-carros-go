package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Kyutz/aluguel-carros-go/models"
)

// GET /carros - listar todos (admin)
func ListarCarrosHandler(db *sql.DB) http.HandlerFunc {
	return AuthMiddleware(db, []string{"admin"}, func(w http.ResponseWriter, r *http.Request) {
		carros, err := models.GetAllCarros(db)
		if err != nil {
			http.Error(w, "Erro ao buscar carros", 500)
			return
		}
		json.NewEncoder(w).Encode(carros)
	})
}

// POST /carros - criar carro (admin)
func CriarCarroHandler(db *sql.DB) http.HandlerFunc {
	return AuthMiddleware(db, []string{"admin"}, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método não permitido", 405)
			return
		}

		var c models.Carro
		err := json.NewDecoder(r.Body).Decode(&c)
		if err != nil {
			http.Error(w, "JSON inválido: "+err.Error(), 400)
			return
		}

		c.Disponibilidade = true // default

		err = models.CreateCarro(db, c)
		if err != nil {
			http.Error(w, "Erro ao criar carro: "+err.Error(), 500)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"message":"Carro criado com sucesso"}`))
	})
}

// PUT /carros?id=123 - atualizar carro (admin)
func AtualizarCarroHandler(db *sql.DB) http.HandlerFunc {
	return AuthMiddleware(db, []string{"admin"}, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "Método não permitido", 405)
			return
		}
		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "ID inválido", 400)
			return
		}

		var c models.Carro
		err = json.NewDecoder(r.Body).Decode(&c)
		if err != nil {
			http.Error(w, "JSON inválido", 400)
			return
		}
		c.ID = id

		err = models.UpdateCarro(db, c)
		if err != nil {
			http.Error(w, "Erro ao atualizar carro", 500)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}

// POST /carros/deletar?id=123 - deletar carro (admin)
func DeletarCarroHandler(db *sql.DB) http.HandlerFunc {
	return AuthMiddleware(db, []string{"admin"}, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método não permitido", 405)
			return
		}
		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "ID inválido", 400)
			return
		}
		err = models.DeleteCarro(db, id)
		if err != nil {
			http.Error(w, "Erro ao deletar carro", 500)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}
