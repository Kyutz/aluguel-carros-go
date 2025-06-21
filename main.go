package main

import (
	"log"
	"net/http"

	"github.com/Kyutz/aluguel-carros-go/handlers"
)

func main() {
	SetupDatabase()
	defer db.Close()

	// Autenticação
	http.HandleFunc("/login", handlers.LoginJSONHandler(db)) // POST /login
	http.HandleFunc("/logout", handlers.LogoutJSONHandler)   // GET /logout

	// CRUD de carros
	http.HandleFunc("/carros", handlers.ListarCarrosHandler(db))             // GET
	http.HandleFunc("/carros/criar", handlers.CriarCarroHandler(db))         // POST
	http.HandleFunc("/carros/atualizar", handlers.AtualizarCarroHandler(db)) // PUT (emulado via POST)
	http.HandleFunc("/carros/deletar", handlers.DeletarCarroHandler(db))     // POST (emulando DELETE)

	// Cliente
	http.HandleFunc("/clientes", handlers.ClientesHandler(db))            // GET
	http.HandleFunc("/clientes/criar", handlers.ClienteCreateHandler(db)) // POST
	http.HandleFunc("/clientes/editar", handlers.ClienteEditHandler(db))  // POST

	// Aluguel
	http.HandleFunc("/carros/disponiveis", handlers.CarrosDisponiveisHandler(db)) // GET
	http.HandleFunc("/aluguel", handlers.CriarLocacaoHandler(db))                 // POST
	http.HandleFunc("/minhas-locacoes", handlers.MinhasLocacoesHandler(db))       // GET

	// Pagamento
	http.HandleFunc("/pagamento", handlers.RealizarPagamentoHandler(db))  // POST
	http.HandleFunc("/pagamentos", handlers.PagamentosClienteHandler(db)) // GET

	log.Println("Servidor rodando na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
