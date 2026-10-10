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

	if err != nil {
		http.Error(w, "URL inválida", http.StatusBadRequest)
		return
	}

	dominio, err := domain.ExtractDomain(urlRecebida)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusAccepted)
		return
	}

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

}

// Função principal do programa
func main() {
	http.HandleFunc("/analisar", analisarURL)

	fmt.Println("Servidor rodando em http://localhost:8080")

	http.ListenAndServe(":8080", nil)
}
