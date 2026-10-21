package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	brutos "github.com/isconaifa/api-guine-bissau/dados"
	"github.com/isconaifa/api-guine-bissau/internal/dados"
)

func benchmarkRota(b *testing.B, caminho string) {
	base, err := dados.Carregar(brutos.DivisaoAdministrativa)
	if err != nil {
		b.Fatal(err)
	}
	r := NovoRouter(base)
	req := httptest.NewRequest(http.MethodGet, caminho, nil)
	b.ReportAllocs()
	for b.Loop() {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			b.Fatalf("status %d", w.Code)
		}
	}
}

func BenchmarkListaRegioes(b *testing.B) { benchmarkRota(b, "/v1/regioes") }

func BenchmarkPesquisa(b *testing.B) { benchmarkRota(b, "/v1/pesquisa?q=ba") }
