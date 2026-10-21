package dados

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/isconaifa/api-guine-bissau/internal/pesquisa"
)

// FormatoSlug é o formato obrigatório de slugs e ids de fonte.
var FormatoSlug = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

// validador acumula os problemas encontrados para os reportar todos de uma vez.
type validador struct {
	erros []error
}

func (v *validador) erro(onde, formato string, args ...any) {
	v.erros = append(v.erros, fmt.Errorf("%s: %s", onde, fmt.Sprintf(formato, args...)))
}

func (v *validador) obrigatorio(onde, campo, valor string) {
	if strings.TrimSpace(valor) == "" {
		v.erro(onde, "campo %q vazio ou ausente", campo)
	}
}

func (v *validador) slug(onde, slug string) {
	if slug != "" && !FormatoSlug.MatchString(slug) {
		v.erro(onde, "slug %q com formato inválido (use minúsculas, dígitos e hífens, ex.: 'bafata')", slug)
	}
}

// nomesUnicos regista um nome no conjunto do seu pai e acusa duplicados após normalização.
func (v *validador) nomesUnicos(vistos map[string]bool, pai, nome string) {
	n := pesquisa.Normalizar(nome)
	if n == "" {
		return
	}
	if vistos[n] {
		v.erro(pai, "nome duplicado %q", nome)
	}
	vistos[n] = true
}

func rotulo(tipo string, i int, slug string) string {
	if slug == "" {
		return fmt.Sprintf("%s #%d", tipo, i+1)
	}
	return fmt.Sprintf("%s %q", tipo, slug)
}

func validar(a *Arquivo) error {
	v := &validador{}

	fontes := map[string]bool{}
	usadas := map[string]bool{}
	for i, f := range a.Fontes {
		onde := rotulo("fonte", i, f.ID)
		v.obrigatorio(onde, "id", f.ID)
		if f.ID != "" && !FormatoSlug.MatchString(f.ID) {
			v.erro(onde, "id com formato inválido (use minúsculas, dígitos e hífens)")
		}
		if fontes[f.ID] && f.ID != "" {
			v.erro(onde, "id duplicado")
		}
		fontes[f.ID] = true
		v.obrigatorio(onde, "nome", f.Nome)
		v.obrigatorio(onde, "referencia", f.Referencia)
		v.obrigatorio(onde, "consultado_em", f.ConsultadoEm)
		if f.ConsultadoEm != "" {
			if _, err := time.Parse(time.DateOnly, f.ConsultadoEm); err != nil {
				v.erro(onde, "consultado_em %q deve estar no formato AAAA-MM-DD", f.ConsultadoEm)
			}
		}
		if f.URL != "" && !strings.HasPrefix(f.URL, "http://") && !strings.HasPrefix(f.URL, "https://") {
			v.erro(onde, "url %q deve começar por http:// ou https://", f.URL)
		}
	}

	fonte := func(onde, id string) {
		v.obrigatorio(onde, "fonte_id", id)
		if id == "" {
			return
		}
		usadas[id] = true
		if !fontes[id] {
			v.erro(onde, "fonte %q inexistente", id)
		}
	}

	slugsRegioes := map[string]bool{}
	slugsSetores := map[string]bool{}
	slugsSeccoes := map[string]bool{}
	nomesRegioes := map[string]bool{}
	autonomos := 0

	for i, r := range a.Regioes {
		onde := rotulo("região", i, r.Slug)
		v.obrigatorio(onde, "slug", r.Slug)
		v.slug(onde, r.Slug)
		if r.Slug != "" && slugsRegioes[r.Slug] {
			v.erro(onde, "slug duplicado entre regiões")
		}
		slugsRegioes[r.Slug] = true
		v.obrigatorio(onde, "nome", r.Nome)
		v.nomesUnicos(nomesRegioes, onde, r.Nome)
		v.obrigatorio(onde, "capital", r.Capital)
		switch r.Tipo {
		case TipoRegiao:
		case TipoSetorAutonomo:
			autonomos++
		default:
			v.erro(onde, "tipo %q inválido (use %q ou %q)", r.Tipo, TipoRegiao, TipoSetorAutonomo)
		}
		fonte(onde, r.FonteID)
		if r.CapitalFonteID != "" {
			fonte(onde+" (capital)", r.CapitalFonteID)
		}

		nomesSetores := map[string]bool{}
		for j, s := range r.Setores {
			ondeS := rotulo("setor", j, s.Slug) + " na " + onde
			v.obrigatorio(ondeS, "slug", s.Slug)
			v.slug(ondeS, s.Slug)
			if s.Slug != "" && slugsSetores[s.Slug] {
				v.erro(ondeS, "slug duplicado entre setores")
			}
			slugsSetores[s.Slug] = true
			v.obrigatorio(ondeS, "nome", s.Nome)
			v.nomesUnicos(nomesSetores, onde, s.Nome)
			fonte(ondeS, s.FonteID)

			nomesSeccoes := map[string]bool{}
			for k, sc := range s.Seccoes {
				ondeSc := fmt.Sprintf("%s no %s (%s)", rotulo("secção", k, sc.Slug), rotulo("setor", j, s.Slug), onde)
				v.obrigatorio(ondeSc, "slug", sc.Slug)
				v.slug(ondeSc, sc.Slug)
				if sc.Slug != "" && slugsSeccoes[sc.Slug] {
					v.erro(ondeSc, "slug duplicado entre secções")
				}
				slugsSeccoes[sc.Slug] = true
				v.obrigatorio(ondeSc, "nome", sc.Nome)
				v.nomesUnicos(nomesSeccoes, rotulo("setor", j, s.Slug), sc.Nome)
				fonte(ondeSc, sc.FonteID)
			}
		}
	}

	if autonomos != 1 {
		v.erros = append(v.erros, fmt.Errorf(
			"deve existir exatamente uma entidade com tipo %q (encontradas: %d)", TipoSetorAutonomo, autonomos))
	}
	for i, f := range a.Fontes {
		if f.ID != "" && !usadas[f.ID] {
			v.erro(rotulo("fonte", i, f.ID), "não é usada por nenhuma entidade")
		}
	}
	return errors.Join(v.erros...)
}
