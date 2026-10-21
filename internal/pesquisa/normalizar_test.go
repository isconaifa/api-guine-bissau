package pesquisa

import "testing"

func TestNormalizar(t *testing.T) {
	casos := []struct {
		entrada, esperado string
	}{
		{"Bafatá", "bafata"},
		{"GABÚ", "gabu"},
		{"Bolama/Bijagós", "bolama bijagos"},
		{"  Setor   Autónomo de Bissau ", "setor autonomo de bissau"},
		{"Quinhámel", "quinhamel"},
		{"ção", "cao"},
		{"sector-autonomo-de-bissau", "sector autonomo de bissau"},
		{"", ""},
		{"!!!", ""},
	}
	for _, c := range casos {
		if obtido := Normalizar(c.entrada); obtido != c.esperado {
			t.Errorf("Normalizar(%q) = %q; esperado %q", c.entrada, obtido, c.esperado)
		}
	}
}
