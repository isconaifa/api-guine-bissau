// Package dados carrega, valida e indexa em memória a divisão administrativa da Guiné-Bissau.
package dados

// Tipos de entidade de primeiro nível.
const (
	TipoRegiao        = "regiao"
	TipoSetorAutonomo = "setor-autonomo"
)

// Tipos de localidade abaixo do primeiro nível (usados na pesquisa).
const (
	TipoSetor  = "setor"
	TipoSeccao = "seccao"
)

// Arquivo é o conteúdo de dados/divisao-administrativa.json.
type Arquivo struct {
	Fontes  []Fonte  `json:"fontes"`
	Regioes []Regiao `json:"regioes"`
}

// Fonte é a origem identificada de um dado (ex.: INE Guiné-Bissau).
type Fonte struct {
	ID           string `json:"id"`
	Nome         string `json:"nome"`
	Referencia   string `json:"referencia"`
	URL          string `json:"url,omitempty"`
	ConsultadoEm string `json:"consultado_em"`
}

// Regiao é uma entidade de primeiro nível: uma região ou o Setor Autónomo de Bissau.
type Regiao struct {
	Slug    string `json:"slug"`
	Nome    string `json:"nome"`
	Tipo    string `json:"tipo"`
	Capital string `json:"capital"`
	// CapitalFonteID é a fonte da capital, quando diferente de FonteID (opcional).
	CapitalFonteID string  `json:"capital_fonte_id,omitempty"`
	FonteID        string  `json:"fonte_id"`
	Verificado     bool    `json:"verificado"`
	Setores        []Setor `json:"setores"`
}

// Setor é uma subdivisão de uma região.
type Setor struct {
	Slug       string   `json:"slug"`
	Nome       string   `json:"nome"`
	FonteID    string   `json:"fonte_id"`
	Verificado bool     `json:"verificado"`
	Seccoes    []Seccao `json:"seccoes"`
}

// Seccao é uma subdivisão de um setor.
type Seccao struct {
	Slug       string `json:"slug"`
	Nome       string `json:"nome"`
	FonteID    string `json:"fonte_id"`
	Verificado bool   `json:"verificado"`
}

// Ref identifica uma entidade superior (região ou setor) de uma localidade.
type Ref struct {
	Slug string
	Nome string
}

// Localidade é uma entrada da lista plana usada na pesquisa.
type Localidade struct {
	Tipo            string // TipoRegiao, TipoSetorAutonomo, TipoSetor ou TipoSeccao
	Slug            string
	Nome            string
	NomeNormalizado string
	Regiao          *Ref // presente em setores e secções
	Setor           *Ref // presente em secções
	FonteID         string
	Verificado      bool
}
