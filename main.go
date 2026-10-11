package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/pablo67toledo-dotcom/URL_CHECKER/database"
	"github.com/pablo67toledo-dotcom/URL_CHECKER/domain"
)

// analisarURL analisa a URL recebida e verifica seu domínio.
func analisarURL(w http.ResponseWriter, r *http.Request) {
	urlRecebida := r.URL.Query().Get("url")

	_, err := url.ParseRequestURI(urlRecebida)
	if err != nil {
		http.Error(w, "URL inválida", http.StatusBadRequest)
		return
	}

	dominio, err := domain.ExtractDomain(urlRecebida)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = domain.Check(dominio)
	if err != nil {
		http.Error(w, "Domínio inválido", http.StatusBadRequest)
		return
	}

	redirectURL, err := domain.CheckRedirect(urlRecebida)
	if err != nil {
		http.Error(w, "Não foi possível verificar o redirecionamento", http.StatusBadRequest)
		return
	}

	if redirectURL != "" {
		fmt.Fprintln(w, "Redirecionou para:", redirectURL)
	} else {
		fmt.Fprintln(w, "Não houve redirecionamento.")
	}

	fmt.Fprintln(w, "Domínio encontrado:", dominio)

	// Consulta a data de criação do domínio.
	dataCriacao, err := domain.GetDomainCreationDate(dominio)
	if err != nil {
		fmt.Fprintln(w, "Não foi possível consultar a idade do domínio:", err)
		return
	}

	// Calcula a idade do domínio em anos completos.
	idade := domain.CalculateDomainAge(dataCriacao)

	fmt.Fprintln(w, "Data de criação:", dataCriacao.Format("02/01/2006"))
	fmt.Fprintf(w, "Idade do domínio: %d anos\n", idade)
}

// main inicializa o banco de dados e registra as rotas HTTP.
func main() {
	// Carrega as variáveis do arquivo .env.
	if err := godotenv.Load(); err != nil {
		fmt.Println("Arquivo .env não carregado; usando variáveis do ambiente")
	}

	// Define um limite de tempo para iniciar a conexão.
	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	// Conecta ao PostgreSQL.
	pool, err := database.NewPool(ctx)
	if err != nil {
		fmt.Println("Erro ao conectar ao PostgreSQL:", err)
		os.Exit(1)
	}
	defer pool.Close()

	fmt.Println("Conexão com PostgreSQL estabelecida.")

	// Registra as rotas uma única vez, durante a inicialização.
	http.HandleFunc("/analisar", analisarURL)
	http.HandleFunc("/brands", brandHandler(pool))

	fmt.Println("Servidor rodando em http://localhost:8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Erro ao iniciar o servidor:", err)
	}
}
