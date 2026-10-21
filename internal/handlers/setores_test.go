package handlers

import (
	"net/http"
	"testing"
)

func TestSetores(t *testing.T) {
	r := routerTeste(t)

	w, corpo := pedir(t, r, http.MethodGet, "/v1/regioes/regiao-a/setores/setor-b")
	dados, _ := corpo["dados"].(map[string]any)
	regiao, _ := dados["regiao"].(map[string]any)
	if w.Code != http.StatusOK || dados["nome"] != "Setor B" || regiao["slug"] != "regiao-a" {
		t.Errorf("detalhe do setor: status %d, corpo %v", w.Code, corpo)
	}

	w, corpo = pedir(t, r, http.MethodGet, "/v1/regioes/regiao-a/setores/regiao-a/seccoes")
	if w.Code != http.StatusOK || corpo["total"] != 2.0 {
		t.Errorf("secções: status %d, corpo %v", w.Code, corpo)
	}

	w, corpo = pedir(t, r, http.MethodGet, "/v1/regioes/regiao-a/setores/setor-b/seccoes")
	if lista, _ := corpo["dados"].([]any); w.Code != http.StatusOK || lista == nil || len(lista) != 0 {
		t.Errorf("setor sem secções: status %d, corpo %v", w.Code, corpo)
	}
}

func TestSetorNaoEncontrado(t *testing.T) {
	r := routerTeste(t)
	casos := []struct{ caminho, codigo, mensagem string }{
		{"/v1/regioes/regiao-a/setores/xyz", "SETOR_NAO_ENCONTRADO", "Setor 'xyz' não encontrado na região 'regiao-a'."},
		{"/v1/regioes/regiao-a/setores/xyz/seccoes", "SETOR_NAO_ENCONTRADO", "Setor 'xyz' não encontrado na região 'regiao-a'."},
		// setor existente pedido sob outra região
		{"/v1/regioes/setor-autonomo-teste/setores/setor-b", "SETOR_NAO_ENCONTRADO", "Setor 'setor-b' não encontrado na região 'setor-autonomo-teste'."},
		{"/v1/regioes/xyz/setores/setor-b", "REGIAO_NAO_ENCONTRADA", "Região 'xyz' não encontrada."},
		{"/v1/regioes/xyz/setores/setor-b/seccoes", "REGIAO_NAO_ENCONTRADA", "Região 'xyz' não encontrada."},
	}
	for _, c := range casos {
		t.Run(c.caminho, func(t *testing.T) {
			w, corpo := pedir(t, r, http.MethodGet, c.caminho)
			if msg := verificarErro(t, w, corpo, http.StatusNotFound, c.codigo); msg != c.mensagem {
				t.Errorf("mensagem %q; esperado %q", msg, c.mensagem)
			}
		})
	}
}
