// Package server define as rotas HTTP do StatusHub.
package server

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/Flavio-Paixao/statushub/internal/checker"
)

// Source é o que o servidor precisa saber do verificador.
// Usar uma interface (e não o Checker direto) permite testar com dados falsos.
type Source interface {
	Results() []checker.Result
	Ready() bool
}

// New monta o roteador com todas as rotas.
func New(src Source) http.Handler {
	mux := http.NewServeMux()

	// Liveness: o processo está vivo. Nunca depende de nada externo.
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok\n"))
	})

	// Readiness: só recebe tráfego depois da primeira rodada de verificação.
	mux.HandleFunc("GET /readyz", func(w http.ResponseWriter, r *http.Request) {
		if !src.Ready() {
			http.Error(w, "ainda verificando", http.StatusServiceUnavailable)
			return
		}
		w.Write([]byte("ready\n"))
	})

	// JSON com os últimos resultados.
	mux.HandleFunc("GET /api/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(src.Results()); err != nil {
			log.Printf("erro escrevendo JSON: %v", err)
		}
	})

	// Página HTML. O "{$}" faz casar só com "/" exato, e não com qualquer caminho.
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := page.Execute(w, src.Results()); err != nil {
			log.Printf("erro gerando HTML: %v", err)
		}
	})

	return logRequests(mux)
}

// logRequests é um "middleware": embrulha o roteador e registra cada pedido.
func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/healthz" && r.URL.Path != "/readyz" { // probes geram ruído demais
			log.Printf("%s %s", r.Method, r.URL.Path)
		}
		next.ServeHTTP(w, r)
	})
}

// html/template escapa o conteúdo automaticamente, evitando injeção de HTML.
var page = template.Must(template.New("page").Parse(`<!doctype html>
<html lang="pt-BR">
<head><title>StatusHub</title></head>
<body style="font-family: sans-serif; max-width: 760px; margin: 2rem auto; padding: 0 1rem">
<h1>StatusHub</h1>
{{if not .}}<p>Primeira verificação em andamento...</p>{{else}}
<table style="width: 100%; border-collapse: collapse" border="1" cellpadding="8">
<tr><th>Alvo</th><th>Situação</th><th>HTTP</th><th>Tempo</th></tr>
{{range .}}<tr>
  <td>{{.Name}}<br><small>{{.URL}}</small></td>
  <td>{{if .Up}}<b style="color: green">no ar</b>{{else}}<b style="color: red">fora do ar</b>{{end}}</td>
  <td>{{if .StatusCode}}{{.StatusCode}}{{else}}<small>{{.Error}}</small>{{end}}</td>
  <td>{{.LatencyMS}} ms</td>
</tr>{{end}}
</table>{{end}}
<p><small>Atualize a página (F5) para ver a verificação mais recente.</small></p>
</body>
</html>
`))
