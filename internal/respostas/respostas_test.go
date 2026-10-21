package respostas

import (
	"os"
	"testing"

	"github.com/isconaifa/api-guine-bissau/internal/dados"
)

func carregarValido(t *testing.T) *dados.Base {
	t.Helper()
	b, err := os.ReadFile("../dados/testdata/valido.json")
	if err != nil {
		t.Fatal(err)
	}
	base, err := dados.Carregar(b)
	if err != nil {
		t.Fatal(err)
	}
	return base
}

func serializar(t *testing.T, v any) string {
	t.Helper()
	b, err := Serializar(v)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestSerializarEnvelopes(t *testing.T) {
	casos := []struct {
		nome     string
		valor    any
		esperado string
	}{
		{"lista", NovaLista([]string{"a", "b"}, "INE"), `{"dados":["a","b"],"total":2,"fonte":"INE"}` + "\n"},
		{"lista vazia", NovaLista([]string{}, "INE"), `{"dados":[],"total":0,"fonte":"INE"}` + "\n"},
		{"lista nil", NovaLista[string](nil, "INE"), `{"dados":[],"total":0,"fonte":"INE"}` + "\n"},
		{"objeto", Objeto{Dados: map[string]int{"a": 1}, Fonte: "INE"}, `{"dados":{"a":1},"fonte":"INE"}` + "\n"},
		{"objeto sem fonte", Objeto{Dados: map[string]int{"a": 1}}, `{"dados":{"a":1}}` + "\n"},
		{"erro", Erro{Status: 404, Erro: DetalheErro{Codigo: "X", Mensagem: "y"}}, `{"erro":{"codigo":"X","mensagem":"y"}}` + "\n"},
		{"sem escapar HTML", Objeto{Dados: "Bafatá & <Gabú>"}, `{"dados":"Bafatá & <Gabú>"}` + "\n"},
	}
	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			if obtido := serializar(t, c.valor); obtido != c.esperado {
				t.Errorf("obtido  %s\nesperado %s", obtido, c.esperado)
			}
		})
	}
}

func TestResumoFontes(t *testing.T) {
	base := carregarValido(t)
	obtido := ResumoFontes(base, "fonte-teste", "fonte-secundaria", "fonte-teste", "", "nao-existe")
	if esperado := "Fonte de Teste; Fonte Secundária & Teste"; obtido != esperado {
		t.Errorf("ResumoFontes = %q; esperado %q", obtido, esperado)
	}
	if obtido := ResumoFontes(base); obtido != "" {
		t.Errorf("ResumoFontes sem ids = %q; esperado vazio", obtido)
	}
}
