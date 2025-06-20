package main

import (
	"log"
	"net/http"

	"github.com/Kyutz/aluguel-carros-go/handlers"
)

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

	log.Println("Servidor rodando na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
