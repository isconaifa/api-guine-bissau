package main

import (
	"os"
	"path/filepath"

	"github.com/isconaifa/api-guine-bissau/internal/dados"
	"github.com/isconaifa/api-guine-bissau/internal/respostas"
)

// itemIndice é um resultado de pesquisa com o nome normalizado, para pesquisa no cliente.
type itemIndice struct {
	respostas.ResultadoPesquisa
	NomeNormalizado string `json:"nome_normalizado"`
}

// Gerar escreve em dir um ficheiro .json por rota da API v1 (contracts/json-estatico.md),
// com o mesmo corpo que a API devolve, e o índice de pesquisa.
func Gerar(base *dados.Base, dir string) error {
	escrever := func(caminho string, v any) error {
		b, err := respostas.Serializar(v)
		if err != nil {
			return err
		}
		destino := filepath.Join(dir, filepath.FromSlash(caminho))
		if err := os.MkdirAll(filepath.Dir(destino), 0o755); err != nil {
			return err
		}
		return os.WriteFile(destino, b, 0o644)
	}

	if err := escrever("v1/regioes.json", respostas.ListaRegioes(base)); err != nil {
		return err
	}
	for i := range base.Regioes() {
		r := &base.Regioes()[i]
		prefixo := "v1/regioes/" + r.Slug
		if err := escrever(prefixo+".json", respostas.DetalheRegiao(base, r)); err != nil {
			return err
		}
		if err := escrever(prefixo+"/setores.json", respostas.SetoresDaRegiao(base, r)); err != nil {
			return err
		}
		for j := range r.Setores {
			s := &r.Setores[j]
			p := prefixo + "/setores/" + s.Slug
			if err := escrever(p+".json", respostas.DetalheSetor(base, s, r)); err != nil {
				return err
			}
			if err := escrever(p+"/seccoes.json", respostas.SeccoesDoSetor(base, s)); err != nil {
				return err
			}
		}
	}
	if err := escrever("v1/sobre.json", respostas.InformacoesProjeto(base)); err != nil {
		return err
	}

	locs := base.Localidades()
	itens := make([]itemIndice, len(locs))
	fontes := make([]string, len(locs))
	for i, l := range locs {
		itens[i] = itemIndice{ResultadoPesquisa: respostas.NovoResultadoPesquisa(l), NomeNormalizado: l.NomeNormalizado}
		fontes[i] = l.FonteID
	}
	return escrever("v1/pesquisa/indice.json", respostas.NovaLista(itens, respostas.ResumoFontes(base, fontes...)))
}
