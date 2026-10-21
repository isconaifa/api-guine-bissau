package pesquisa

import (
	"sort"
	"strings"
)

// Item é uma entrada pesquisável: o nome já normalizado e o nível na hierarquia
// (0 = região ou setor autónomo, 1 = setor, 2 = secção).
type Item struct {
	NomeNormalizado string
	Nivel           int
}

// Pesquisar devolve os índices dos itens cujo nome contém q (ignorando acentos e maiúsculas),
// ordenados por correspondência exata, depois prefixo, depois "contém"; em empate, por nível e
// por nome. Sem correspondência, devolve um slice vazio (nunca nil).
func Pesquisar(itens []Item, q string) []int {
	termo := Normalizar(q)
	resultado := []int{}
	if termo == "" {
		return resultado
	}
	classe := make(map[int]int)
	for i, it := range itens {
		switch {
		case it.NomeNormalizado == termo:
			classe[i] = 0
		case strings.HasPrefix(it.NomeNormalizado, termo):
			classe[i] = 1
		case strings.Contains(it.NomeNormalizado, termo):
			classe[i] = 2
		default:
			continue
		}
		resultado = append(resultado, i)
	}
	sort.SliceStable(resultado, func(a, b int) bool {
		ia, ib := resultado[a], resultado[b]
		if classe[ia] != classe[ib] {
			return classe[ia] < classe[ib]
		}
		if itens[ia].Nivel != itens[ib].Nivel {
			return itens[ia].Nivel < itens[ib].Nivel
		}
		return itens[ia].NomeNormalizado < itens[ib].NomeNormalizado
	})
	return resultado
}
