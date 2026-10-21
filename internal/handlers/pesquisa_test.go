package handlers

import (
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

func TestPesquisa(t *testing.T) {
	r := routerTeste(t)

	_, corpo1 := pedir(t, r, http.MethodGet, "/v1/pesquisa?q=regiao-a")
	_, corpo2 := pedir(t, r, http.MethodGet, "/v1/pesquisa?q="+url.QueryEscape("REGIÃO A"))
	if !reflect.DeepEqual(corpo1, corpo2) {
		t.Errorf("resultados diferentes:\n%v\n%v", corpo1, corpo2)
	}

	w, corpo := pedir(t, r, http.MethodGet, "/v1/pesquisa?q="+url.QueryEscape("secção um"))
	itens, _ := corpo["dados"].([]any)
	if w.Code != http.StatusOK || len(itens) != 1 {
		t.Fatalf("pesquisa por secção: status %d, corpo %v", w.Code, corpo)
	}
	seccao := itens[0].(map[string]any)
	if seccao["tipo"] != "seccao" || seccao["caminho"] != "/v1/regioes/regiao-a/setores/regiao-a/seccoes" ||
		seccao["regiao"] == nil || seccao["setor"] == nil || seccao["fonte_id"] != "fonte-teste" {
		t.Errorf("secção: %v", seccao)
	}

	_, corpo = pedir(t, r, http.MethodGet, "/v1/pesquisa?q=regiao")
	porTipo := map[string]map[string]any{}
	for _, it := range corpo["dados"].([]any) {
		m := it.(map[string]any)
		porTipo[m["tipo"].(string)] = m
	}
	if reg := porTipo["regiao"]; reg == nil || reg["caminho"] != "/v1/regioes/regiao-a" || reg["regiao"] != nil || reg["setor"] != nil {
		t.Errorf("região: %v", reg)
	}
	if set := porTipo["setor"]; set == nil || set["caminho"] != "/v1/regioes/regiao-a/setores/regiao-a" || set["regiao"] == nil || set["setor"] != nil {
		t.Errorf("setor: %v", set)
	}

	w, corpo = pedir(t, r, http.MethodGet, "/v1/pesquisa?q=zzzz")
	if lista, _ := corpo["dados"].([]any); w.Code != http.StatusOK || lista == nil || len(lista) != 0 || corpo["total"] != 0.0 {
		t.Errorf("sem correspondência: status %d, corpo %v", w.Code, corpo)
	}
}

func TestPesquisaParametroInvalido(t *testing.T) {
	r := routerTeste(t)
	for _, caminho := range []string{
		"/v1/pesquisa",
		"/v1/pesquisa?q=",
		"/v1/pesquisa?q=%20%20",
		"/v1/pesquisa?q=" + strings.Repeat("a", 101),
	} {
		w, corpo := pedir(t, r, http.MethodGet, caminho)
		msg := verificarErro(t, w, corpo, http.StatusBadRequest, "PARAMETRO_INVALIDO")
		if msg != "O parâmetro 'q' é obrigatório e deve ter entre 1 e 100 caracteres." {
			t.Errorf("%s: mensagem %q", caminho, msg)
		}
	}
	if w, _ := pedir(t, r, http.MethodGet, "/v1/pesquisa?q="+strings.Repeat("a", 100)); w.Code != http.StatusOK {
		t.Errorf("q com 100 caracteres deveria ser aceite: %d", w.Code)
	}
}
