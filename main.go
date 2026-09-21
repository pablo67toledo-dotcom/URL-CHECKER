package main

import (
	"fmt"
	"net/http"
	"net/url"
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

	fmt.Fprintln(w, "Domínio:", dominio)
}

// Função principal do programa
func main() {
	http.HandleFunc("/analisar", analisarURL)

	fmt.Println("Servidor rodando em http://localhost:8080")

	http.ListenAndServe(":8080", nil)
}
