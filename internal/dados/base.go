package dados

import "github.com/isconaifa/api-guine-bissau/internal/pesquisa"

// Base guarda os dados validados e os índices em memória. É só de leitura e segura para uso
// concorrente.
type Base struct {
	fontes        []Fonte
	porIDFonte    map[string]Fonte
	regioes       []Regiao
	porSlugRegiao map[string]*Regiao
	porSlugSetor  map[string]setorIndexado
	localidades   []Localidade
	itens         []pesquisa.Item // paralelo a localidades
}

type setorIndexado struct {
	setor  *Setor
	regiao *Regiao
}

func novaBase(a Arquivo) *Base {
	b := &Base{
		fontes:        a.Fontes,
		porIDFonte:    make(map[string]Fonte, len(a.Fontes)),
		regioes:       a.Regioes,
		porSlugRegiao: make(map[string]*Regiao, len(a.Regioes)),
		porSlugSetor:  map[string]setorIndexado{},
	}
	for _, f := range b.fontes {
		b.porIDFonte[f.ID] = f
	}
	for i := range b.regioes {
		r := &b.regioes[i]
		if r.Setores == nil {
			r.Setores = []Setor{}
		}
		b.porSlugRegiao[r.Slug] = r
		b.localidades = append(b.localidades, Localidade{
			Tipo: r.Tipo, Slug: r.Slug, Nome: r.Nome, NomeNormalizado: pesquisa.Normalizar(r.Nome),
			FonteID: r.FonteID, Verificado: r.Verificado,
		})
		refRegiao := &Ref{Slug: r.Slug, Nome: r.Nome}
		for j := range r.Setores {
			s := &r.Setores[j]
			if s.Seccoes == nil {
				s.Seccoes = []Seccao{}
			}
			b.porSlugSetor[s.Slug] = setorIndexado{setor: s, regiao: r}
			b.localidades = append(b.localidades, Localidade{
				Tipo: TipoSetor, Slug: s.Slug, Nome: s.Nome, NomeNormalizado: pesquisa.Normalizar(s.Nome),
				Regiao: refRegiao, FonteID: s.FonteID, Verificado: s.Verificado,
			})
			refSetor := &Ref{Slug: s.Slug, Nome: s.Nome}
			for _, sc := range s.Seccoes {
				b.localidades = append(b.localidades, Localidade{
					Tipo: TipoSeccao, Slug: sc.Slug, Nome: sc.Nome, NomeNormalizado: pesquisa.Normalizar(sc.Nome),
					Regiao: refRegiao, Setor: refSetor, FonteID: sc.FonteID, Verificado: sc.Verificado,
				})
			}
		}
	}
	for _, l := range b.localidades {
		b.itens = append(b.itens, pesquisa.Item{NomeNormalizado: l.NomeNormalizado, Nivel: nivel(l.Tipo)})
	}
	return b
}

func nivel(tipo string) int {
	switch tipo {
	case TipoSetor:
		return 1
	case TipoSeccao:
		return 2
	default:
		return 0
	}
}

// Fontes devolve as fontes na ordem do ficheiro.
func (b *Base) Fontes() []Fonte { return b.fontes }

// Fonte procura uma fonte pelo id.
func (b *Base) Fonte(id string) (Fonte, bool) {
	f, ok := b.porIDFonte[id]
	return f, ok
}

// Regioes devolve as entidades de primeiro nível na ordem do ficheiro.
func (b *Base) Regioes() []Regiao { return b.regioes }

// Regiao procura uma entidade de primeiro nível pelo slug.
func (b *Base) Regiao(slug string) (*Regiao, bool) {
	r, ok := b.porSlugRegiao[slug]
	return r, ok
}

// Setor procura um setor pelo slug, exigindo que pertença à região indicada.
func (b *Base) Setor(slugRegiao, slugSetor string) (*Setor, *Regiao, bool) {
	s, ok := b.porSlugSetor[slugSetor]
	if !ok || s.regiao.Slug != slugRegiao {
		return nil, nil, false
	}
	return s.setor, s.regiao, true
}

// Localidades devolve a lista plana de todas as regiões, setores e secções.
func (b *Base) Localidades() []Localidade { return b.localidades }

// Pesquisar devolve as localidades cujo nome contém q, ignorando acentos e maiúsculas, pela
// ordem de pesquisa.Pesquisar.
func (b *Base) Pesquisar(q string) []Localidade {
	indices := pesquisa.Pesquisar(b.itens, q)
	r := make([]Localidade, len(indices))
	for i, j := range indices {
		r[i] = b.localidades[j]
	}
	return r
}
