package pesquisa

import (
	"reflect"
	"testing"
)

func itens(nomes ...string) []Item {
	var r []Item
	for i, n := range nomes {
		r = append(r, Item{NomeNormalizado: Normalizar(n), Nivel: i % 3})
	}
	return r
}

func TestPesquisarOrdenacao(t *testing.T) {
	lista := []Item{
		{NomeNormalizado: "bairro de gabu", Nivel: 2}, // 0: contém
		{NomeNormalizado: "gabu", Nivel: 1},           // 1: exata, setor
		{NomeNormalizado: "gabu", Nivel: 0},           // 2: exata, região
		{NomeNormalizado: "gabu novo", Nivel: 2},      // 3: prefixo
		{NomeNormalizado: "bafata", Nivel: 0},         // 4: sem correspondência
		{NomeNormalizado: "a gabu", Nivel: 0},         // 5: contém, região
	}
	obtido := Pesquisar(lista, "Gabú")
	if esperado := []int{2, 1, 3, 5, 0}; !reflect.DeepEqual(obtido, esperado) {
		t.Errorf("Pesquisar = %v; esperado %v", obtido, esperado)
	}
}

func TestPesquisarIgnoraAcentosEMaiusculas(t *testing.T) {
	lista := itens("Gabú", "Bafatá", "Bolama/Bijagós")
	base := Pesquisar(lista, "gabu")
	for _, q := range []string{"GABÚ", "Gabú", "  gabu  ", "gAbU"} {
		if obtido := Pesquisar(lista, q); !reflect.DeepEqual(obtido, base) {
			t.Errorf("Pesquisar(%q) = %v; esperado %v", q, obtido, base)
		}
	}
	if obtido := Pesquisar(lista, "bolama bijagos"); !reflect.DeepEqual(obtido, []int{2}) {
		t.Errorf("Pesquisar(bolama bijagos) = %v", obtido)
	}
	if obtido := Pesquisar(lista, "bafa"); !reflect.DeepEqual(obtido, []int{1}) {
		t.Errorf("Pesquisar(bafa) = %v", obtido)
	}
}

func TestPesquisarSemResultados(t *testing.T) {
	obtido := Pesquisar(itens("Gabú"), "zzzz")
	if obtido == nil || len(obtido) != 0 {
		t.Errorf("Pesquisar sem correspondência = %#v; esperado slice vazio não-nil", obtido)
	}
	if obtido := Pesquisar(itens("Gabú"), "!!!"); obtido == nil || len(obtido) != 0 {
		t.Errorf("Pesquisar com texto vazio após normalização = %#v", obtido)
	}
}
