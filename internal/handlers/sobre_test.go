package handlers

import (
	"net/http"
	"testing"
)

func TestSobre(t *testing.T) {
	w, corpo := pedir(t, routerTeste(t), http.MethodGet, "/v1/sobre")
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
	if _, ok := corpo["total"]; ok {
		t.Error("/v1/sobre não deve ter 'total'")
	}
	if _, ok := corpo["fonte"]; ok {
		t.Error("/v1/sobre não deve ter 'fonte' no envelope")
	}
	d := corpo["dados"].(map[string]any)
	esperados := map[string]any{
		"nome":        "API Guiné-Bissau",
		"autor":       "Francisco Naifa",
		"licenca":     "MIT",
		"versao":      "dev",
		"versao_api":  "v1",
		"repositorio": "https://github.com/isconaifa/api-guine-bissau",
	}
	for campo, valor := range esperados {
		if d[campo] != valor {
			t.Errorf("%s = %v; esperado %v", campo, d[campo], valor)
		}
	}
	fontes, _ := d["fontes"].([]any)
	if len(fontes) != 2 {
		t.Fatalf("fontes = %v", d["fontes"])
	}
	primeira, segunda := fontes[0].(map[string]any), fontes[1].(map[string]any)
	for _, campo := range []string{"id", "nome", "referencia", "consultado_em", "url"} {
		if primeira[campo] == nil {
			t.Errorf("fonte sem %q: %v", campo, primeira)
		}
	}
	if _, ok := segunda["url"]; ok {
		t.Errorf("url ausente não deve ser serializada: %v", segunda)
	}
}
