# Data Model: Divisão Administrativa da Guiné-Bissau

**Feature**: `001-divisao-administrativa` | **Date**: 2026-10-21

Modelo do ficheiro `dados/divisao-administrativa.json` (fonte de verdade) e das estruturas em
memória construídas no arranque. No JSON, o campo `fonte_id` de cada entidade é o id da fonte; o campo `fonte` do envelope de
resposta é o nome legível da(s) fonte(s). Os formatos de resposta estão em
[contracts/openapi.yaml](contracts/openapi.yaml).

## Ficheiro de dados

```json
{
  "fontes": [
    {
      "id": "ine-rgph",
      "nome": "INE Guiné-Bissau",
      "referencia": "<título e ano do documento>",
      "url": "<url, se existir>",
      "consultado_em": "AAAA-MM-DD"
    }
  ],
  "regioes": [
    {
      "slug": "bafata",
      "nome": "Bafatá",
      "tipo": "regiao",
      "capital": "<capital, segundo a fonte>",
      "fonte_id": "ine-rgph",
      "verificado": false,
      "setores": [
        {
          "slug": "<slug>",
          "nome": "<nome oficial>",
          "fonte_id": "ine-rgph",
          "verificado": false,
          "seccoes": [
            { "slug": "<slug>", "nome": "<nome oficial>", "fonte_id": "ine-rgph", "verificado": false }
          ]
        }
      ]
    }
  ]
}
```

Os valores entre `<>` são ilustrativos; nenhum dado real é definido neste plano (Princípio I).

## Entidades

### Fonte

| Campo | Tipo | Obrigatório | Regra |
|-------|------|-------------|-------|
| `id` | string | sim | formato de slug; único entre fontes |
| `nome` | string | sim | não vazio (ex.: "INE Guiné-Bissau") |
| `referencia` | string | sim | documento citado (título, ano) |
| `url` | string | não | URL absoluto `http(s)://` quando presente |
| `consultado_em` | string | sim | data ISO `AAAA-MM-DD` |

### Região (entidade de primeiro nível)

| Campo | Tipo | Obrigatório | Regra |
|-------|------|-------------|-------|
| `slug` | string | sim | `^[a-z0-9]+(-[a-z0-9]+)*$`; único entre regiões |
| `nome` | string | sim | nome oficial com acentos; não vazio |
| `tipo` | enum | sim | `regiao` ou `setor-autonomo` |
| `capital` | string | sim | não vazio |
| `capital_fonte_id` | string | não | `Fonte.id` existente; fonte da capital quando diferente de `fonte_id` |
| `fonte_id` | string | sim | referencia um `Fonte.id` existente |
| `verificado` | bool | sim | `false` = pendente de verificação |
| `setores` | lista | sim | pode ser vazia |

O Setor Autónomo de Bissau é uma Região com `tipo: "setor-autonomo"`; as suas subdivisões seguem
a fonte oficial (lista de setores pode ficar vazia).

### Setor

| Campo | Tipo | Obrigatório | Regra |
|-------|------|-------------|-------|
| `slug` | string | sim | formato de slug; único entre **todos** os setores |
| `nome` | string | sim | não vazio |
| `fonte_id` | string | sim | `Fonte.id` existente |
| `verificado` | bool | sim | |
| `seccoes` | lista | sim | pode ser vazia |

Relação: pertence a exatamente uma Região (pelo aninhamento).

### Secção

| Campo | Tipo | Obrigatório | Regra |
|-------|------|-------------|-------|
| `slug` | string | sim | formato de slug; único entre **todas** as secções |
| `nome` | string | sim | não vazio |
| `fonte_id` | string | sim | `Fonte.id` existente |
| `verificado` | bool | sim | |

Relação: pertence a exatamente um Setor (pelo aninhamento).

### Informações do projeto (não está no ficheiro de dados)

Constantes em `internal/sobre` + versão injetada no build:
`nome`, `autor`, `licenca` (`MIT`), `versao`, `versao_api` (`v1`), `repositorio`, `fontes`
(lista de Fonte, lida do ficheiro de dados).

## Regras de validação (executadas no arranque e nos testes)

1. JSON bem formado e sem campos desconhecidos (`DisallowUnknownFields`).
2. Todos os campos obrigatórios presentes e não vazios.
3. Slugs no formato e únicos por tipo; ids de fonte únicos.
4. Toda referência `fonte_id` aponta para uma `Fonte.id` existente; toda fonte é usada por ≥ 1 entidade.
5. `tipo` ∈ {`regiao`, `setor-autonomo`}; exatamente uma entidade `setor-autonomo`.
6. Nenhuma lista contém nomes duplicados no mesmo pai (após normalização).
7. Órfãos impossíveis por construção (aninhamento); o teste confirma que todo setor indexado tem
   região e toda secção tem setor.

Falha em qualquer regra → o servidor não arranca e o gerador não publica.

## Estruturas em memória

- `regioes []Regiao` na ordem do ficheiro (ordem de listagem).
- `porSlugRegiao map[string]*Regiao`.
- `porSlugSetor map[string]*Setor` + referência à região do setor.
- `localidades []Localidade` para pesquisa: `{tipo, slug, nome, nome_normalizado, regiao, setor}`.

Ordem das listas: a do ficheiro (definida pelos curadores, por norma alfabética pelo nome).

## Estados

Uma entidade tem dois estados: **pendente** (`verificado: false`) → **verificado**
(`verificado: true`), transição feita por PR que cita o documento da fonte. Mudar o `slug` de uma
entidade publicada é uma mudança incompatível (Princípio VII) e só ocorre numa nova versão da API.
A estabilidade é garantida pelo teste `TestSlugsPublicadosEstaveis` e pela lista
`internal/dados/testdata/slugs-publicados.txt`.
