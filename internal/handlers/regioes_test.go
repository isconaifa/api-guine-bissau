package handlers

import (
	"net/http"
	"strings"
	"testing"
)

func TestRegioes(t *testing.T) {
	r := routerTeste(t)

	w, corpo := pedir(t, r, http.MethodGet, "/v1/regioes")
	if w.Code != http.StatusOK || corpo["total"] != 2.0 {
		t.Errorf("/v1/regioes: status %d, corpo %v", w.Code, corpo)
	}

	w, corpo = pedir(t, r, http.MethodGet, "/v1/regioes/regiao-a")
	dados, _ := corpo["dados"].(map[string]any)
	if w.Code != http.StatusOK || dados["nome"] != "Região A" || len(dados["setores"].([]any)) != 2 {
		t.Errorf("/v1/regioes/regiao-a: status %d, corpo %v", w.Code, corpo)
	}

	w, corpo = pedir(t, r, http.MethodGet, "/v1/regioes/regiao-a/setores")
	if w.Code != http.StatusOK || corpo["total"] != 2.0 {
		t.Errorf("/v1/regioes/regiao-a/setores: status %d, corpo %v", w.Code, corpo)
	}
}

func TestRegiaoNaoEncontrada(t *testing.T) {
	r := routerTeste(t)
	for _, caminho := range []string{"/v1/regioes/xyz", "/v1/regioes/xyz/setores"} {
		w, corpo := pedir(t, r, http.MethodGet, caminho)
		msg := verificarErro(t, w, corpo, http.StatusNotFound, "REGIAO_NAO_ENCONTRADA")
		if msg != "Região 'xyz' não encontrada." {
			t.Errorf("%s: mensagem %q", caminho, msg)
		}
	}

	w, corpo := pedir(t, r, http.MethodGet, "/v1/regioes/Bafat%C3%A1")
	msg := verificarErro(t, w, corpo, http.StatusNotFound, "REGIAO_NAO_ENCONTRADA")
	if !strings.Contains(msg, "'Bafatá'") || !strings.Contains(msg, "minúsculas e sem acentos") {
		t.Errorf("mensagem sem dica de formato: %q", msg)
	}
}
