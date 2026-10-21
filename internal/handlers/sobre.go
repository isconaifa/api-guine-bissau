package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isconaifa/api-guine-bissau/internal/respostas"
)

func (a *api) sobre(c *gin.Context) {
	responder(c, http.StatusOK, respostas.InformacoesProjeto(a.base))
}
