package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

func GerarHashSenhaTemporario() {
	senha := "admin123"
	hash, err := bcrypt.GenerateFromPassword([]byte(senha), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Hash da senha:", string(hash))
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ---------- Handlers de Login ----------

// GET /login - Mostra o formulário
func LoginFormHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}
	err := templates.ExecuteTemplate(w, "login.html", nil)
	if err != nil {
		http.Error(w, "Erro ao carregar a página", http.StatusInternalServerError)
	}
}

// POST /login - Processa o login
func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}

		usuario := r.FormValue("usuario")
		senha := r.FormValue("senha")

		var hashSenha string
		err := db.QueryRow("SELECT senha_hash FROM usuarios WHERE usuario = ?", usuario).Scan(&hashSenha)
		if err != nil {
			http.Error(w, "Usuário ou senha inválidos", http.StatusUnauthorized)
			return
		}

		if !CheckPasswordHash(senha, hashSenha) {
			http.Error(w, "Usuário ou senha inválidos", http.StatusUnauthorized)
			return
		}

		// Cria um cookie simples para manter sessão
		cookie := http.Cookie{
			Name:  "session",
			Value: usuario,
			Path:  "/",
			// Pode adicionar Secure: true e HttpOnly: true no HTTPS
		}
		http.SetCookie(w, &cookie)

		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}
}

// GET /logout - Faz logout
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1, // Faz expirar
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
