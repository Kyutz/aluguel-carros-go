package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/Kyutz/aluguel-carros-go/handlers"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

func main() {
	SetupDatabase()
	defer db.Close()

	http.HandleFunc("/clientes/criar", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handlers.ClienteCreateFormHandler(w, r)
		} else if r.Method == http.MethodPost {
			handlers.ClienteCreateHandler(db)(w, r)
		} else {
			http.Error(w, "Método não permitido", 405)
		}
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handlers.LoginFormHandler(w, r)
		} else if r.Method == http.MethodPost {
			handlers.LoginHandler(db)(w, r)
		} else {
			http.Error(w, "Método não permitido", 405)
		}
	})

	http.HandleFunc("/logout", handlers.LogoutHandler)

	log.Println("Servidor rodando na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
