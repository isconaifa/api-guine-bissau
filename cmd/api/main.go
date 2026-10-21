// Comando api inicia o servidor HTTP da API Guiné-Bissau.
package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	brutos "github.com/isconaifa/api-guine-bissau/dados"
	"github.com/isconaifa/api-guine-bissau/internal/dados"
	"github.com/isconaifa/api-guine-bissau/internal/handlers"
	"github.com/isconaifa/api-guine-bissau/internal/sobre"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	base, err := dados.Carregar(brutos.DivisaoAdministrativa)
	if err != nil {
		log.Fatal(err)
	}

	porta := os.Getenv("PORT")
	if porta == "" {
		porta = "8080"
	}

	srv := &http.Server{
		Addr:              ":" + porta,
		Handler:           handlers.NovoRouter(base),
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Printf("%s %s a escutar na porta %s (%d entidades de 1.º nível)",
		sobre.Nome, sobre.Versao, porta, len(base.Regioes()))
	log.Fatal(srv.ListenAndServe())
}
