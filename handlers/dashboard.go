package handlers

import (
	"log"
	"net/http"
)

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	// Verifica cookie de sessão
	cookie, err := r.Cookie("session")
	if err != nil || cookie.Value == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Você pode passar dados para o template, como o nome do usuário
	data := struct {
		Usuario string
	}{
		Usuario: cookie.Value,
	}

	err = templates.ExecuteTemplate(w, "dashboard.html", data)
	if err != nil {
		log.Println("Erro executando template:", err)
		http.Error(w, "Erro interno", http.StatusInternalServerError)
	}
}
