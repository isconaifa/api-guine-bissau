package respostas

import (
	"fmt"
	"net/http"

	"github.com/isconaifa/api-guine-bissau/internal/dados"
)

// Códigos de erro estáveis (ver contracts/openapi.yaml).
const (
	CodigoRegiaoNaoEncontrada = "REGIAO_NAO_ENCONTRADA"
	CodigoSetorNaoEncontrado  = "SETOR_NAO_ENCONTRADO"
	CodigoParametroInvalido   = "PARAMETRO_INVALIDO"
	CodigoRotaNaoEncontrada   = "ROTA_NAO_ENCONTRADA"
	CodigoErroInterno         = "ERRO_INTERNO"
)

const dicaSlug = " Use o slug em minúsculas e sem acentos (ex.: 'bafata')."

func novoErro(status int, codigo, mensagem string) Erro {
	return Erro{Status: status, Erro: DetalheErro{Codigo: codigo, Mensagem: mensagem}}
}

// comDica acrescenta a dica de formato quando o valor pedido não tem formato de slug.
func comDica(mensagem string, slugs ...string) string {
	for _, s := range slugs {
		if !dados.FormatoSlug.MatchString(s) {
			return mensagem + dicaSlug
		}
	}
	return mensagem
}

// RegiaoNaoEncontrada: não existe região com o slug pedido.
func RegiaoNaoEncontrada(slug string) Erro {
	return novoErro(http.StatusNotFound, CodigoRegiaoNaoEncontrada,
		comDica(fmt.Sprintf("Região '%s' não encontrada.", slug), slug))
}

// SetorNaoEncontrado: o setor não existe ou não pertence à região do caminho.
func SetorNaoEncontrado(regiao, setor string) Erro {
	return novoErro(http.StatusNotFound, CodigoSetorNaoEncontrado,
		comDica(fmt.Sprintf("Setor '%s' não encontrado na região '%s'.", setor, regiao), setor))
}

// ParametroInvalido: um parâmetro do pedido está ausente ou é inválido.
func ParametroInvalido(mensagem string) Erro {
	return novoErro(http.StatusBadRequest, CodigoParametroInvalido, mensagem)
}

// RotaNaoEncontrada: a rota ou o método não existe.
func RotaNaoEncontrada(caminho string) Erro {
	return novoErro(http.StatusNotFound, CodigoRotaNaoEncontrada,
		fmt.Sprintf("Rota '%s' não encontrada. Consulte a documentação no README.", caminho))
}

// ErroInterno: falha inesperada no servidor.
func ErroInterno() Erro {
	return novoErro(http.StatusInternalServerError, CodigoErroInterno, "Erro interno do servidor.")
}
