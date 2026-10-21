package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isconaifa/api-guine-bissau/internal/dados"
	"github.com/isconaifa/api-guine-bissau/internal/respostas"
)

func (a *api) listarRegioes(c *gin.Context) {
	responder(c, http.StatusOK, respostas.ListaRegioes(a.base))
}

func (a *api) detalharRegiao(c *gin.Context) {
	if r, ok := a.regiao(c); ok {
		responder(c, http.StatusOK, respostas.DetalheRegiao(a.base, r))
	}
}

func (a *api) listarSetores(c *gin.Context) {
	if r, ok := a.regiao(c); ok {
		responder(c, http.StatusOK, respostas.SetoresDaRegiao(a.base, r))
	}
}

// regiao procura a região do caminho; se não existir, responde com o erro e devolve false.
func (a *api) regiao(c *gin.Context) (*dados.Regiao, bool) {
	slug := c.Param("regiao")
	r, ok := a.base.Regiao(slug)
	if !ok {
		responderErro(c, respostas.RegiaoNaoEncontrada(slug))
	}
	return r, ok
}
