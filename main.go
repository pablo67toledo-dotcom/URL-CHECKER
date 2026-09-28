package main

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/pablo67toledo-dotcom/URL_CHECKER/domain"
)

// Função para analisar a URL recebida
func analisarURL(w http.ResponseWriter, r *http.Request) {
	urlRecebida := r.URL.Query().Get("url")

	_, err := url.ParseRequestURI(urlRecebida)

	urlParsed, err := url.Parse(urlRecebida)

	if err != nil {
		http.Error(w, "Não foi possível analisar a URL", http.StatusBadRequest)
		return
	}
	dominio := urlParsed.Host

	err = domain.Validate(dominio)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resultado := domain.Check(dominio)

	fmt.Fprintln(w, resultado)
}

// Função principal do programa
func main() {
	http.HandleFunc("/analisar", analisarURL)

	fmt.Println("Servidor rodando em http://localhost:8080")

	http.ListenAndServe(":8080", nil)
}
