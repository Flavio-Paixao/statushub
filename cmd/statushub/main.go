// Todo programa executável em Go começa no pacote "main".
package main

import (
	"log"      // escreve mensagens no terminal (com data e hora)
	"net/http" // servidor e cliente HTTP da biblioteca padrão
	"os"       // acesso a variáveis de ambiente

	"github.com/Flavio-Paixao/statushub/internal/config"
)

func main() {
	// O caminho do arquivo vem da variável de ambiente CONFIG_PATH.
	// No Kubernetes vamos apontar para o ConfigMap montado como arquivo.
	path := getenv("CONFIG_PATH", "config/targets.yaml")

	cfg, err := config.Load(path)
	if err != nil {
		// Fail fast: com configuração errada o serviço nem sobe.
		log.Fatalf("configuração inválida: %v", err)
	}
	log.Printf("config carregada de %s: %d alvos, intervalo %s, timeout %s",
		path, len(cfg.Targets), cfg.Interval, cfg.Timeout)
	for _, t := range cfg.Targets {
		log.Printf("  - %s (%s)", t.Name, t.URL)
	}

	// ServeMux é o "roteador": decide qual função atende cada caminho.
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", healthz)

	addr := ":" + getenv("PORT", "8080")
	log.Printf("StatusHub ouvindo em %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

// healthz responde se o processo está vivo (liveness).
func healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok\n"))
}

// getenv lê uma variável de ambiente ou devolve o valor padrão.
func getenv(chave, padrao string) string {
	if v := os.Getenv(chave); v != "" {
		return v
	}
	return padrao
}
