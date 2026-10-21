package handlers

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSaude(t *testing.T) {
	w, _ := pedir(t, routerTeste(t), http.MethodGet, "/v1/saude")
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
	if w.Body.String() != `{"dados":{"estado":"ok"}}`+"\n" {
		t.Errorf("corpo = %q", w.Body.String())
	}
	if cc := w.Header().Get("Cache-Control"); cc != "no-store" {
		t.Errorf("Cache-Control = %q; esperado no-store", cc)
	}
}

func TestRotasInexistentes(t *testing.T) {
	r := routerTeste(t)
	for _, c := range []struct{ metodo, caminho string }{
		{http.MethodGet, "/saude"},
		{http.MethodGet, "/v1/inexistente"},
		{http.MethodPost, "/v1/regioes"},
	} {
		t.Run(c.metodo+" "+c.caminho, func(t *testing.T) {
			w, corpo := pedir(t, r, c.metodo, c.caminho)
			verificarErro(t, w, corpo, http.StatusNotFound, "ROTA_NAO_ENCONTRADA")
		})
	}
}

func TestCabecalhos(t *testing.T) {
	r := routerTeste(t)
	r.Group("/v1", cabecalhosDados).GET("/teste-dados", func(c *gin.Context) {
		responder(c, http.StatusOK, map[string]any{"dados": []int{}})
	})
	w, _ := pedir(t, r, http.MethodGet, "/v1/teste-dados")
	esperados := map[string]string{
		"Access-Control-Allow-Origin": "*",
		"Cache-Control":               "public, max-age=3600",
		"Content-Type":                "application/json; charset=utf-8",
	}
	for nome, valor := range esperados {
		if obtido := w.Header().Get(nome); obtido != valor {
			t.Errorf("%s = %q; esperado %q", nome, obtido, valor)
		}
	}
	// Erros também são JSON e têm CORS, para clientes no browser lerem a mensagem.
	w, _ = pedir(t, routerTeste(t), http.MethodGet, "/v1/nada")
	if w.Header().Get("Access-Control-Allow-Origin") != "*" || w.Header().Get("Content-Type") != "application/json; charset=utf-8" {
		t.Errorf("cabeçalhos do erro: %v", w.Header())
	}
}

func TestPanicoDevolveErroInterno(t *testing.T) {
	r := routerTeste(t)
	r.GET("/v1/panico", func(*gin.Context) { panic("falha de teste") })
	w, corpo := pedir(t, r, http.MethodGet, "/v1/panico")
	verificarErro(t, w, corpo, http.StatusInternalServerError, "ERRO_INTERNO")
}
