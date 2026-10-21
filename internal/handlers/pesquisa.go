package handlers

import (
	"net/http"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/isconaifa/api-guine-bissau/internal/pesquisa"
	"github.com/isconaifa/api-guine-bissau/internal/respostas"
)

const maxQ = 100

func (a *api) pesquisar(c *gin.Context) {
	q, _ := c.GetQuery("q")
	if utf8.RuneCountInString(q) > maxQ || pesquisa.Normalizar(q) == "" {
		responderErro(c, respostas.ParametroInvalido(
			"O parâmetro 'q' é obrigatório e deve ter entre 1 e 100 caracteres."))
		return
	}
	responder(c, http.StatusOK, respostas.ResultadosPesquisa(a.base, a.base.Pesquisar(q)))
}
