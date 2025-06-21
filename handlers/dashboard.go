package handlers

import (
	"encoding/json"
	"net/http"
)

// DashboardHandler retorna dados básicos do usuário autenticado
func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil || cookie.Value == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Usuário não autenticado"})
		return
	}

	// Retorna o nome do usuário da sessão
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"usuario": cookie.Value,
		"msg":     "Bem-vindo ao dashboard",
	})
}
