package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

func ClienteCreateFormHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "cliente_create.html", nil)
	if err != nil {
		http.Error(w, "Erro ao carregar formulário", 500)
	}
}

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
