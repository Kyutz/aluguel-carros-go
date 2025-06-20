package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/Kyutz/aluguel-carros-go/models"
)

// Lista clientes
func ClientesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id_cliente, nome, email, telefone, endereco, documento_identidade FROM clientes")
		if err != nil {
			http.Error(w, "Erro ao buscar clientes", 500)
			log.Println("Erro query clientes:", err)
			return
		}
		defer rows.Close()

		var clientes []models.Cliente
		for rows.Next() {
			var c models.Cliente
			err := rows.Scan(&c.ID, &c.Nome, &c.Email, &c.Telefone, &c.Endereco, &c.DocumentoIdentidade)
			if err != nil {
				http.Error(w, "Erro ao ler dados", 500)
				log.Println("Erro scan cliente:", err)
				return
			}
			clientes = append(clientes, c)
		}

		err = templates.ExecuteTemplate(w, "clientes.html", clientes)
		if err != nil {
			http.Error(w, "Erro ao renderizar template", 500)
			log.Println("Erro template clientes:", err)
		}
	}
}

// Formulário para criar cliente
func ClienteCreateFormHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "cliente_create.html", nil)
	if err != nil {
		http.Error(w, "Erro ao carregar formulário", 500)
	}
}

// Handler para salvar novo cliente
func ClienteCreateHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/clientes/criar", http.StatusSeeOther)
			return
		}

		nome := r.FormValue("nome")
		email := r.FormValue("email")
		telefone := r.FormValue("telefone")
		endereco := r.FormValue("endereco")
		documento := r.FormValue("documento_identidade")

		if nome == "" {
			http.Error(w, "Nome é obrigatório", 400)
			return
		}

		_, err := db.Exec(`INSERT INTO clientes (nome, email, telefone, endereco, documento_identidade) VALUES (?, ?, ?, ?, ?)`,
			nome, email, telefone, endereco, documento)
		if err != nil {
			log.Println("Erro inserindo cliente:", err)
			http.Error(w, "Erro ao salvar cliente", 500)
			return
		}

		http.Redirect(w, r, "/clientes", http.StatusSeeOther)
	}
}

func ClienteDeleteHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}

		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "ID do cliente é obrigatório", http.StatusBadRequest)
			return
		}

		_, err := db.Exec("DELETE FROM clientes WHERE id_cliente = ?", id)
		if err != nil {
			http.Error(w, "Erro ao deletar cliente", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/clientes", http.StatusSeeOther)
	}
}

func ClienteEditFormHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			http.Error(w, "ID do cliente não fornecido", http.StatusBadRequest)
			return
		}

		// Buscar cliente no DB pelo ID
		var c models.Cliente
		err := db.QueryRow("SELECT id_cliente, nome, email, telefone, endereco, documento_identidade FROM clientes WHERE id_cliente = ?", idStr).
			Scan(&c.ID, &c.Nome, &c.Email, &c.Telefone, &c.Endereco, &c.DocumentoIdentidade)
		if err != nil {
			http.Error(w, "Cliente não encontrado", http.StatusNotFound)
			return
		}

		// Renderizar template de edição com os dados do cliente
		err = templates.ExecuteTemplate(w, "cliente_edit.html", c)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func ClienteEditHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}

		id := r.FormValue("id")
		nome := r.FormValue("nome")
		email := r.FormValue("email")
		telefone := r.FormValue("telefone")
		endereco := r.FormValue("endereco")
		documento := r.FormValue("documento_identidade")

		_, err := db.Exec(`UPDATE clientes SET nome = ?, email = ?, telefone = ?, endereco = ?, documento_identidade = ? WHERE id_cliente = ?`,
			nome, email, telefone, endereco, documento, id)
		if err != nil {
			http.Error(w, "Erro ao atualizar cliente", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/clientes", http.StatusSeeOther)
	}
}
