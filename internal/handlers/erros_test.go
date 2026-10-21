package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// codigosContrato são os códigos de erro de contracts/openapi.yaml.
var codigosContrato = map[string]bool{
	"REGIAO_NAO_ENCONTRADA": true,
	"SETOR_NAO_ENCONTRADO":  true,
	"PARAMETRO_INVALIDO":    true,
	"ROTA_NAO_ENCONTRADA":   true,
	"ERRO_INTERNO":          true,
}

func TestConformidadeDosErros(t *testing.T) {
	r := routerTeste(t)
	casos := []struct {
		caminho string
		status  int
		codigo  string
	}{
		{"/v1/regioes/xyz", 404, "REGIAO_NAO_ENCONTRADA"},
		{"/v1/regioes/xyz/setores", 404, "REGIAO_NAO_ENCONTRADA"},
		{"/v1/regioes/regiao-a/setores/xyz", 404, "SETOR_NAO_ENCONTRADO"},
		{"/v1/regioes/regiao-a/setores/xyz/seccoes", 404, "SETOR_NAO_ENCONTRADO"},
		{"/v1/pesquisa", 400, "PARAMETRO_INVALIDO"},
		{"/v1/nada", 404, "ROTA_NAO_ENCONTRADA"},
		{"/nada", 404, "ROTA_NAO_ENCONTRADA"},
	}
	for _, c := range casos {
		t.Run(c.caminho, func(t *testing.T) {
			w, corpo := pedir(t, r, http.MethodGet, c.caminho)
			if ct := w.Header().Get("Content-Type"); ct != tipoJSON {
				t.Errorf("Content-Type = %q", ct)
			}
			verificarErro(t, w, corpo, c.status, c.codigo)
			if !codigosContrato[c.codigo] {
				t.Errorf("código %s fora do contrato", c.codigo)
			}
		})
	}
}

func TestSucessoNuncaTemErro(t *testing.T) {
	r := routerTeste(t)
	for _, caminho := range []string{
		"/v1/saude", "/v1/regioes", "/v1/regioes/regiao-a", "/v1/regioes/regiao-a/setores",
		"/v1/regioes/regiao-a/setores/setor-b", "/v1/regioes/regiao-a/setores/regiao-a/seccoes",
		"/v1/pesquisa?q=a", "/v1/sobre",
	} {
		w, corpo := pedir(t, r, http.MethodGet, caminho)
		if _, ok := corpo["erro"]; ok || w.Code != http.StatusOK {
			t.Errorf("%s: status %d, corpo %v", caminho, w.Code, corpo)
		}
		if _, ok := corpo["dados"]; !ok {
			t.Errorf("%s: sem 'dados'", caminho)
		}
	}
}

// Uma barra final redireciona (301) para a rota sem barra; não é uma resposta de erro.
func TestBarraFinalRedireciona(t *testing.T) {
	w := httptest.NewRecorder()
	routerTeste(t).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/v1/regioes/", nil))
	if w.Code != http.StatusMovedPermanently || w.Header().Get("Location") != "/v1/regioes" {
		t.Errorf("status %d, Location %q", w.Code, w.Header().Get("Location"))
	}
}
