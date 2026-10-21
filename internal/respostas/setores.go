package respostas

import "github.com/isconaifa/api-guine-bissau/internal/dados"

// SetorDetalhe é um setor com a região a que pertence.
type SetorDetalhe struct {
	SetorResumo
	Regiao       Referencia `json:"regiao"`
	TotalSeccoes int        `json:"total_seccoes"`
}

// Seccao é uma secção de um setor.
type Seccao struct {
	Slug       string `json:"slug"`
	Nome       string `json:"nome"`
	FonteID    string `json:"fonte_id"`
	Verificado bool   `json:"verificado"`
}

// DetalheSetor responde a GET /v1/regioes/{regiao}/setores/{setor}.
func DetalheSetor(base *dados.Base, s *dados.Setor, r *dados.Regiao) Objeto {
	return Objeto{
		Dados: SetorDetalhe{
			SetorResumo:  resumoSetor(s),
			Regiao:       Referencia{Slug: r.Slug, Nome: r.Nome},
			TotalSeccoes: len(s.Seccoes),
		},
		Fonte: ResumoFontes(base, s.FonteID),
	}
}

// SeccoesDoSetor responde a GET /v1/regioes/{regiao}/setores/{setor}/seccoes. Uma lista vazia
// cita a fonte do próprio setor.
func SeccoesDoSetor(base *dados.Base, s *dados.Setor) Lista {
	itens := make([]Seccao, len(s.Seccoes))
	fontes := make([]string, len(s.Seccoes))
	for i, sc := range s.Seccoes {
		itens[i] = Seccao{Slug: sc.Slug, Nome: sc.Nome, FonteID: sc.FonteID, Verificado: sc.Verificado}
		fontes[i] = sc.FonteID
	}
	if len(fontes) == 0 {
		fontes = []string{s.FonteID}
	}
	return NovaLista(itens, ResumoFontes(base, fontes...))
}
