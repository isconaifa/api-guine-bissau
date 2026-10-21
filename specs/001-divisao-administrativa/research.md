# Research: Divisão Administrativa da Guiné-Bissau

**Feature**: `001-divisao-administrativa` | **Date**: 2026-10-21

Decisões técnicas que resolvem as incógnitas do Technical Context. Decisões tomadas pelo autor do
projeto estão indicadas como tal; as restantes são escolhas deste plano.

## R1. Linguagem e versão

- **Decision**: Go 1.26 (`go 1.26` no `go.mod`), versão estável instalada no ambiente.
- **Rationale**: decisão do autor; binário único estático, arranque em milissegundos, pouca
  memória (Princípio V).
- **Alternatives considered**: nenhuma (decisão prévia).

## R2. Framework HTTP

- **Decision**: Gin (`github.com/gin-gonic/gin`), última versão estável v1.x fixada no `go.mod`,
  em modo `release`, criado com `gin.New()` + `gin.Recovery()` (sem o logger por omissão, para
  reduzir ruído e custo).
- **Rationale**: decisão do autor; roteamento com parâmetros, `NoRoute` para erros uniformes e
  helpers de teste conhecidos.
- **Alternatives considered**: `net/http` com o `ServeMux` de padrões (Go ≥ 1.22) cobriria todas
  as rotas sem dependências. Rejeitado por decisão do autor; o desvio ao Princípio II está
  registado em *Complexity Tracking* no plano.

## R3. Armazenamento e carregamento dos dados

- **Decision**: um único ficheiro `dados/divisao-administrativa.json` (hierárquico: fontes +
  regiões → setores → secções) embutido com `//go:embed` num pequeno pacote Go dentro de
  `dados/`. No arranque, `internal/dados` faz o *parse*, valida e constrói índices em memória
  (mapas por slug, lista plana para pesquisa). Se a validação falhar, o processo termina com erro
  (falha rápida) — dados inválidos nunca são servidos.
- **Rationale**: Princípio II (sem banco de dados) e V (binário único). A estrutura aninhada
  garante por construção que nenhum setor fica sem região e nenhuma secção sem setor.
- **Alternatives considered**: ficheiros separados por nível com chaves estrangeiras (mais
  propenso a órfãos); CSV (perde a hierarquia e é pior para revisão em PR).
- **Nota**: `go:embed` não aceita caminhos com `..`, por isso o ficheiro `dados/embed.go`
  (`package dados`) expõe os bytes; `internal/dados` importa-o com alias (`brutos`).

## R4. Normalização de texto e pesquisa

- **Decision**: normalização = NFD (`golang.org/x/text/unicode/norm`) → remover marcas
  combinatórias (`unicode.Mn`) → minúsculas → substituir tudo o que não seja letra/dígito por
  espaço → colapsar espaços → *trim*. Ex.: `"Bolama/Bijagós"` → `"bolama bijagos"`.
  Correspondência por "contém" sobre o nome normalizado, pré-calculado no arranque.
  Ordenação: correspondência exata → prefixo → contém; depois por tipo (região, setor, secção);
  depois por nome normalizado.
- **Rationale**: cobre acentos do português (á, â, ã, ç, é, ê, í, ó, ô, õ, ú) e de nomes Kriol;
  `golang.org/x/text` já é dependência transitiva do Gin, logo não acrescenta peso.
- **Alternatives considered**: tabela manual de substituição (frágil); pesquisa aproximada
  (fuzzy) — fora do âmbito (spec, Assumptions).
- **Validação de `q`**: obrigatório; vazio após normalização ou com mais de 100 caracteres →
  `400 PARAMETRO_INVALIDO`.

## R5. Formato das respostas

- **Decision**: envelope único (Princípio VIII):
  - Sucesso: `{"dados": <objeto|lista>, "total": <n, só em listas>, "fonte": "<texto>"}`.
    `fonte` resume as fontes das entidades devolvidas (nomes distintos separados por `"; "`);
    é omitido em respostas sem dados territoriais (`/v1/sobre`, `/v1/saude`).
  - Erro: `{"erro": {"codigo": "<CODIGO>", "mensagem": "<texto em português>"}}`.
  - Cada entidade traz também `fonte_id` (id da fonte) e `verificado` (booleano), cumprindo
    FR-012 ao nível do registo. O `fonte` do envelope é o nome legível; nomes distintos evitam
    que o mesmo campo tenha dois significados.
- **Rationale**: formato proposto pelo autor, estendido para o detalhe (objeto) e para a
  marcação por entidade exigida pelo Princípio I.
- **Alternatives considered**: `total` em todas as respostas (sem significado num detalhe).

## R6. Códigos de erro

| Código | HTTP | Quando |
|--------|------|--------|
| `REGIAO_NAO_ENCONTRADA` | 404 | slug de região inexistente |
| `SETOR_NAO_ENCONTRADO` | 404 | setor inexistente ou que não pertence à região do caminho |
| `PARAMETRO_INVALIDO` | 400 | `q` ausente, vazio ou longo demais |
| `ROTA_NAO_ENCONTRADA` | 404 | qualquer rota/método não definido (`NoRoute`) |
| `ERRO_INTERNO` | 500 | pânico recuperado |

Se o valor pedido não tiver formato de slug (ex.: `Bafatá`), a mensagem acrescenta uma dica:
"Use o slug em minúsculas e sem acentos (ex.: 'bafata')."

## R7. Slugs

- **Decision**: slugs escritos à mão no JSON (revistos em PR), validados por teste com
  `^[a-z0-9]+(-[a-z0-9]+)*$` e unicidade **por tipo** em todo o conjunto (região, setor, secção).
  Em colisão dentro do mesmo tipo, acrescenta-se o slug do pai (ex.: `<seccao>-<setor>`).
- **Rationale**: as rotas distinguem o tipo pelo caminho, logo região `bafata` e setor `bafata`
  podem coexistir sem ambiguidade; slugs estáveis não devem depender de um algoritmo que mude.
  Clarificação aplicada à spec (FR-011 e Edge Cases).
- **Alternatives considered**: unicidade global entre tipos (forçaria `bafata-setor`, pouco
  amigável); gerar slugs automaticamente (risco de mudança involuntária).

## R7a. Grafia ("sector" vs "setor")

- **Decision**: o `nome` de cada entidade é **exatamente** o da fonte oficial (se o INE escreve
  "Sector Autónomo de Bissau", é esse o nome) e o `slug` deriva desse nome (ex.:
  `sector-autonomo-de-bissau`). O vocabulário da própria API (rotas, campos e valores de `tipo`)
  segue a grafia do Acordo Ortográfico usada na spec: `setores`, `seccoes`,
  `tipo: "setor-autonomo"`.
- **Rationale**: Princípio I (nome oficial sem alterações) e Princípio IV (slug derivado do nome
  oficial); o vocabulário da API mantém-se estável e coerente, independentemente da grafia de
  cada fonte.
- **Alternatives considered**: uniformizar tudo para "setor" (alteraria nomes oficiais);
  uniformizar tudo para "sector" (contrariaria o vocabulário definido na spec e na constituição).

## R8. Publicação estática (GitHub Pages)

- **Decision**: `cmd/gerador` reutiliza os mesmos construtores de resposta da API e escreve em
  `publico/` (ignorado pelo git): `v1/regioes.json`, `v1/regioes/<r>.json`,
  `v1/regioes/<r>/setores.json`, `v1/regioes/<r>/setores/<s>.json`,
  `v1/regioes/<r>/setores/<s>/seccoes.json`, `v1/sobre.json` e
  `v1/pesquisa/indice.json` (lista plana com `nome_normalizado` para pesquisa no cliente).
  Ver [contracts/json-estatico.md](contracts/json-estatico.md).
- **Rationale**: GitHub Pages só serve ficheiros; extensão `.json` garante o `Content-Type`
  correto. Partilhar o código de construção garante respostas idênticas (testado).
- **Alternatives considered**: pastas com `index.json` (o Pages não as resolve como índice).
- **Limitação**: não há pesquisa no servidor nem erros JSON no modo estático (404 do Pages).

## R9. CORS e cache

- **Decision**: middleware próprio (sem dependência) que define
  `Access-Control-Allow-Origin: *` e `Cache-Control: public, max-age=3600` nas rotas de dados
  `/v1` (exceto `/v1/saude`, que usa `no-store`).
  O GitHub Pages já envia CORS aberto.
- **Rationale**: a API é pública e de leitura (FR-014); clientes em browser precisam de CORS.

## R10. Configuração, saúde e versão

- **Decision**: `PORT` (padrão `8080`). `GET /v1/saude` → `200 {"dados":{"estado":"ok"}}` com
  `Cache-Control: no-store`.
  Versão do projeto injetada com `-ldflags "-X <módulo>/internal/sobre.Versao=vX.Y.Z"`
  (padrão `dev`); autor, licença e URL do repositório em constantes em `internal/sobre`.
- **Rationale**: as plataformas gratuitas (Render, Koyeb, Cloud Run) definem `PORT` e precisam de
  um *health check* simples. Fica dentro de `/v1`, como todas as rotas (Princípio VII).

## R11. Testes e CI

- **Decision**: `testing` + `net/http/httptest` com `gin.SetMode(gin.TestMode)`:
  - testes de cada endpoint (sucesso, 404, 400, rota inexistente, formato do envelope);
  - testes de validação de dados (slugs, unicidade, fontes referenciadas existem, campos
    obrigatórios, região com capital);
  - testes unitários da normalização e ordenação da pesquisa;
  - teste de paridade: saída do gerador == corpo das respostas da API.
  GitHub Actions: `gofmt -l`, `go vet ./...`, `go test ./...` em cada PR; em `main`, gerar
  `publico/` e publicar com `actions/upload-pages-artifact` + `actions/deploy-pages`.
- **Rationale**: Princípio VI e fluxo de qualidade da constituição.

## R12. Origem dos dados (bloqueante para publicação, não para implementação)

- **Decision**: os dados só são preenchidos a partir de documento identificado (prioridade:
  INE Guiné-Bissau — Recenseamento Geral da População e Habitação). Até confirmação, cada
  registo entra com `verificado: false`. Nenhum setor/secção é acrescentado sem fonte.
- **Rationale**: Princípio I (INEGOCIÁVEL). O plano não inclui dados reais; a recolha é uma
  tarefa própria.

## R13. Autoria

- **Decision**: autor **Francisco Naifa**, Engenheiro de Computação pela Universidade Federal de
  Mato Grosso (UFMT). Registado em `LICENSE` (`Copyright (c) 2026 Francisco Naifa`), no `README`
  (secção Autor, com a formação e o GitHub `@isconaifa`) e em `/v1/sobre` (`autor: "Francisco Naifa"`,
  `repositorio: "https://github.com/isconaifa/api-guine-bissau"`), conforme o
  Princípio IX.

## Itens em aberto (não bloqueiam o plano)

- Documento oficial exato (título, ano, URL) da fonte INE.
