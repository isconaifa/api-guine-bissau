package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isconaifa/api-guine-bissau/internal/respostas"
)

const tipoJSON = "application/json; charset=utf-8"

// responder escreve v como JSON com respostas.Serializar. Todas as respostas da API passam
// por aqui (nunca c.JSON), para serem idênticas aos ficheiros estáticos.
func responder(c *gin.Context, status int, v any) {
	b, err := respostas.Serializar(v)
	if err != nil {
		log.Printf("serializar resposta de %s: %v", c.Request.URL.Path, err)
		status = http.StatusInternalServerError
		b, _ = respostas.Serializar(respostas.ErroInterno())
	}
	c.Data(status, tipoJSON, b)
}

// responderErro escreve um erro no envelope comum, com o status HTTP do próprio erro.
func responderErro(c *gin.Context, e respostas.Erro) {
	responder(c, e.Status, e)
	c.Abort()
}
