package checker

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Flavio-Paixao/statushub/internal/config"
)

// httptest.NewServer sobe um servidor de mentira só para o teste,
// assim não dependemos da internet nem de sites de verdade.
func TestCheckAll(t *testing.T) {
	ok := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ok.Close()

	erro := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer erro.Close()

	lento := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(300 * time.Millisecond) // mais que o timeout abaixo
	}))
	defer lento.Close()

	c := New(config.Config{
		Interval: time.Second,
		Timeout:  100 * time.Millisecond,
		Targets: []config.Target{
			{Name: "ok", URL: ok.URL},
			{Name: "erro", URL: erro.URL},
			{Name: "lento", URL: lento.URL},
		},
	})

	if c.Ready() {
		t.Fatal("não deveria estar pronto antes da primeira rodada")
	}

	c.checkAll(context.Background())

	if !c.Ready() {
		t.Fatal("deveria estar pronto depois da primeira rodada")
	}

	res := c.Results()
	if !res[0].Up || res[0].StatusCode != 200 {
		t.Errorf("ok: esperava no ar com 200, veio %+v", res[0])
	}
	if res[1].Up || res[1].StatusCode != 500 {
		t.Errorf("erro: esperava fora do ar com 500, veio %+v", res[1])
	}
	if res[2].Up || res[2].Error == "" {
		t.Errorf("lento: esperava fora do ar por timeout, veio %+v", res[2])
	}
}
