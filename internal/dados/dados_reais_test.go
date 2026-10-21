package dados

import (
	"bufio"
	"os"
	"strings"
	"testing"

	brutos "github.com/isconaifa/api-guine-bissau/dados"
)

func carregarReais(t *testing.T) *Base {
	t.Helper()
	b, err := Carregar(brutos.DivisaoAdministrativa)
	if err != nil {
		t.Fatalf("dados/divisao-administrativa.json inválido:\n%v", err)
	}
	return b
}

func TestDadosReais(t *testing.T) {
	b := carregarReais(t)
	if n := len(b.Regioes()); n != 9 {
		t.Errorf("entidades de 1.º nível = %d; esperado 9 (8 regiões + setor autónomo)", n)
	}
	autonomos := 0
	for _, r := range b.Regioes() {
		if r.Tipo == TipoSetorAutonomo {
			autonomos++
		}
	}
	if autonomos != 1 {
		t.Errorf("entidades %q = %d; esperado 1", TipoSetorAutonomo, autonomos)
	}
	for _, l := range b.Localidades() {
		if _, ok := b.Fonte(l.FonteID); !ok {
			t.Errorf("%s %q sem fonte válida (%q)", l.Tipo, l.Slug, l.FonteID)
		}
	}
}

// TestSlugsPublicadosEstaveis garante que nenhum slug publicado desaparece ou muda de nome
// (Princípios IV e VII). Slugs novos são permitidos, mas devem ser acrescentados à lista.
func TestSlugsPublicadosEstaveis(t *testing.T) {
	b := carregarReais(t)

	f, err := os.Open("testdata/slugs-publicados.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	publicados := map[string]bool{}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		linha := strings.TrimSpace(sc.Text())
		if linha == "" || strings.HasPrefix(linha, "#") {
			continue
		}
		publicados[linha] = true
		tipo, caminho, _ := strings.Cut(linha, ":")
		partes := strings.Split(caminho, "/")
		var existe bool
		switch {
		case tipo == "regiao" && len(partes) == 1:
			_, existe = b.Regiao(partes[0])
		case tipo == "setor" && len(partes) == 2:
			_, _, existe = b.Setor(partes[0], partes[1])
		case tipo == "seccao" && len(partes) == 3:
			if s, _, ok := b.Setor(partes[0], partes[1]); ok {
				for _, sc := range s.Seccoes {
					existe = existe || sc.Slug == partes[2]
				}
			}
		default:
			t.Errorf("linha inválida em slugs-publicados.txt: %q", linha)
			continue
		}
		if !existe {
			t.Errorf("slug publicado removido ou renomeado: %q (só permitido numa nova versão da API)", linha)
		}
	}
	if err := sc.Err(); err != nil {
		t.Fatal(err)
	}

	for _, r := range b.Regioes() {
		avisarNovo(t, publicados, "regiao:"+r.Slug)
		for _, s := range r.Setores {
			avisarNovo(t, publicados, "setor:"+r.Slug+"/"+s.Slug)
			for _, sc := range s.Seccoes {
				avisarNovo(t, publicados, "seccao:"+r.Slug+"/"+s.Slug+"/"+sc.Slug)
			}
		}
	}
}

func avisarNovo(t *testing.T, publicados map[string]bool, linha string) {
	if !publicados[linha] {
		t.Logf("slug novo, acrescente a testdata/slugs-publicados.txt: %s", linha)
	}
}
