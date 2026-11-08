// Todo programa executável em Go começa no pacote "main".
package main

import (
	"log"      // escreve mensagens no terminal (com data e hora)
	"net/http" // servidor e cliente HTTP da biblioteca padrão
)

func main() {
	// ServeMux é o "roteador": decide qual função atende cada caminho.
	mux := http.NewServeMux()

	// Desde o Go 1.22 dá para dizer o método junto com o caminho.
	// Só GET /healthz chega na função healthz; um POST recebe 405.
	mux.HandleFunc("GET /healthz", healthz)

	addr := ":8080" // porta 8080 em todas as interfaces
	log.Printf("StatusHub ouvindo em %s", addr)

	// ListenAndServe fica rodando para sempre e só retorna se der erro.
	// log.Fatal imprime o erro e encerra o programa.
	log.Fatal(http.ListenAndServe(addr, mux))
}

// healthz responde se o processo está vivo (liveness).
// w é onde escrevemos a resposta; r é o pedido que chegou.
func healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK) // código 200
	w.Write([]byte("ok\n"))      // corpo da resposta
}
