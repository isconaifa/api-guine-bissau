package handlers

import "github.com/gin-gonic/gin"

// cors permite pedidos de qualquer origem: a API é pública e só de leitura.
func cors(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Next()
}

// cabecalhosDados permite cache dos dados, que só mudam com uma nova versão do projeto.
func cabecalhosDados(c *gin.Context) {
	c.Header("Cache-Control", "public, max-age=3600")
	c.Next()
}
