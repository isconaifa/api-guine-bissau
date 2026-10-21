package respostas

import "testing"

func TestDetalheSetor(t *testing.T) {
	base := carregarValido(t)
	s, r, _ := base.Setor("regiao-a", "regiao-a")
	obtido := serializar(t, DetalheSetor(base, s, r))
	esperado := `{"dados":{"slug":"regiao-a","nome":"Região A","fonte_id":"fonte-teste","verificado":true,` +
		`"regiao":{"slug":"regiao-a","nome":"Região A"},"total_seccoes":2},"fonte":"Fonte de Teste"}` + "\n"
	if obtido != esperado {
		t.Errorf("obtido  %s\nesperado %s", obtido, esperado)
	}
}

func TestSeccoesDoSetor(t *testing.T) {
	base := carregarValido(t)
	s, _, _ := base.Setor("regiao-a", "regiao-a")
	obtido := serializar(t, SeccoesDoSetor(base, s))
	esperado := `{"dados":[` +
		`{"slug":"seccao-um","nome":"Secção Um","fonte_id":"fonte-teste","verificado":false},` +
		`{"slug":"seccao-dois","nome":"Secção Dois","fonte_id":"fonte-teste","verificado":false}` +
		`],"total":2,"fonte":"Fonte de Teste"}` + "\n"
	if obtido != esperado {
		t.Errorf("obtido  %s\nesperado %s", obtido, esperado)
	}

	s, _, _ = base.Setor("regiao-a", "setor-b")
	obtido = serializar(t, SeccoesDoSetor(base, s))
	// Lista vazia: a fonte é a do próprio setor.
	if esperado := `{"dados":[],"total":0,"fonte":"Fonte Secundária & Teste"}` + "\n"; obtido != esperado {
		t.Errorf("setor sem secções: obtido %s; esperado %s", obtido, esperado)
	}
}
