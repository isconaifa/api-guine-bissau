package respostas

import "github.com/isconaifa/api-guine-bissau/internal/dados"

// ResultadoPesquisa é uma localidade encontrada, com a sua posição na hierarquia e a rota
// onde pode ser consultada.
type ResultadoPesquisa struct {
	Tipo       string      `json:"tipo"`
	Slug       string      `json:"slug"`
	Nome       string      `json:"nome"`
	Regiao     *Referencia `json:"regiao,omitempty"`
	Setor      *Referencia `json:"setor,omitempty"`
	Caminho    string      `json:"caminho"`
	FonteID    string      `json:"fonte_id"`
	Verificado bool        `json:"verificado"`
}

// NovoResultadoPesquisa converte uma localidade num resultado de pesquisa.
func NovoResultadoPesquisa(l dados.Localidade) ResultadoPesquisa {
	r := ResultadoPesquisa{Tipo: l.Tipo, Slug: l.Slug, Nome: l.Nome, FonteID: l.FonteID, Verificado: l.Verificado}
	switch l.Tipo {
	case dados.TipoSetor:
		r.Regiao = &Referencia{Slug: l.Regiao.Slug, Nome: l.Regiao.Nome}
		r.Caminho = "/v1/regioes/" + l.Regiao.Slug + "/setores/" + l.Slug
	case dados.TipoSeccao:
		r.Regiao = &Referencia{Slug: l.Regiao.Slug, Nome: l.Regiao.Nome}
		r.Setor = &Referencia{Slug: l.Setor.Slug, Nome: l.Setor.Nome}
		r.Caminho = "/v1/regioes/" + l.Regiao.Slug + "/setores/" + l.Setor.Slug + "/seccoes"
	default:
		r.Caminho = "/v1/regioes/" + l.Slug
	}
	return r
}

// ResultadosPesquisa responde a GET /v1/pesquisa.
func ResultadosPesquisa(base *dados.Base, locs []dados.Localidade) Lista {
	itens := make([]ResultadoPesquisa, len(locs))
	fontes := make([]string, len(locs))
	for i, l := range locs {
		itens[i] = NovoResultadoPesquisa(l)
		fontes[i] = l.FonteID
	}
	return NovaLista(itens, ResumoFontes(base, fontes...))
}
