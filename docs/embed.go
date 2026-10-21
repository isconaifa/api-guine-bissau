// Package docs expõe a documentação da API (Swagger UI e contrato OpenAPI) embutida no binário.
package docs

import _ "embed"

// SwaggerUI é a página index.html com o Swagger UI, que lê openapi.yaml do mesmo diretório.
//
//go:embed index.html
var SwaggerUI []byte

// OpenAPI é o contrato da API v1 (openapi.yaml).
//
//go:embed openapi.yaml
var OpenAPI []byte
