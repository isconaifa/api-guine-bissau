package dados

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

// lerValido devolve os bytes da fixture válida.
func lerValido(t *testing.T) []byte {
	t.Helper()
	b, err := os.ReadFile("testdata/valido.json")
	if err != nil {
		t.Fatalf("ler fixture: %v", err)
	}
	return b
}

// variante aplica uma alteração à fixture válida e devolve o JSON resultante.
func variante(t *testing.T, alterar func(a *Arquivo)) []byte {
	t.Helper()
	var a Arquivo
	if err := json.Unmarshal(lerValido(t), &a); err != nil {
		t.Fatalf("decodificar fixture: %v", err)
	}
	alterar(&a)
	b, err := json.Marshal(a)
	if err != nil {
		t.Fatalf("codificar variante: %v", err)
	}
	return b
}

func TestCarregarValido(t *testing.T) {
	if _, err := Carregar(lerValido(t)); err != nil {
		t.Fatalf("fixture válida rejeitada: %v", err)
	}
}

func TestCarregarJSONInvalido(t *testing.T) {
	casos := map[string]string{
		"mal formado":        `{"fontes": [`,
		"campo desconhecido": `{"fontes": [], "regioes": [], "extra": 1}`,
		"dados a mais":       `{"fontes": [], "regioes": []} {}`,
	}
	for nome, entrada := range casos {
		t.Run(nome, func(t *testing.T) {
			if _, err := Carregar([]byte(entrada)); err == nil {
				t.Fatal("esperava erro")
			}
		})
	}
}

func TestCarregarRegras(t *testing.T) {
	casos := []struct {
		nome    string
		alterar func(a *Arquivo)
		contem  []string // trechos que a mensagem de erro deve conter
	}{
		// (2) campos obrigatórios
		{"nome da região vazio", func(a *Arquivo) { a.Regioes[0].Nome = "" },
			[]string{`região "regiao-a"`, "nome"}},
		{"slug da região vazio", func(a *Arquivo) { a.Regioes[0].Slug = "" },
			[]string{"região #1", "slug"}},
		{"capital vazia", func(a *Arquivo) { a.Regioes[0].Capital = "" },
			[]string{`região "regiao-a"`, "capital"}},
		{"fonte_id do setor vazio", func(a *Arquivo) { a.Regioes[0].Setores[1].FonteID = "" },
			[]string{`setor "setor-b"`, "fonte_id"}},
		{"nome da secção vazio", func(a *Arquivo) { a.Regioes[0].Setores[0].Seccoes[0].Nome = "" },
			[]string{`secção "seccao-um"`, "nome"}},
		{"referência da fonte vazia", func(a *Arquivo) { a.Fontes[0].Referencia = "" },
			[]string{`fonte "fonte-teste"`, "referencia"}},
		{"consultado_em vazio", func(a *Arquivo) { a.Fontes[0].ConsultadoEm = "" },
			[]string{`fonte "fonte-teste"`, "consultado_em"}},

		// (3) formato e unicidade de slugs
		{"slug com acento", func(a *Arquivo) { a.Regioes[0].Slug = "Bafatá" },
			[]string{`"Bafatá"`, "slug"}},
		{"slug com hífen final", func(a *Arquivo) { a.Regioes[0].Setores[1].Slug = "bafata-" },
			[]string{`"bafata-"`, "slug"}},
		{"slug com hífen duplo", func(a *Arquivo) { a.Regioes[0].Setores[0].Seccoes[0].Slug = "ba--fata" },
			[]string{`"ba--fata"`, "slug"}},
		{"slug de região duplicado", func(a *Arquivo) { a.Regioes[1].Slug = "regiao-a" },
			[]string{`região "regiao-a"`, "duplicado"}},
		{"slug de setor duplicado em regiões diferentes", func(a *Arquivo) {
			a.Regioes[1].Setores = []Setor{{Slug: "setor-b", Nome: "Outro Setor B", FonteID: "fonte-teste"}}
		}, []string{`setor "setor-b"`, "duplicado"}},
		{"slug de secção duplicado em setores diferentes", func(a *Arquivo) {
			a.Regioes[0].Setores[1].Seccoes = []Seccao{{Slug: "seccao-um", Nome: "Outra Secção", FonteID: "fonte-teste"}}
		}, []string{`secção "seccao-um"`, "duplicado"}},
		{"id de fonte duplicado", func(a *Arquivo) { a.Fontes[1].ID = "fonte-teste" },
			[]string{`fonte "fonte-teste"`, "duplicado"}},
		{"id de fonte com formato inválido", func(a *Arquivo) {
			a.Fontes[1].ID = "Fonte 2"
			a.Regioes[0].Setores[1].FonteID = "Fonte 2"
		}, []string{`fonte "Fonte 2"`, "id"}},

		// (4) referências a fontes
		{"fonte inexistente", func(a *Arquivo) { a.Regioes[0].Setores[0].FonteID = "nao-existe" },
			[]string{`setor "regiao-a" na região "regiao-a"`, `fonte "nao-existe" inexistente`}},
		{"fonte não usada", func(a *Arquivo) {
			a.Fontes = append(a.Fontes, Fonte{ID: "sem-uso", Nome: "Sem uso", Referencia: "x", ConsultadoEm: "2026-01-01"})
		}, []string{`fonte "sem-uso"`, "não é usada"}},

		{"fonte da capital inexistente", func(a *Arquivo) { a.Regioes[0].CapitalFonteID = "nao-existe" },
			[]string{`região "regiao-a" (capital)`, `fonte "nao-existe" inexistente`}},

		// (5) tipos
		{"tipo inválido", func(a *Arquivo) { a.Regioes[0].Tipo = "provincia" },
			[]string{`região "regiao-a"`, `tipo "provincia"`}},
		{"nenhum setor autónomo", func(a *Arquivo) { a.Regioes[1].Tipo = TipoRegiao },
			[]string{"setor-autonomo", "exatamente uma"}},
		{"dois setores autónomos", func(a *Arquivo) { a.Regioes[0].Tipo = TipoSetorAutonomo },
			[]string{"setor-autonomo", "exatamente uma"}},

		// (6) nomes duplicados no mesmo pai, datas e URLs
		{"nomes de setores duplicados na mesma região", func(a *Arquivo) { a.Regioes[0].Setores[1].Nome = "REGIÃO a" },
			[]string{`região "regiao-a"`, "nome duplicado", "REGIÃO a"}},
		{"nomes de regiões duplicados", func(a *Arquivo) { a.Regioes[1].Nome = "Regiao A" },
			[]string{"nome duplicado", "Regiao A"}},
		{"nomes de secções duplicados no mesmo setor", func(a *Arquivo) {
			a.Regioes[0].Setores[0].Seccoes[1].Nome = "Seccao um"
		}, []string{`setor "regiao-a"`, "nome duplicado"}},
		{"data com formato inválido", func(a *Arquivo) { a.Fontes[0].ConsultadoEm = "21/10/2026" },
			[]string{`fonte "fonte-teste"`, "consultado_em", "AAAA-MM-DD"}},
		{"url sem http", func(a *Arquivo) { a.Fontes[0].URL = "exemplo.org" },
			[]string{`fonte "fonte-teste"`, "url"}},
	}
	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			_, err := Carregar(variante(t, c.alterar))
			if err == nil {
				t.Fatal("esperava erro, obteve nil")
			}
			for _, trecho := range c.contem {
				if !strings.Contains(err.Error(), trecho) {
					t.Errorf("erro %q não contém %q", err.Error(), trecho)
				}
			}
		})
	}
}
