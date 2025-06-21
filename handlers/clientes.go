package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Kyutz/aluguel-carros-go/models"
)

// Middleware simples para checar sessão via cookie "session"
func checkSession(w http.ResponseWriter, r *http.Request) bool {
	cookie, err := r.Cookie("session")
	if err != nil || cookie.Value == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Usuário não autenticado"})
		return false
	}
	return true
}

// Listar todos clientes (GET /clientes)
func ClientesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !checkSession(w, r) {
			return
		}

		clientes, err := models.GetAllClientes(db)
		if err != nil {
			log.Println("Erro buscando clientes:", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Erro ao buscar clientes"})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(clientes)
	}
}

// Criar novo cliente (POST /clientes)
// Em handlers/clientes.go
func ClienteCreateHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Verificação de sessão (middleware)
		if !checkSession(w, r) {
			// checkSession já envia a resposta de erro (401 Unauthorized)
			return
		}

		// 2. Verificação do método HTTP
		if r.Method != http.MethodPost {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}

		var input struct {
			models.Cliente
			Senha    string `json:"senha"`
			Username string `json:"username"` // Capturar username da entrada JSON
		}

		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "JSON inválido"})
			return
		}

		// Atribua o username capturado à struct Cliente
		input.Cliente.Username = input.Username

		if input.Nome == "" || input.Senha == "" || input.Username == "" { // Também verificar Username
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Nome, senha e nome de usuário são obrigatórios"})
			return
		}

		err = models.CreateCliente(db, input.Cliente, input.Senha)
		if err != nil {
			log.Println("Erro criando cliente:", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Erro ao salvar cliente"})
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Cliente criado com sucesso"})
	}
}

// Deletar cliente (DELETE /clientes?id=)
func ClienteDeleteHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !checkSession(w, r) {
			return
		}

		if r.Method != http.MethodDelete {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{"error": "Método não permitido"})
			return
		}

		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "ID do cliente é obrigatório"})
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "ID inválido"})
			return
		}

		err = models.DeleteCliente(db, id)
		if err != nil {
			log.Println("Erro deletando cliente:", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Erro ao deletar cliente"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Cliente deletado com sucesso"})
	}
}

// Editar cliente (PUT /clientes?id=)
func ClienteEditHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !checkSession(w, r) {
			return
		}

		if r.Method != http.MethodPut {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{"error": "Método não permitido"})
			return
		}

		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "ID do cliente é obrigatório"})
			return
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "ID inválido"})
			return
		}

		var c models.Cliente
		err = json.NewDecoder(r.Body).Decode(&c)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "JSON inválido"})
			return
		}

		c.ID = id

		err = models.UpdateCliente(db, c)
		if err != nil {
			log.Println("Erro atualizando cliente:", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Erro ao atualizar cliente"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Cliente atualizado com sucesso"})
	}
}
