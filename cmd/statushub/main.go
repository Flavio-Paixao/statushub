// Ponto de entrada do StatusHub: carrega a configuração, liga o verificador
// e sobe o servidor HTTP, desligando com calma quando recebe SIGTERM.
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Flavio-Paixao/statushub/internal/checker"
	"github.com/Flavio-Paixao/statushub/internal/config"
	"github.com/Flavio-Paixao/statushub/internal/server"
)

func main() {
	path := getenv("CONFIG_PATH", "config/targets.yaml")
	cfg, err := config.Load(path)
	if err != nil {
		log.Fatalf("configuração inválida: %v", err) // fail fast
	}
	log.Printf("config carregada de %s: %d alvos, intervalo %s, timeout %s",
		path, len(cfg.Targets), cfg.Interval, cfg.Timeout)

	// ctx é cancelado quando chega Ctrl+C (SIGINT) ou SIGTERM.
	// O Kubernetes manda SIGTERM antes de matar um Pod.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	chk := checker.New(cfg)
	go chk.Run(ctx) // "go" roda a função em paralelo (goroutine)

	srv := &http.Server{
		Addr:              ":" + getenv("PORT", "8080"),
		Handler:           server.New(chk),
		ReadHeaderTimeout: 5 * time.Second, // protege contra clientes lentos
	}

	go func() {
		log.Printf("StatusHub ouvindo em %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("servidor: %v", err)
		}
	}()

	<-ctx.Done() // fica parado aqui até chegar o sinal de desligar
	log.Println("desligando...")

	// Dá até 10s para terminar os pedidos em andamento (graceful shutdown).
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("erro no desligamento: %v", err)
	}
	log.Println("desligado")
}

// getenv lê uma variável de ambiente ou devolve o valor padrão.
func getenv(chave, padrao string) string {
	if v := os.Getenv(chave); v != "" {
		return v
	}
	return padrao
}
