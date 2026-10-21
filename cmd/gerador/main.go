// Comando gerador escreve a versão estática da API (ficheiros JSON) para o GitHub Pages.
package main

import (
	"flag"
	"log"
	"os"

	brutos "github.com/isconaifa/api-guine-bissau/dados"
	"github.com/isconaifa/api-guine-bissau/internal/dados"
)

func main() {
	saida := flag.String("saida", "publico", "pasta de saída (apagada e recriada)")
	flag.Parse()

	base, err := dados.Carregar(brutos.DivisaoAdministrativa)
	if err != nil {
		log.Fatal(err)
	}
	if err := os.RemoveAll(*saida); err != nil {
		log.Fatal(err)
	}
	if err := Gerar(base, *saida); err != nil {
		log.Fatal(err)
	}
	log.Printf("versão estática gerada em %s", *saida)
}
