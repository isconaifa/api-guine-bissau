// Package dados expõe o ficheiro da divisão administrativa embutido no binário.
package dados

import _ "embed"

// DivisaoAdministrativa é o conteúdo de divisao-administrativa.json.
//
//go:embed divisao-administrativa.json
var DivisaoAdministrativa []byte
