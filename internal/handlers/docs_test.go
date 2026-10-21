package handlers

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/isconaifa/api-guine-bissau/docs"
)

func TestDocumentacao(t *testing.T) {
	r := routerTeste(t)
	for _, c := range []struct{ caminho, tipo, contem string }{
		{"/v1/docs/", "text/html; charset=utf-8", "swagger-ui"},
		{"/v1/docs/openapi.yaml", "application/yaml; charset=utf-8", "openapi: 3.1.0"},
	} {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, c.caminho, nil))
		if w.Code != http.StatusOK || w.Header().Get("Content-Type") != c.tipo || !strings.Contains(w.Body.String(), c.contem) {
			t.Errorf("%s: status %d, tipo %q", c.caminho, w.Code, w.Header().Get("Content-Type"))
		}
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/v1/docs", nil))
	if w.Code != http.StatusMovedPermanently || w.Header().Get("Location") != "/v1/docs/" {
		t.Errorf("/v1/docs: status %d, Location %q", w.Code, w.Header().Get("Location"))
	}
}

// TestContratoCobreTodasAsRotas garante que docs/openapi.yaml documenta exatamente as rotas de
// dados do router (Princípio X).
func TestContratoCobreTodasAsRotas(t *testing.T) {
	spec := string(docs.OpenAPI)
	documentadas := map[string]bool{}
	for _, m := range regexp.MustCompile(`(?m)^  (/\S+):$`).FindAllStringSubmatch(spec, -1) {
		documentadas[m[1]] = true
	}
	registadas := map[string]bool{}
	param := regexp.MustCompile(`:(\w+)`)
	for _, rota := range routerTeste(t).Routes() {
		if strings.HasPrefix(rota.Path, "/v1/docs") {
			continue // a própria documentação
		}
		caminho := param.ReplaceAllString(rota.Path, "{$1}")
		registadas[caminho] = true
		if !documentadas[caminho] {
			t.Errorf("rota %s %s não documentada em docs/openapi.yaml", rota.Method, caminho)
		}
	}
	for caminho := range documentadas {
		if !registadas[caminho] {
			t.Errorf("docs/openapi.yaml documenta %s, que não existe no router", caminho)
		}
	}
}

// Pedidos de páginas HTTPS públicas (ex.: Swagger no GitHub Pages) para o servidor local exigem
// uma resposta à pré-verificação CORS / Private Network Access.
func TestPreVerificacaoCORS(t *testing.T) {
	req := httptest.NewRequest(http.MethodOptions, "/v1/regioes/oio", nil)
	req.Header.Set("Origin", "https://isconaifa.github.io")
	req.Header.Set("Access-Control-Request-Method", "GET")
	req.Header.Set("Access-Control-Request-Headers", "accept")
	req.Header.Set("Access-Control-Request-Private-Network", "true")
	w := httptest.NewRecorder()
	routerTeste(t).ServeHTTP(w, req)
	esperados := map[string]string{
		"Access-Control-Allow-Origin":          "*",
		"Access-Control-Allow-Methods":         "GET, OPTIONS",
		"Access-Control-Allow-Headers":         "accept",
		"Access-Control-Allow-Private-Network": "true",
	}
	if w.Code != http.StatusNoContent {
		t.Errorf("status = %d; esperado 204", w.Code)
	}
	for nome, valor := range esperados {
		if obtido := w.Header().Get(nome); obtido != valor {
			t.Errorf("%s = %q; esperado %q", nome, obtido, valor)
		}
	}
}
