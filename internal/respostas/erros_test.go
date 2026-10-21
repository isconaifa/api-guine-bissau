package respostas

import "testing"

func TestErros(t *testing.T) {
	casos := []struct {
		nome             string
		erro             Erro
		status           int
		codigo, mensagem string
	}{
		{"região", RegiaoNaoEncontrada("xyz"), 404, CodigoRegiaoNaoEncontrada,
			"Região 'xyz' não encontrada."},
		{"região com acento", RegiaoNaoEncontrada("Bafatá"), 404, CodigoRegiaoNaoEncontrada,
			"Região 'Bafatá' não encontrada. Use o slug em minúsculas e sem acentos (ex.: 'bafata')."},
		{"setor", SetorNaoEncontrado("bafata", "xyz"), 404, CodigoSetorNaoEncontrado,
			"Setor 'xyz' não encontrado na região 'bafata'."},
		{"setor com maiúsculas", SetorNaoEncontrado("bafata", "Bambadinca"), 404, CodigoSetorNaoEncontrado,
			"Setor 'Bambadinca' não encontrado na região 'bafata'. Use o slug em minúsculas e sem acentos (ex.: 'bafata')."},
		{"parâmetro", ParametroInvalido("msg."), 400, CodigoParametroInvalido, "msg."},
		{"rota", RotaNaoEncontrada("/v1/nada"), 404, CodigoRotaNaoEncontrada,
			"Rota '/v1/nada' não encontrada. Consulte a documentação no README."},
		{"interno", ErroInterno(), 500, CodigoErroInterno, "Erro interno do servidor."},
	}
	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			if c.erro.Status != c.status || c.erro.Erro.Codigo != c.codigo || c.erro.Erro.Mensagem != c.mensagem {
				t.Errorf("obtido %+v; esperado status %d, código %s, mensagem %q", c.erro, c.status, c.codigo, c.mensagem)
			}
		})
	}
}
