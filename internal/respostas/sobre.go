package respostas

import (
	"github.com/isconaifa/api-guine-bissau/internal/dados"
	"github.com/isconaifa/api-guine-bissau/internal/sobre"
)

// Sobre descreve o projeto: autoria, licença, versão e fontes dos dados.
type Sobre struct {
	Nome        string        `json:"nome"`
	Autor       string        `json:"autor"`
	Licenca     string        `json:"licenca"`
	Versao      string        `json:"versao"`
	VersaoAPI   string        `json:"versao_api"`
	Repositorio string        `json:"repositorio"`
	Fontes      []dados.Fonte `json:"fontes"`
}

// InformacoesProjeto responde a GET /v1/sobre.
func InformacoesProjeto(base *dados.Base) Objeto {
	return Objeto{Dados: Sobre{
		Nome:        sobre.Nome,
		Autor:       sobre.Autor,
		Licenca:     sobre.Licenca,
		Versao:      sobre.Versao,
		VersaoAPI:   sobre.VersaoAPI,
		Repositorio: sobre.Repositorio,
		Fontes:      base.Fontes(),
	}}
}
