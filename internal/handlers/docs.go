package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isconaifa/api-guine-bissau/docs"
)

// A documentação é servida pelo próprio servidor para que "Try it out" funcione sem CORS.

func paginaDocs(c *gin.Context) {
	c.Data(http.StatusOK, "text/html; charset=utf-8", docs.SwaggerUI)
}

func contratoOpenAPI(c *gin.Context) {
	c.Data(http.StatusOK, "application/yaml; charset=utf-8", docs.OpenAPI)
}
