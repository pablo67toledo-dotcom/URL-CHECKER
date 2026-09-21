package main

import (
	"fmt"
	"net/http"
)
// Função para analisar a ULR recebida
func analisarURL(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Servidor recebeu uma URL")
}
// Função principal do programa
func main (){
	http.HandleFunc("/analisar", analisarURL)

	fmt.Println("Servidor rodando em http://localhost:8080")

	http.ListenAndServe(":8080", nil)
}
