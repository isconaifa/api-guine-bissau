// Package respostas define o envelope JSON comum a todas as respostas e os construtores
// partilhados pela API e pelo gerador estático. Não depende do Gin.
package respostas

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/isconaifa/api-guine-bissau/internal/dados"
)

// Lista é o envelope de respostas com uma lista de itens.
type Lista struct {
	Dados any    `json:"dados"`
	Total int    `json:"total"`
	Fonte string `json:"fonte"`
}

// NovaLista cria o envelope de uma lista; uma lista nil é devolvida como [].
func NovaLista[T any](itens []T, fonte string) Lista {
	if itens == nil {
		itens = []T{}
	}
	return Lista{Dados: itens, Total: len(itens), Fonte: fonte}
}

// Objeto é o envelope de respostas com um único objeto.
type Objeto struct {
	Dados any    `json:"dados"`
	Fonte string `json:"fonte,omitempty"`
}

// Erro é o envelope de respostas de erro; Status é o código HTTP e não é serializado.
type Erro struct {
	Status int         `json:"-"`
	Erro   DetalheErro `json:"erro"`
}

// DetalheErro descreve um erro com código estável e mensagem em português.
type DetalheErro struct {
	Codigo   string `json:"codigo"`
	Mensagem string `json:"mensagem"`
}

// Referencia identifica uma entidade pelo slug e pelo nome.
type Referencia struct {
	Slug string `json:"slug"`
	Nome string `json:"nome"`
}

// ResumoFontes devolve os nomes distintos das fontes indicadas, separados por "; ", na ordem
// da primeira ocorrência. Ids vazios ou desconhecidos são ignorados.
func ResumoFontes(base *dados.Base, ids ...string) string {
	vistos := map[string]bool{}
	var nomes []string
	for _, id := range ids {
		f, ok := base.Fonte(id)
		if !ok || vistos[f.Nome] {
			continue
		}
		vistos[f.Nome] = true
		nomes = append(nomes, f.Nome)
	}
	return strings.Join(nomes, "; ")
}

// Serializar produz o corpo JSON de uma resposta: sem escapar HTML, sem indentação e com
// quebra de linha final. É a única forma de gerar corpos JSON no projeto, para que a API e
// os ficheiros estáticos sejam idênticos byte a byte.
func Serializar(v any) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
