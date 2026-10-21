package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// cors permite pedidos de qualquer origem: a API é pública e só de leitura. Responde também à
// pré-verificação (OPTIONS), incluindo Private Network Access, usada pelos browsers quando uma
// página pública (ex.: o Swagger no GitHub Pages) chama um servidor local.
func cors(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	if c.Request.Method != http.MethodOptions {
		c.Next()
		return
	}
	c.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
	if h := c.GetHeader("Access-Control-Request-Headers"); h != "" {
		c.Header("Access-Control-Allow-Headers", h)
	}
	if c.GetHeader("Access-Control-Request-Private-Network") == "true" {
		c.Header("Access-Control-Allow-Private-Network", "true")
	}
	c.Header("Access-Control-Max-Age", "86400")
	c.AbortWithStatus(http.StatusNoContent)
}

// cabecalhosDados permite cache dos dados, que só mudam com uma nova versão do projeto.
func cabecalhosDados(c *gin.Context) {
	c.Header("Cache-Control", "public, max-age=3600")
	c.Next()
}
