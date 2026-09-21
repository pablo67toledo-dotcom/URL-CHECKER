package main

import (
	"fmt"
	"net/http"
)

// Função para analisar a URL recebida
func analisarURL(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	fmt.Fprintln(w, "URL recebida:", url)
}

// Função principal do programa
func main() {
	http.HandleFunc("/analisar", analisarURL)

	fmt.Println("Servidor rodando em http://localhost:8080")

	http.ListenAndServe(":8080", nil)
}
