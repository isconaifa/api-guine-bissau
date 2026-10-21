// Package sobre contém os metadados do projeto devolvidos por /v1/sobre.
package sobre

// Versao é a versão do projeto, injetada no build com
// -ldflags "-X github.com/isconaifa/api-guine-bissau/internal/sobre.Versao=vX.Y.Z".
var Versao = "dev"

const (
	Nome        = "API Guiné-Bissau"
	Autor       = "Francisco Naifa"
	Licenca     = "MIT"
	VersaoAPI   = "v1"
	Repositorio = "https://github.com/isconaifa/api-guine-bissau"
)
