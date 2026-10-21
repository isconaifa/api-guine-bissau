package dados

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

// Carregar decodifica e valida o JSON da divisão administrativa e constrói os índices em
// memória. Dados inválidos nunca são servidos: qualquer violação das regras do modelo de
// dados devolve erro, com todos os problemas encontrados.
func Carregar(b []byte) (*Base, error) {
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields()
	var a Arquivo
	if err := dec.Decode(&a); err != nil {
		return nil, fmt.Errorf("dados: JSON inválido: %w", err)
	}
	if _, err := dec.Token(); !errors.Is(err, io.EOF) {
		return nil, errors.New("dados: JSON inválido: conteúdo depois do objeto principal")
	}
	if err := validar(&a); err != nil {
		return nil, fmt.Errorf("dados inválidos:\n%w", err)
	}
	return novaBase(a), nil
}
