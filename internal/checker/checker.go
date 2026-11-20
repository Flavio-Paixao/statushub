// Package checker verifica os alvos de tempos em tempos e guarda o último resultado.
package checker

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/Flavio-Paixao/statushub/internal/config"
)

// Result é o resultado da última verificação de um alvo.
// As tags `json:"..."` definem os nomes no /api/status.
type Result struct {
	Name       string    `json:"name"`
	URL        string    `json:"url"`
	Up         bool      `json:"up"`
	StatusCode int       `json:"status_code"`
	LatencyMS  int64     `json:"latency_ms"`
	Error      string    `json:"error,omitempty"`
	CheckedAt  time.Time `json:"checked_at"`
}

// Checker guarda a configuração e os resultados.
// Os resultados são escritos por uma goroutine e lidos pelos handlers HTTP
// ao mesmo tempo, por isso ficam protegidos por um mutex.
type Checker struct {
	targets  []config.Target
	interval time.Duration
	client   *http.Client

	mu      sync.RWMutex // protege results e ready
	results []Result
	ready   bool // vira true depois da primeira rodada completa
}

// New cria um Checker a partir da configuração.
func New(cfg config.Config) *Checker {
	return &Checker{
		targets:  cfg.Targets,
		interval: cfg.Interval,
		// O Timeout do cliente garante que um site travado não prende a verificação.
		client: &http.Client{Timeout: cfg.Timeout},
	}
}

// Run verifica todos os alvos agora e depois a cada intervalo,
// até o contexto ser cancelado (quando o serviço está desligando).
func (c *Checker) Run(ctx context.Context) {
	c.checkAll(ctx)

	ticker := time.NewTicker(c.interval) // um "despertador" que toca a cada intervalo
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.checkAll(ctx)
		}
	}
}

// checkAll verifica todos os alvos em paralelo e troca os resultados de uma vez.
func (c *Checker) checkAll(ctx context.Context) {
	results := make([]Result, len(c.targets))

	var wg sync.WaitGroup // espera todas as goroutines terminarem
	for i, t := range c.targets {
		wg.Add(1)
		go func() { // cada alvo é verificado numa goroutine (uma "linha de execução" leve)
			defer wg.Done()
			results[i] = c.check(ctx, t)
		}()
	}
	wg.Wait()

	c.mu.Lock() // trava para escrita: ninguém lê enquanto trocamos
	c.results = results
	c.ready = true
	c.mu.Unlock()
}

// check faz um GET no alvo e mede o tempo de resposta.
func (c *Checker) check(ctx context.Context, t config.Target) Result {
	r := Result{Name: t.Name, URL: t.URL, CheckedAt: time.Now()}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, t.URL, nil)
	if err != nil {
		r.Error = err.Error()
		return r
	}
	req.Header.Set("User-Agent", "StatusHub/1.0")

	start := time.Now()
	resp, err := c.client.Do(req)
	r.LatencyMS = time.Since(start).Milliseconds()
	if err != nil {
		r.Error = err.Error() // DNS, conexão recusada, timeout...
		return r
	}
	defer resp.Body.Close()

	r.StatusCode = resp.StatusCode
	r.Up = resp.StatusCode >= 200 && resp.StatusCode < 400
	return r
}

// Results devolve uma cópia dos últimos resultados.
// A cópia evita que quem chamou altere os dados guardados.
func (c *Checker) Results() []Result {
	c.mu.RLock() // trava só para leitura: vários leitores ao mesmo tempo podem
	defer c.mu.RUnlock()
	out := make([]Result, len(c.results))
	copy(out, c.results)
	return out
}

// Ready diz se a primeira rodada de verificação já terminou.
func (c *Checker) Ready() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.ready
}
