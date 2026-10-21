// Package handlers expõe a API HTTP v1 com o Gin.
package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isconaifa/api-guine-bissau/internal/dados"
	"github.com/isconaifa/api-guine-bissau/internal/respostas"
)

// api agrupa os handlers e os dados que servem.
type api struct {
	base *dados.Base
}

// NovoRouter cria o router com todas as rotas /v1, erros no envelope comum e CORS aberto.
func NovoRouter(base *dados.Base) *gin.Engine {
	a := &api{base: base}

	r := gin.New()
	r.Use(cors, gin.CustomRecovery(func(c *gin.Context, recuperado any) {
		log.Printf("pânico em %s: %v", c.Request.URL.Path, recuperado)
		responderErro(c, respostas.ErroInterno())
	}))
	semRota := func(c *gin.Context) {
		responderErro(c, respostas.RotaNaoEncontrada(c.Request.URL.Path))
	}
	r.NoRoute(semRota)
	r.NoMethod(semRota)

	r.GET("/v1/saude", a.saude)

	v1 := r.Group("/v1", cabecalhosDados)
	v1.GET("/regioes", a.listarRegioes)
	v1.GET("/regioes/:regiao", a.detalharRegiao)
	v1.GET("/regioes/:regiao/setores", a.listarSetores)
	v1.GET("/regioes/:regiao/setores/:setor", a.detalharSetor)
	v1.GET("/regioes/:regiao/setores/:setor/seccoes", a.listarSeccoes)
	v1.GET("/pesquisa", a.pesquisar)
	v1.GET("/sobre", a.sobre)

	return r
}

func (a *api) saude(c *gin.Context) {
	c.Header("Cache-Control", "no-store")
	responder(c, http.StatusOK, respostas.Objeto{Dados: map[string]string{"estado": "ok"}})
}
