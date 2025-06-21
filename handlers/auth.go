package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// Função auxiliar que pega o papel do usuário pelo nome
func GetUserRole(db *sql.DB, username string) (string, error) {
	var papel string
	err := db.QueryRow("SELECT papel FROM usuarios WHERE usuario = ?", username).Scan(&papel)
	return papel, err
}

// Middleware que verifica se há cookie e se papel é permitido
func AuthMiddleware(db *sql.DB, allowedRoles []string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil || cookie.Value == "" {
			http.Error(w, "Não autorizado", http.StatusUnauthorized)
			return
		}

		papel, err := GetUserRole(db, cookie.Value)
		if err != nil {
			log.Println("Erro ao buscar papel do usuário:", err)
			http.Error(w, "Usuário não encontrado", http.StatusUnauthorized)
			return
		}

		permitido := false
		for _, ar := range allowedRoles {
			if ar == papel {
				permitido = true
				break
			}
		}

		if !permitido {
			http.Error(w, "Acesso negado", http.StatusForbidden)
			return
		}

		next(w, r)
	}
}

// Funções auxiliares para senhas
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// LoginJSONHandler realiza o login do usuário
func LoginJSONHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}

		var creds struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			http.Error(w, "Dados inválidos", http.StatusBadRequest)
			return
		}

		var hashedPassword string
		// Corrigindo a consulta SQL para buscar na coluna "senha_hash"
		err := db.QueryRow("SELECT senha_hash FROM usuarios WHERE usuario = ?", creds.Username).Scan(&hashedPassword)
		if err != nil {
			http.Error(w, "Usuário ou senha inválidos", http.StatusUnauthorized)
			return
		}

		if !CheckPasswordHash(creds.Password, hashedPassword) {
			http.Error(w, "Usuário ou senha inválidos", http.StatusUnauthorized)
			return
		}

		// Simula a criação de um cookie de sessão
		http.SetCookie(w, &http.Cookie{
			Name:  "session",
			Value: creds.Username,
			Path:  "/",
		})

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Login realizado com sucesso"))
	}
}

// LogoutJSONHandler realiza o logout do usuário
func LogoutJSONHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Remove o cookie de sessão
	http.SetCookie(w, &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logout realizado com sucesso"))
}
