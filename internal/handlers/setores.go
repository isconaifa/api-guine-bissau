package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isconaifa/api-guine-bissau/internal/dados"
	"github.com/isconaifa/api-guine-bissau/internal/respostas"
)

func (a *api) detalharSetor(c *gin.Context) {
	if s, r, ok := a.setor(c); ok {
		responder(c, http.StatusOK, respostas.DetalheSetor(a.base, s, r))
	}
}

func (a *api) listarSeccoes(c *gin.Context) {
	if s, _, ok := a.setor(c); ok {
		responder(c, http.StatusOK, respostas.SeccoesDoSetor(a.base, s))
	}
}

// setor procura primeiro a região e depois o setor do caminho; se algum não existir, responde
// com o erro correspondente e devolve false.
func (a *api) setor(c *gin.Context) (*dados.Setor, *dados.Regiao, bool) {
	r, ok := a.regiao(c)
	if !ok {
		return nil, nil, false
	}
	slug := c.Param("setor")
	s, _, ok := a.base.Setor(r.Slug, slug)
	if !ok {
		responderErro(c, respostas.SetorNaoEncontrado(r.Slug, slug))
	}
	return s, r, ok
}
