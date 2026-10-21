# Implementation Plan: Divisão Administrativa da Guiné-Bissau

**Branch**: `001-divisao-administrativa` | **Date**: 2026-10-21 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `/specs/001-divisao-administrativa/spec.md`

## Summary

API pública, só de leitura, com a divisão administrativa da Guiné-Bissau (regiões e Setor
Autónomo de Bissau → setores → secções), pesquisa por nome sem acentos e informações do projeto.
Os dados vivem num ficheiro JSON versionado, embutido no binário com `go:embed`, validado e
indexado em memória no arranque. Um servidor Gin expõe as rotas `/v1`; um gerador reutiliza os
mesmos construtores de resposta para publicar JSON estático idêntico no GitHub Pages.
Detalhes e justificações em [research.md](research.md).

## Technical Context

**Language/Version**: Go 1.26

**Primary Dependencies**: Gin v1.x (`github.com/gin-gonic/gin`); `golang.org/x/text` (já
transitiva do Gin) para normalização Unicode

**Storage**: ficheiro `dados/divisao-administrativa.json` embutido (`go:embed`), em memória; sem
banco de dados

**Testing**: `testing` + `net/http/httptest` (`gin.TestMode`)

**Target Platform**: Linux (contentor ou binário) em plano gratuito — Render, Koyeb ou Cloud
Run; GitHub Pages para a variante estática

**Project Type**: web-service (API REST só de leitura) + gerador de site estático

**Performance Goals**: p95 < 50 ms no servidor (excluindo arranque a frio da plataforma), para
cumprir SC-004 (< 1 s percebido)

**Constraints**: memória residente < 30 MB; arranque < 1 s; binário único sem ficheiros
externos; `PORT` por variável de ambiente (padrão 8080)

**Scale/Scope**: 9 entidades de 1.º nível, dezenas de setores, centenas de secções; 8 rotas

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| # | Princípio | Avaliação | Estado |
|---|-----------|-----------|--------|
| I | Precisão dos Dados | `fonte` + `verificado` por entidade; validação de fontes referenciadas; sem dados inventados no plano (R12) | ✅ |
| II | Simplicidade | Sem BD; um ficheiro de dados. **Gin é dependência dispensável** (stdlib bastaria) | ⚠️ justificado |
| III | Português | Rotas, campos, códigos e mensagens em português; Go idiomático | ✅ |
| IV | Slugs estáveis | Slugs manuais, formato e unicidade por tipo testados; `nome` com acentos (R7) | ✅ |
| V | Baixo consumo | Binário único com dados embutidos; < 30 MB; arranque < 1 s | ✅ |
| VI | Testes | Testes por endpoint, de dados, de pesquisa e de paridade do gerador (R11) | ✅ |
| VII | API versionada | Todas as rotas em `/v1`, incluindo o health check `/v1/saude` | ✅ |
| VIII | Respostas consistentes | Envelope único `dados`/`total`/`fonte` e `erro{codigo,mensagem}`, incl. `NoRoute` (R5, R6) | ✅ |
| IX | Código aberto | `LICENSE` MIT, autoria (Francisco Naifa) no LICENSE, README e `/v1/sobre` | ✅ |
| X | Documentação | README por endpoint; OpenAPI já produzido em `contracts/` | ✅ |

**Resultado (pré-research)**: PASS com 1 desvio justificado em *Complexity Tracking*.

**Re-check pós-design (Phase 1)**: PASS. O design não introduziu novas dependências nem
violações; o gerador estático partilha código com a API (sem duplicação) e a paridade é testada.

## Project Structure

### Documentation (this feature)

```text
specs/001-divisao-administrativa/
├── plan.md              # Este ficheiro
├── research.md          # Phase 0
├── data-model.md        # Phase 1
├── quickstart.md        # Phase 1
├── contracts/
│   ├── openapi.yaml     # Contrato HTTP v1
│   └── json-estatico.md # Contrato da publicação estática
└── tasks.md             # Phase 2 (/speckit-tasks)
```

### Source Code (repository root)

```text
.
├── cmd/
│   ├── api/                 # main: carrega dados, monta router Gin, lê PORT
│   └── gerador/             # main: escreve publico/v1/**.json para o GitHub Pages
├── internal/
│   ├── dados/               # parse, validação e índices em memória (importa dados/ via embed)
│   ├── pesquisa/            # Normalizar(), correspondência e ordenação
│   ├── respostas/           # envelope, erros e construtores partilhados por handlers e gerador
│   ├── handlers/            # handlers Gin, NoRoute, middleware CORS/cache, NovoRouter()
│   └── sobre/               # autor, licença, repositório, Versao (ldflags)
├── dados/
│   ├── divisao-administrativa.json
│   └── embed.go             # //go:embed → var DivisaoAdministrativa []byte
├── .github/workflows/
│   ├── ci.yml               # gofmt, go vet, go test em PRs e push
│   └── pages.yml            # gera publico/ e publica no GitHub Pages (main)
├── LICENSE
├── README.md
└── go.mod                   # module github.com/isconaifa/api-guine-bissau
```

Testes ficam ao lado do código (`*_test.go`), conforme a convenção Go: `internal/dados`
(validação dos dados reais), `internal/pesquisa`, `internal/handlers` (cada endpoint) e
`cmd/gerador` (paridade).

**Structure Decision**: a estrutura sugerida pelo autor, na raiz deste repositório, com dois
pacotes adicionais: `internal/respostas` (para o gerador produzir exatamente as mesmas respostas
sem depender do Gin) e `internal/sobre` (metadados do projeto e versão injetada no build).

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| Dependência Gin (Princípio II) | Decisão explícita do autor: familiaridade, `NoRoute`, ecossistema e convenções conhecidas | `net/http` com `ServeMux` (Go ≥ 1.22) cobriria as rotas sem dependências; rejeitado por preferência do autor. Mitigação: Gin confinado a `internal/handlers` e `cmd/api`; a lógica (`dados`, `pesquisa`, `respostas`) não depende dele, permitindo trocar depois |
