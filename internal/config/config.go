// Package config lê e valida o arquivo de alvos do StatusHub.
package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"go.yaml.in/yaml/v3"
)

// Valores usados quando o arquivo não informa interval ou timeout.
const (
	DefaultInterval = 30 * time.Second
	DefaultTimeout  = 5 * time.Second
)

// Target é um site ou API que vamos verificar.
// As "tags" `yaml:"..."` dizem qual chave do arquivo vai para cada campo.
type Target struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

// Config é o arquivo inteiro.
type Config struct {
	Interval time.Duration `yaml:"interval"` // de quanto em quanto tempo verificar (ex.: 30s)
	Timeout  time.Duration `yaml:"timeout"`  // quanto esperar cada resposta (ex.: 5s)
	Targets  []Target      `yaml:"targets"`  // a lista de alvos
}

// Load lê o arquivo em path, aplica os padrões e valida.
// Em Go, funções que podem falhar devolvem (resultado, erro).
func Load(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("lendo %s: %w", path, err)
	}
	return Parse(data)
}

// Parse faz o trabalho de verdade a partir dos bytes do arquivo.
// Separar Load de Parse deixa o teste fácil: dá para testar sem criar arquivos.
func Parse(data []byte) (Config, error) {
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("YAML inválido: %w", err)
	}

	if cfg.Interval == 0 {
		cfg.Interval = DefaultInterval
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = DefaultTimeout
	}

	if err := cfg.validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

// validate confere as regras mínimas antes de o serviço começar.
func (c Config) validate() error {
	if c.Interval < 0 || c.Timeout < 0 {
		return errors.New("interval e timeout não podem ser negativos")
	}
	if c.Timeout >= c.Interval {
		return errors.New("timeout precisa ser menor que interval")
	}
	if len(c.Targets) == 0 {
		return errors.New("a lista de targets está vazia")
	}
	for i, t := range c.Targets {
		if t.Name == "" {
			return fmt.Errorf("target %d: name vazio", i+1)
		}
		if !strings.HasPrefix(t.URL, "http://") && !strings.HasPrefix(t.URL, "https://") {
			return fmt.Errorf("target %q: url precisa começar com http:// ou https://", t.Name)
		}
	}
	return nil
}
