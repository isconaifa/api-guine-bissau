package main

import (
	"bytes"
	"encoding/json"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/isconaifa/api-guine-bissau/internal/dados"
	"github.com/isconaifa/api-guine-bissau/internal/handlers"
)

func baseTeste(t *testing.T) *dados.Base {
	t.Helper()
	b, err := os.ReadFile("../../internal/dados/testdata/valido.json")
	if err != nil {
		t.Fatal(err)
	}
	base, err := dados.Carregar(b)
	if err != nil {
		t.Fatal(err)
	}
	return base
}

// TestParidadeComAPI garante que cada ficheiro estático é byte a byte igual à resposta 200
// da rota equivalente da API (contracts/json-estatico.md).
func TestParidadeComAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)
	base := baseTeste(t)
	dir := t.TempDir()
	if err := Gerar(base, dir); err != nil {
		t.Fatal(err)
	}

	esperados := []string{
		"v1/regioes.json",
		"v1/regioes/regiao-a.json",
		"v1/regioes/regiao-a/setores.json",
		"v1/regioes/regiao-a/setores/setor-b.json",
		"v1/regioes/regiao-a/setores/regiao-a/seccoes.json",
		"v1/regioes/setor-autonomo-teste.json",
		"v1/regioes/setor-autonomo-teste/setores.json",
		"v1/sobre.json",
		"v1/pesquisa/indice.json",
	}
	for _, f := range esperados {
		if _, err := os.Stat(filepath.Join(dir, f)); err != nil {
			t.Errorf("ficheiro em falta: %s", f)
		}
	}

	router := handlers.NovoRouter(base)
	comparados := 0
	err := filepath.WalkDir(dir, func(caminho string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, _ := filepath.Rel(dir, caminho)
		rel = filepath.ToSlash(rel)
		if rel == "v1/pesquisa/indice.json" {
			return nil // sem rota equivalente
		}
		rota := "/" + strings.TrimSuffix(rel, ".json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, httptest.NewRequest(http.MethodGet, rota, nil))
		conteudo, _ := os.ReadFile(caminho)
		if w.Code != http.StatusOK || !bytes.Equal(w.Body.Bytes(), conteudo) {
			t.Errorf("%s difere de GET %s (status %d)\nficheiro: %s\napi:      %s", rel, rota, w.Code, conteudo, w.Body.Bytes())
		}
		comparados++
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	// 1 lista + 2×(detalhe + setores) + 2×(detalhe de setor + secções) + sobre
	if comparados != 10 {
		t.Errorf("ficheiros comparados = %d; esperado 10", comparados)
	}
}

func TestIndicePesquisa(t *testing.T) {
	dir := t.TempDir()
	if err := Gerar(baseTeste(t), dir); err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(filepath.Join(dir, "v1/pesquisa/indice.json"))
	var indice struct {
		Dados []map[string]any `json:"dados"`
		Total int              `json:"total"`
	}
	if err := json.Unmarshal(b, &indice); err != nil {
		t.Fatal(err)
	}
	if indice.Total != 6 || len(indice.Dados) != 6 {
		t.Fatalf("índice com %d itens; esperado 6", indice.Total)
	}
	for _, it := range indice.Dados {
		if it["nome_normalizado"] == "" || it["caminho"] == "" {
			t.Errorf("item incompleto: %v", it)
		}
	}
	if indice.Dados[0]["nome_normalizado"] != "regiao a" {
		t.Errorf("primeiro item: %v", indice.Dados[0])
	}
}
