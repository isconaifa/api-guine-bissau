package respostas

import "testing"

func TestListaRegioes(t *testing.T) {
	base := carregarValido(t)
	obtido := serializar(t, ListaRegioes(base))
	esperado := `{"dados":[` +
		`{"slug":"regiao-a","nome":"Região A","tipo":"regiao","capital":"Cidade A","capital_fonte_id":"fonte-secundaria","fonte_id":"fonte-teste","verificado":true},` +
		`{"slug":"setor-autonomo-teste","nome":"Setor Autónomo Teste","tipo":"setor-autonomo","capital":"Cidade Autónoma","fonte_id":"fonte-teste","verificado":false}` +
		`],"total":2,"fonte":"Fonte de Teste; Fonte Secundária & Teste"}` + "\n"
	if obtido != esperado {
		t.Errorf("obtido  %s\nesperado %s", obtido, esperado)
	}
}

func TestDetalheRegiao(t *testing.T) {
	base := carregarValido(t)
	r, _ := base.Regiao("regiao-a")
	obtido := serializar(t, DetalheRegiao(base, r))
	esperado := `{"dados":{"slug":"regiao-a","nome":"Região A","tipo":"regiao","capital":"Cidade A","capital_fonte_id":"fonte-secundaria","fonte_id":"fonte-teste","verificado":true,` +
		`"setores":[{"slug":"regiao-a","nome":"Região A"},{"slug":"setor-b","nome":"Setor B"}]},` +
		`"fonte":"Fonte de Teste; Fonte Secundária & Teste"}` + "\n"
	if obtido != esperado {
		t.Errorf("obtido  %s\nesperado %s", obtido, esperado)
	}

	r, _ = base.Regiao("setor-autonomo-teste")
	obtido = serializar(t, DetalheRegiao(base, r))
	esperado = `{"dados":{"slug":"setor-autonomo-teste","nome":"Setor Autónomo Teste","tipo":"setor-autonomo","capital":"Cidade Autónoma","fonte_id":"fonte-teste","verificado":false,"setores":[]},"fonte":"Fonte de Teste"}` + "\n"
	if obtido != esperado {
		t.Errorf("sem setores: obtido  %s\nesperado %s", obtido, esperado)
	}
}

func TestSetoresDaRegiao(t *testing.T) {
	base := carregarValido(t)
	r, _ := base.Regiao("regiao-a")
	obtido := serializar(t, SetoresDaRegiao(base, r))
	esperado := `{"dados":[` +
		`{"slug":"regiao-a","nome":"Região A","fonte_id":"fonte-teste","verificado":true},` +
		`{"slug":"setor-b","nome":"Setor B","fonte_id":"fonte-secundaria","verificado":false}` +
		`],"total":2,"fonte":"Fonte de Teste; Fonte Secundária & Teste"}` + "\n"
	if obtido != esperado {
		t.Errorf("obtido  %s\nesperado %s", obtido, esperado)
	}
}
