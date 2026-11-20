package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Flavio-Paixao/statushub/internal/checker"
)

// fakeSource finge ser o verificador, com dados que o teste controla.
type fakeSource struct {
	results []checker.Result
	ready   bool
}

func (f fakeSource) Results() []checker.Result { return f.results }
func (f fakeSource) Ready() bool               { return f.ready }

// get faz um pedido de mentira e devolve a resposta gravada.
func get(h http.Handler, path string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
	return rec
}

func TestReadyz(t *testing.T) {
	if code := get(New(fakeSource{ready: false}), "/readyz").Code; code != http.StatusServiceUnavailable {
		t.Errorf("antes da 1ª rodada: esperava 503, veio %d", code)
	}
	if code := get(New(fakeSource{ready: true}), "/readyz").Code; code != http.StatusOK {
		t.Errorf("depois da 1ª rodada: esperava 200, veio %d", code)
	}
}

func TestHealthz(t *testing.T) {
	if code := get(New(fakeSource{}), "/healthz").Code; code != http.StatusOK {
		t.Errorf("esperava 200, veio %d", code)
	}
}

func TestAPIStatus(t *testing.T) {
	src := fakeSource{ready: true, results: []checker.Result{
		{Name: "Exemplo", URL: "https://example.com", Up: true, StatusCode: 200},
	}}
	rec := get(New(src), "/api/status")

	var got []checker.Result
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("JSON inválido: %v", err)
	}
	if len(got) != 1 || got[0].Name != "Exemplo" || !got[0].Up {
		t.Errorf("JSON inesperado: %+v", got)
	}
}

func TestPaginaHTML(t *testing.T) {
	src := fakeSource{ready: true, results: []checker.Result{{Name: "Exemplo", Up: true, StatusCode: 200}}}
	body := get(New(src), "/").Body.String()
	if !strings.Contains(body, "Exemplo") || !strings.Contains(body, "no ar") {
		t.Errorf("página não mostra o alvo: %s", body)
	}
	if code := get(New(src), "/qualquer-coisa").Code; code != http.StatusNotFound {
		t.Errorf("caminho desconhecido: esperava 404, veio %d", code)
	}
}
