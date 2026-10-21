package respostas

import "github.com/isconaifa/api-guine-bissau/internal/dados"

// RegiaoResumo é uma entidade de primeiro nível sem a lista de setores.
type RegiaoResumo struct {
	Slug           string `json:"slug"`
	Nome           string `json:"nome"`
	Tipo           string `json:"tipo"`
	Capital        string `json:"capital"`
	CapitalFonteID string `json:"capital_fonte_id,omitempty"`
	FonteID        string `json:"fonte_id"`
	Verificado     bool   `json:"verificado"`
}

// RegiaoDetalhe é uma entidade de primeiro nível com os seus setores.
type RegiaoDetalhe struct {
	RegiaoResumo
	Setores []Referencia `json:"setores"`
}

// SetorResumo é um setor sem as suas secções.
type SetorResumo struct {
	Slug       string `json:"slug"`
	Nome       string `json:"nome"`
	FonteID    string `json:"fonte_id"`
	Verificado bool   `json:"verificado"`
}

func resumoRegiao(r *dados.Regiao) RegiaoResumo {
	return RegiaoResumo{
		Slug: r.Slug, Nome: r.Nome, Tipo: r.Tipo, Capital: r.Capital,
		CapitalFonteID: r.CapitalFonteID, FonteID: r.FonteID, Verificado: r.Verificado,
	}
}

func resumoSetor(s *dados.Setor) SetorResumo {
	return SetorResumo{Slug: s.Slug, Nome: s.Nome, FonteID: s.FonteID, Verificado: s.Verificado}
}

// ListaRegioes responde a GET /v1/regioes.
func ListaRegioes(base *dados.Base) Lista {
	regioes := base.Regioes()
	itens := make([]RegiaoResumo, len(regioes))
	var fontes []string
	for i := range regioes {
		itens[i] = resumoRegiao(&regioes[i])
		fontes = append(fontes, regioes[i].FonteID, regioes[i].CapitalFonteID)
	}
	return NovaLista(itens, ResumoFontes(base, fontes...))
}

// DetalheRegiao responde a GET /v1/regioes/{regiao}.
func DetalheRegiao(base *dados.Base, r *dados.Regiao) Objeto {
	setores := make([]Referencia, len(r.Setores))
	fontes := []string{r.FonteID, r.CapitalFonteID}
	for i, s := range r.Setores {
		setores[i] = Referencia{Slug: s.Slug, Nome: s.Nome}
		fontes = append(fontes, s.FonteID)
	}
	return Objeto{
		Dados: RegiaoDetalhe{RegiaoResumo: resumoRegiao(r), Setores: setores},
		Fonte: ResumoFontes(base, fontes...),
	}
}

// SetoresDaRegiao responde a GET /v1/regioes/{regiao}/setores. Uma lista vazia cita a fonte da
// própria região.
func SetoresDaRegiao(base *dados.Base, r *dados.Regiao) Lista {
	itens := make([]SetorResumo, len(r.Setores))
	fontes := make([]string, len(r.Setores))
	for i := range r.Setores {
		itens[i] = resumoSetor(&r.Setores[i])
		fontes[i] = r.Setores[i].FonteID
	}
	if len(fontes) == 0 {
		fontes = []string{r.FonteID}
	}
	return NovaLista(itens, ResumoFontes(base, fontes...))
}
