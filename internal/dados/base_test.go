package dados

import "testing"

func carregarValido(t *testing.T) *Base {
	t.Helper()
	b, err := Carregar(lerValido(t))
	if err != nil {
		t.Fatalf("Carregar: %v", err)
	}
	return b
}

func TestBaseRegioes(t *testing.T) {
	b := carregarValido(t)
	rs := b.Regioes()
	if len(rs) != 2 || rs[0].Slug != "regiao-a" || rs[1].Slug != "setor-autonomo-teste" {
		t.Fatalf("regiões fora da ordem do ficheiro: %+v", rs)
	}
	if r, ok := b.Regiao("regiao-a"); !ok || r.Nome != "Região A" {
		t.Errorf("Regiao(regiao-a) = %v, %v", r, ok)
	}
	if _, ok := b.Regiao("xyz"); ok {
		t.Error("Regiao(xyz) deveria falhar")
	}
	if r, _ := b.Regiao("setor-autonomo-teste"); r.Setores == nil {
		t.Error("lista de setores vazia deve ser [] e não nil")
	}
}

func TestBaseSetor(t *testing.T) {
	b := carregarValido(t)
	s, r, ok := b.Setor("regiao-a", "setor-b")
	if !ok || s.Nome != "Setor B" || r.Slug != "regiao-a" {
		t.Fatalf("Setor(regiao-a, setor-b) = %v, %v, %v", s, r, ok)
	}
	if s.Seccoes == nil {
		t.Error("lista de secções vazia deve ser [] e não nil")
	}
	// Região e setor com o mesmo slug coexistem (unicidade por tipo).
	if _, _, ok := b.Setor("regiao-a", "regiao-a"); !ok {
		t.Error("setor 'regiao-a' na região 'regiao-a' deveria existir")
	}
	if _, _, ok := b.Setor("setor-autonomo-teste", "setor-b"); ok {
		t.Error("setor pedido na região errada deveria falhar")
	}
	if _, _, ok := b.Setor("regiao-a", "xyz"); ok {
		t.Error("setor inexistente deveria falhar")
	}
}

func TestBaseFontes(t *testing.T) {
	b := carregarValido(t)
	if len(b.Fontes()) != 2 {
		t.Fatalf("fontes = %d; esperado 2", len(b.Fontes()))
	}
	if f, ok := b.Fonte("fonte-secundaria"); !ok || f.Nome != "Fonte Secundária & Teste" {
		t.Errorf("Fonte(fonte-secundaria) = %v, %v", f, ok)
	}
}

func TestBaseLocalidades(t *testing.T) {
	b := carregarValido(t)
	locs := b.Localidades()
	// 2 regiões + 2 setores + 2 secções
	if len(locs) != 6 {
		t.Fatalf("localidades = %d; esperado 6", len(locs))
	}
	tipos := map[string]int{}
	for _, l := range locs {
		tipos[l.Tipo]++
		switch l.Tipo {
		case TipoSetor:
			if l.Regiao == nil || l.Setor != nil {
				t.Errorf("setor %q com hierarquia errada", l.Slug)
			}
		case TipoSeccao:
			if l.Regiao == nil || l.Setor == nil {
				t.Errorf("secção %q sem região ou setor", l.Slug)
			}
		default:
			if l.Regiao != nil || l.Setor != nil {
				t.Errorf("entidade de 1.º nível %q com pai", l.Slug)
			}
		}
		if l.NomeNormalizado == "" {
			t.Errorf("%q sem nome normalizado", l.Slug)
		}
	}
	esperado := map[string]int{TipoRegiao: 1, TipoSetorAutonomo: 1, TipoSetor: 2, TipoSeccao: 2}
	for tipo, n := range esperado {
		if tipos[tipo] != n {
			t.Errorf("tipo %s: %d; esperado %d", tipo, tipos[tipo], n)
		}
	}
}

func TestBasePesquisar(t *testing.T) {
	b := carregarValido(t)
	r := b.Pesquisar("REGIÃO a")
	if len(r) != 2 || r[0].Tipo != TipoRegiao || r[1].Tipo != TipoSetor {
		t.Errorf("Pesquisar(REGIÃO a) = %+v; esperado região antes do setor", r)
	}
	if r := b.Pesquisar("secção"); len(r) != 2 || r[0].Slug != "seccao-dois" || r[1].Slug != "seccao-um" {
		t.Errorf("Pesquisar(secção) = %+v", r)
	}
	if r := b.Pesquisar("zzz"); r == nil || len(r) != 0 {
		t.Errorf("Pesquisar(zzz) = %#v; esperado vazio não-nil", r)
	}
}
