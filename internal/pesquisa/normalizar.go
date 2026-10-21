// Package pesquisa contém a normalização de texto e a pesquisa de localidades por nome.
package pesquisa

import (
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

// Normalizar converte um texto para a forma usada na pesquisa: sem acentos, em minúsculas,
// com qualquer caractere que não seja letra ou dígito trocado por espaço e sem espaços
// repetidos. Ex.: "Bolama/Bijagós" → "bolama bijagos".
func Normalizar(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	espaco := true // evita espaço inicial e espaços repetidos
	for _, r := range norm.NFD.String(s) {
		switch {
		case unicode.Is(unicode.Mn, r):
			continue
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			b.WriteRune(unicode.ToLower(r))
			espaco = false
		case !espaco:
			b.WriteByte(' ')
			espaco = true
		}
	}
	return strings.TrimRight(b.String(), " ")
}
