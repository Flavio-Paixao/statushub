package config

import (
	"testing"
	"time"
)

// TestParseValido confere um arquivo correto e os valores padrão.
func TestParseValido(t *testing.T) {
	yaml := []byte(`
targets:
  - name: Exemplo
    url: https://example.com
`)
	cfg, err := Parse(yaml)
	if err != nil {
		t.Fatalf("não esperava erro, veio: %v", err)
	}
	if len(cfg.Targets) != 1 {
		t.Fatalf("esperava 1 target, veio %d", len(cfg.Targets))
	}
	if cfg.Interval != 30*time.Second {
		t.Errorf("interval padrão: esperava 30s, veio %v", cfg.Interval)
	}
	if cfg.Timeout != 5*time.Second {
		t.Errorf("timeout padrão: esperava 5s, veio %v", cfg.Timeout)
	}
}

// TestParseInvalido usa "tabela de casos": cada linha é um arquivo errado
// e todos precisam dar erro. É o estilo de teste mais comum em Go.
func TestParseInvalido(t *testing.T) {
	casos := []struct {
		nome string
		yaml string
	}{
		{"lista vazia", "targets: []"},
		{"sem name", "targets:\n  - url: https://example.com"},
		{"url sem http", "targets:\n  - name: X\n    url: example.com"},
		{"timeout maior que interval", "interval: 5s\ntimeout: 10s\ntargets:\n  - name: X\n    url: https://example.com"},
		{"yaml quebrado", "targets: [nome"},
	}
	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			if _, err := Parse([]byte(c.yaml)); err == nil {
				t.Errorf("esperava erro para %q, mas não veio", c.nome)
			}
		})
	}
}

// TestLoadArquivoInexistente confere a mensagem quando o arquivo não existe.
func TestLoadArquivoInexistente(t *testing.T) {
	if _, err := Load("nao-existe.yaml"); err == nil {
		t.Error("esperava erro para arquivo inexistente")
	}
}
