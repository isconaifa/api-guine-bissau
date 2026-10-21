---

description: "Lista de tarefas da funcionalidade Divisão Administrativa da Guiné-Bissau"
---

# Tasks: Divisão Administrativa da Guiné-Bissau

**Input**: Design documents from `/specs/001-divisao-administrativa/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/openapi.yaml,
contracts/json-estatico.md, quickstart.md

**Tests**: OBRIGATÓRIOS — a constituição (Princípio VI, INEGOCIÁVEL) exige testes para todo
endpoint e validação dos dados por testes. Em cada fase, os testes são escritos primeiro e devem
falhar antes da implementação.

**Organization**: tarefas agrupadas por história de utilizador (US1–US5 da spec), para que cada
uma seja implementada e testada de forma independente.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: pode correr em paralelo (ficheiros diferentes, sem dependências pendentes)
- **[Story]**: história a que a tarefa pertence (US1…US5)
- Caminhos relativos à raiz do repositório (`/home/isco/Documentos/PROJETOS PESSOAIS/API Guiné-Bissau`)

## Convenções para todas as tarefas

- Código Go formatado com `gofmt`; testes ao lado do código (`*_test.go`), pacote `testing` +
  `net/http/httptest`; nos testes de handlers usar `gin.SetMode(gin.TestMode)`.
- Campos JSON, rotas, códigos de erro e mensagens em português (Princípio III).
- `<módulo>` = `github.com/isconaifa/api-guine-bissau` (definido em T001).
- Nenhum dado territorial é inventado (Princípio I): só entram dados com fonte citada; dados não
  confirmados pelo INE levam `"verificado": false`.

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: inicialização do projeto Go e ficheiros de base

- [ ] T001 Inicializar o módulo em `go.mod` na raiz com `go mod init github.com/isconaifa/api-guine-bissau` e diretiva `go 1.26`
- [ ] T002 Adicionar dependências com `go get github.com/gin-gonic/gin@latest golang.org/x/text@latest` e fixar versões em `go.mod`/`go.sum`
- [ ] T003 [P] Criar `.gitignore` com `publico/`, `/api`, `/gerador`, `*.exe`, `.idea/workspace.xml`, `coverage.out`
- [ ] T004 [P] Criar `LICENSE` com o texto MIT padrão e a linha `Copyright (c) 2026 Francisco Naifa`
- [ ] T005 [P] Criar `.github/workflows/ci.yml` (em `push` e `pull_request`): `actions/checkout`, `actions/setup-go` com `go-version-file: go.mod`, passos `test -z "$(gofmt -l .)"`, `go vet ./...`, `go test ./...`
- [ ] T006 [P] Criar `README.md` inicial: descrição do projeto (problema da spec), secção "Executar localmente" (`go run ./cmd/api`, `PORT`), secção "Fontes dos dados", secção "Contribuir" (PR, citar fonte na descrição, `go test ./...`), secção "Autor" (Francisco Naifa, Engenheiro de Computação pela Universidade Federal de Mato Grosso – UFMT, GitHub [@isconaifa](https://github.com/isconaifa)), link do repositório `https://github.com/isconaifa/api-guine-bissau` e secção "Licença" (MIT); deixar secção "Endpoints" vazia para as histórias preencherem

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: normalização de texto, modelo/carregamento/validação dos dados, envelope de respostas,
router base e servidor — necessários a todas as histórias

**⚠️ CRITICAL**: nenhuma história começa antes desta fase estar concluída

### Normalização de texto

- [ ] T007 [P] Escrever testes de tabela em `internal/pesquisa/normalizar_test.go` para `Normalizar(s string) string`: `"Bafatá"`→`"bafata"`, `"GABÚ"`→`"gabu"`, `"Bolama/Bijagós"`→`"bolama bijagos"`, `"  Setor   Autónomo de Bissau "`→`"setor autonomo de bissau"`, `"Quinhámel"`→`"quinhamel"`, `"ção"`→`"cao"`, `""`→`""`, `"!!!"`→`""`
- [ ] T008 Implementar `Normalizar` em `internal/pesquisa/normalizar.go`: NFD (`golang.org/x/text/unicode/norm`) → remover runas `unicode.Mn` → `strings.ToLower` → todo caractere que não seja letra/dígito vira espaço → colapsar espaços → `strings.TrimSpace` (research R4)

### Modelo, carregamento e validação dos dados

- [ ] T009 [P] Definir tipos em `internal/dados/modelo.go` conforme data-model.md: `Fonte{ID "id", Nome "nome", Referencia "referencia", URL "url,omitempty", ConsultadoEm "consultado_em"}`, `Regiao{Slug, Nome, Tipo, Capital, FonteID, Verificado, Setores []Setor}`, `Setor{Slug, Nome, FonteID, Verificado, Seccoes []Seccao}`, `Seccao{Slug, Nome, FonteID, Verificado}`, `Arquivo{Fontes []Fonte, Regioes []Regiao}` com tags JSON em português (`slug`, `nome`, `tipo`, `capital`, `fonte_id`, `verificado`, `setores`, `seccoes`); constantes `TipoRegiao = "regiao"`, `TipoSetorAutonomo = "setor-autonomo"`
- [ ] T010 [P] Criar fixture válida `internal/dados/testdata/valido.json` com 1 fonte (`id: "fonte-teste"`), 2 regiões fictícias claramente de teste (ex.: `regiao-a` tipo `regiao` com 2 setores, um com 2 secções e outro com 0; `setor-autonomo-teste` tipo `setor-autonomo` sem setores), mais uma região e um setor com o mesmo slug para provar unicidade por tipo
- [ ] T011 Escrever testes em `internal/dados/validacao_test.go` que carregam variações inválidas (inline ou em `internal/dados/testdata/`) e verificam que `Carregar` devolve erro que identifica a entidade, para cada regra do data-model.md: (1) JSON mal formado e campo desconhecido (`DisallowUnknownFields`); (2) campo obrigatório vazio/ausente — `nome`, `slug`, `capital`, `fonte_id`, `referencia`, `consultado_em`; (3) slug fora de `^[a-z0-9]+(-[a-z0-9]+)*$` (ex.: `Bafatá`, `bafata-`, `ba--fata`) e slug duplicado no mesmo tipo (regiões; setores em regiões diferentes; secções em setores diferentes), id de fonte duplicado; (4) `fonte_id` que referencia id inexistente e fonte não usada por nenhuma entidade; (5) `tipo` fora de {`regiao`, `setor-autonomo`} e número de entidades `setor-autonomo` diferente de exatamente uma; (6) nomes duplicados no mesmo pai após `pesquisa.Normalizar`; `consultado_em` fora de `AAAA-MM-DD`; `url` presente que não começa por `http://` ou `https://`; e que `valido.json` carrega sem erro
- [ ] T012 Implementar `Carregar(b []byte) (*Base, error)` em `internal/dados/carregar.go` e as regras em `internal/dados/validar.go` (decoder com `DisallowUnknownFields`; todas as regras de T011; mensagens em português, ex.: `setor "x" na região "y": fonte "z" inexistente`)
- [ ] T013 Implementar índices em `internal/dados/base.go`: `type Base` com `Fontes() []Fonte`, `Fonte(id) (Fonte, bool)`, `Regioes() []Regiao` (ordem do ficheiro), `Regiao(slug) (*Regiao, bool)`, `Setor(slugRegiao, slugSetor) (*Setor, *Regiao, bool)` (falha se o setor não pertencer à região), e `Localidades() []Localidade` com `Localidade{Tipo, Slug, Nome, NomeNormalizado, Regiao *Ref, Setor *Ref, FonteID, Verificado}` (Tipo ∈ `regiao`, `setor-autonomo`, `setor`, `seccao`); testes em `internal/dados/base_test.go` com `valido.json` (incl. setor procurado na região errada → não encontrado; toda localidade setor/secção tem região e secção tem setor)

### Dados reais embutidos

- [ ] T014 Criar `dados/embed.go` (`package dados`, `//go:embed divisao-administrativa.json`, `var DivisaoAdministrativa []byte`)
- [ ] T015 Criar `dados/divisao-administrativa.json` com as 9 entidades de 1.º nível da spec — Bafatá (`bafata`), Biombo (`biombo`), Bolama/Bijagós (`bolama-bijagos`), Cacheu (`cacheu`), Gabú (`gabu`), Oio (`oio`), Quinara (`quinara`), Tombali (`tombali`) com `tipo: "regiao"` e o Sector/Setor Autónomo de Bissau com `tipo: "setor-autonomo"` e nome e slug conforme a grafia da fonte (research R7a) — provisoriamente `nome: "Sector Autónomo de Bissau"`, `slug: "sector-autonomo-de-bissau"`, todas com `"verificado": false` e `setores: []`; a `capital` de cada uma e a(s) fonte(s) DEVEM vir de um documento público identificável (preferir publicações do INE Guiné-Bissau; se só houver fonte secundária, citá-la em `fontes` com `referencia`, `url` e `consultado_em`); NÃO acrescentar setores/secções sem documento que os liste; se uma capital não puder ser confirmada por nenhuma fonte, parar e perguntar ao autor
- [ ] T016 Escrever `internal/dados/dados_reais_test.go`: `Carregar(brutos.DivisaoAdministrativa)` (import de `<módulo>/dados` com alias `brutos`) sem erro; exatamente 9 entidades de 1.º nível; existe exatamente uma entidade com `tipo: "setor-autonomo"` (sem fixar o slug no teste, ver research R7a); todas as entidades têm `fonte_id` válida

- [ ] T017 Criar `internal/dados/testdata/slugs-publicados.txt` (uma linha por entidade: `<tipo>:<caminho de slugs>`, ex.: `regiao:bafata`, `setor:bafata/<s>`, `seccao:bafata/<s>/<sc>`) e o teste `TestSlugsPublicadosEstaveis` em `internal/dados/dados_reais_test.go`: falha se algum slug da lista não existir em `brutos.DivisaoAdministrativa`; slugs novos são permitidos, mas o teste avisa (via `t.Log`) para os acrescentar à lista. Remover ou renomear uma linha da lista só é permitido numa nova versão da API (Princípios IV e VII)

### Metadados do projeto

- [ ] T018 [P] Criar `internal/sobre/sobre.go`: `var Versao = "dev"` (injetada por `-ldflags "-X <módulo>/internal/sobre.Versao=vX.Y.Z"`), constantes `Nome = "API Guiné-Bissau"`, `Autor = "Francisco Naifa"`, `Licenca = "MIT"`, `VersaoAPI = "v1"`, `Repositorio = "https://github.com/isconaifa/api-guine-bissau"`

### Envelope de respostas e erros

- [ ] T019 [P] Escrever testes em `internal/respostas/respostas_test.go`: serialização de `Lista` → `{"dados":[...],"total":n,"fonte":"..."}` (lista vazia serializa `[]`, nunca `null`), `Objeto` → `{"dados":{...},"fonte":"..."}` (sem `total`), `Objeto` sem fonte omite `fonte`, `Erro` → `{"erro":{"codigo":"...","mensagem":"..."}}`; `ResumoFontes` junta nomes distintos por `"; "` na ordem de primeira ocorrência; `Serializar` não escapa HTML (`"Bafatá"` e `&` aparecem literais), não indenta e termina com `\n`
- [ ] T020 Implementar `internal/respostas/respostas.go`: tipos `Lista{Dados any "dados"; Total int "total"; Fonte string "fonte"}`, `Objeto{Dados any "dados"; Fonte string "fonte,omitempty"}`, `Erro{Erro DetalheErro "erro"}`; constantes `CodigoRegiaoNaoEncontrada = "REGIAO_NAO_ENCONTRADA"`, `CodigoSetorNaoEncontrado = "SETOR_NAO_ENCONTRADO"`, `CodigoParametroInvalido = "PARAMETRO_INVALIDO"`, `CodigoRotaNaoEncontrada = "ROTA_NAO_ENCONTRADA"`, `CodigoErroInterno = "ERRO_INTERNO"`; função `ResumoFontes(base *dados.Base, ids ...string) string`; função `Serializar(v any) ([]byte, error)` (`json.Encoder` com `SetEscapeHTML(false)`, sem indentação, `\n` final) — a única forma de produzir corpos JSON no projeto (API e gerador); este pacote NÃO importa Gin
- [ ] T021 Implementar construtores de erro em `internal/respostas/erros.go`: `RegiaoNaoEncontrada(slug)` → `"Região 'xyz' não encontrada."`, `SetorNaoEncontrado(regiao, setor)` → `"Setor 'x' não encontrado na região 'y'."`, `ParametroInvalido(msg)`, `RotaNaoEncontrada(caminho)`, `ErroInterno()`; quando o slug pedido não casa `^[a-z0-9]+(-[a-z0-9]+)*$`, acrescentar `" Use o slug em minúsculas e sem acentos (ex.: 'bafata')."` (research R6); testes em `internal/respostas/erros_test.go`

### Router base e servidor

- [ ] T022 Escrever testes em `internal/handlers/router_test.go` (com `valido.json`): `GET /v1/saude` → 200 `{"dados":{"estado":"ok"}}` com `Cache-Control: no-store`; `GET /saude` → 404 `ROTA_NAO_ENCONTRADA`; `GET /v1/inexistente` → 404 `ROTA_NAO_ENCONTRADA`; `POST /v1/regioes` → 404 `ROTA_NAO_ENCONTRADA`; rotas de dados `/v1` têm `Access-Control-Allow-Origin: *` e `Cache-Control: public, max-age=3600`; `Content-Type` `application/json; charset=utf-8`; pânico num handler de teste → 500 `ERRO_INTERNO`
- [ ] T023 Implementar `NovoRouter(base *dados.Base) *gin.Engine` em `internal/handlers/router.go`: `gin.New()`, `gin.CustomRecovery` que responde `respostas.ErroInterno()` com 500, `NoRoute` e `NoMethod` com `ROTA_NAO_ENCONTRADA`, grupo `/v1` com o middleware de `internal/handlers/middleware.go` (CORS `*` e `Cache-Control: public, max-age=3600`), rota `GET /v1/saude` no grupo `/v1` (o handler define `Cache-Control: no-store`); helper `responder(c *gin.Context, status int, v any)` em `internal/handlers/responder.go` que usa `respostas.Serializar` e `c.Data(status, "application/json; charset=utf-8", b)` — proibido usar `c.JSON`; `NoRoute`, `NoMethod`, `Recovery` e todos os erros passam por `responder`
- [ ] T024 Implementar `cmd/api/main.go`: `gin.SetMode(gin.ReleaseMode)`, `dados.Carregar(brutos.DivisaoAdministrativa)` (em erro: `log.Fatal` com a mensagem — não abre a porta), porta de `PORT` (padrão `8080`), `http.Server` com `ReadHeaderTimeout` de 5 s, log de uma linha no arranque

**Checkpoint**: `go test ./...` passa; `go run ./cmd/api` arranca e `/v1/saude` responde

---

## Phase 2b: Recolha de dados (bloqueia a validação real de US2 e de SC-001)

**Purpose**: preencher os dados reais com fonte oficial (Princípio I). Pode correr em paralelo com
as fases de histórias; as histórias são testadas com a fixture, mas SC-001 só se cumpre aqui.

- [ ] T025 Obter o documento oficial do INE Guiné-Bissau com a lista de regiões, setores e secções (título, ano, URL); se não estiver disponível publicamente, perguntar ao autor; registar em `fontes` de `dados/divisao-administrativa.json` (`id`, `nome`, `referencia`, `url`, `consultado_em`)
- [ ] T026 Confirmar as 9 entidades de 1.º nível e as capitais com o documento de T025: passar `verificado` para `true` só no que o documento confirma; ajustar nomes oficiais e slugs à grafia da fonte (research R7a) antes da primeira publicação
- [ ] T027 Para cada região (um PR por região, fonte citada na descrição do PR), acrescentar os setores a `dados/divisao-administrativa.json` com `nome`, `slug`, `fonte_id`, `verificado` conforme o documento
- [ ] T028 Para cada setor, acrescentar as secções a `dados/divisao-administrativa.json` conforme o documento e acrescentar os novos slugs a `internal/dados/testdata/slugs-publicados.txt`

**Checkpoint**: hierarquia real completa; cenários de US2 e SC-001 verificáveis com dados reais

---

## Phase 3: User Story 1 - Consultar regiões e os seus setores (Priority: P1) 🎯 MVP

**Goal**: listar as entidades de 1.º nível, ver o detalhe de uma região e listar os seus setores

**Independent Test**: `GET /v1/regioes` devolve as 9 entidades com `slug`, `nome`, `tipo`,
`capital`, `fonte_id`, `verificado`; `GET /v1/regioes/bafata` e `/v1/regioes/bafata/setores`
respondem conforme `contracts/openapi.yaml`

### Tests for User Story 1 ⚠️

- [ ] T029 [P] [US1] Testes dos construtores em `internal/respostas/regioes_test.go` (com `valido.json`): `ListaRegioes` (total, ordem do ficheiro, campos `slug`, `nome`, `tipo`, `capital`, `fonte_id`, `verificado`, sem `setores`), `DetalheRegiao` (inclui `setores` como `[{slug, nome}]`, `[]` quando vazio), `SetoresDaRegiao` (itens `slug`, `nome`, `fonte_id`, `verificado`); `fonte` do envelope = `ResumoFontes`
- [ ] T030 [P] [US1] Testes HTTP em `internal/handlers/regioes_test.go`: `GET /v1/regioes` 200; `GET /v1/regioes/regiao-a` 200; `GET /v1/regioes/regiao-a/setores` 200 com `total` correto; `GET /v1/regioes/xyz` e `/v1/regioes/xyz/setores` → 404 `REGIAO_NAO_ENCONTRADA` com `"Região 'xyz' não encontrada."`; `GET /v1/regioes/Bafat%C3%A1` → 404 com a dica de formato

### Implementation for User Story 1

- [ ] T031 [US1] Implementar em `internal/respostas/regioes.go` os tipos `RegiaoResumo`, `RegiaoDetalhe`, `SetorResumo`, `Referencia{Slug "slug", Nome "nome"}` e as funções `ListaRegioes(base) Lista`, `DetalheRegiao(base, *dados.Regiao) Objeto`, `SetoresDaRegiao(base, *dados.Regiao) Lista`, conforme os schemas `RegiaoResumo`, `RegiaoDetalhe` e `SetorResumo` de `contracts/openapi.yaml`
- [ ] T032 [US1] Implementar handlers em `internal/handlers/regioes.go` e registar em `NovoRouter`: `GET /v1/regioes`, `GET /v1/regioes/:regiao`, `GET /v1/regioes/:regiao/setores` (404 via `respostas.RegiaoNaoEncontrada`); responder sempre via `responder(...)`
- [ ] T033 [US1] Documentar os três endpoints na secção "Endpoints" do `README.md` com exemplo de pedido `curl` e resposta (Princípio X)

**Checkpoint**: US1 funcional e testada de forma independente (MVP)

---

## Phase 4: User Story 2 - Consultar setores e as suas secções (Priority: P2)

**Goal**: detalhe de um setor com a sua região e lista das secções do setor

**Independent Test**: `GET /v1/regioes/<r>/setores/<s>` inclui `regiao.slug = <r>`;
`.../seccoes` devolve lista (possivelmente vazia); setor de outra região → 404

### Tests for User Story 2 ⚠️

- [ ] T034 [P] [US2] Testes dos construtores em `internal/respostas/setores_test.go`: `DetalheSetor` com `slug`, `nome`, `fonte_id`, `verificado`, `regiao{slug,nome}` e `total_seccoes`; `SeccoesDoSetor` com itens `slug`, `nome`, `fonte_id`, `verificado`, e `dados: []`, `total: 0` para setor sem secções
- [ ] T035 [P] [US2] Testes HTTP em `internal/handlers/setores_test.go` (com `valido.json`): detalhe e secções de um setor existente → 200; setor sem secções → 200 com lista vazia; setor inexistente → 404 `SETOR_NAO_ENCONTRADO`; setor existente pedido sob outra região → 404 `SETOR_NAO_ENCONTRADO`; região inexistente → 404 `REGIAO_NAO_ENCONTRADA`

### Implementation for User Story 2

- [ ] T036 [US2] Implementar em `internal/respostas/setores.go` os tipos `SetorDetalhe`, `Seccao` e as funções `DetalheSetor(base, *dados.Setor, *dados.Regiao) Objeto` e `SeccoesDoSetor(base, *dados.Setor) Lista`, conforme os schemas `SetorDetalhe` e `Seccao` de `contracts/openapi.yaml`
- [ ] T037 [US2] Implementar handlers em `internal/handlers/setores.go` e registar: `GET /v1/regioes/:regiao/setores/:setor` e `GET /v1/regioes/:regiao/setores/:setor/seccoes` (verificar primeiro a região → `REGIAO_NAO_ENCONTRADA`, depois `base.Setor(r, s)` → `SETOR_NAO_ENCONTRADO`); responder sempre via `responder(...)`
- [ ] T038 [US2] Documentar os dois endpoints no `README.md` com exemplos

**Checkpoint**: US1 e US2 funcionam; hierarquia completa navegável (SC-001)

---

## Phase 5: User Story 3 - Pesquisar localidades pelo nome (Priority: P3)

**Goal**: `GET /v1/pesquisa?q=` encontra regiões, setores e secções ignorando acentos e maiúsculas

**Independent Test**: `q=gabu`, `q=GABÚ` e `q=Gabú` devolvem resultados idênticos;
`q=zzzz` → lista vazia; `q=` → 400

### Tests for User Story 3 ⚠️

- [ ] T039 [P] [US3] Testes em `internal/pesquisa/pesquisa_test.go` para `Pesquisar(itens []Item, q string) []int` (com `Item{NomeNormalizado string; Nivel int}`, `Nivel` 0 = região/setor autónomo, 1 = setor, 2 = secção; devolve os índices dos itens encontrados): correspondência por "contém" no `NomeNormalizado` de `Normalizar(q)`; ordenação exata → prefixo → contém, depois `Nivel`, depois `NomeNormalizado`; mesmo resultado para variações de acentos/maiúsculas; sem correspondência → slice vazio não-nil
- [ ] T040 [P] [US3] Testes HTTP em `internal/handlers/pesquisa_test.go`: resultados iguais para `q=regiao-a`/`q=REGIÃO A`; itens com `tipo`, `slug`, `nome`, `fonte_id`, `verificado`, `caminho`, e `regiao`/`setor` só quando aplicável; `caminho` = `/v1/regioes/<r>` (região), `/v1/regioes/<r>/setores/<s>` (setor), `/v1/regioes/<r>/setores/<s>/seccoes` (secção); `q` ausente, `q=`, `q=%20%20` e `q` com mais de 100 caracteres → 400 `PARAMETRO_INVALIDO`; sem correspondência → 200 `dados: []`, `total: 0`

### Implementation for User Story 3

- [ ] T041 [US3] Implementar `Item` e `Pesquisar` em `internal/pesquisa/pesquisa.go` (research R4); `internal/pesquisa` NÃO importa `internal/dados` (que já importa `internal/pesquisa` para `Normalizar`), por isso opera só sobre `[]Item`; acrescentar a `internal/dados/base.go` o método `Pesquisar(q string) []Localidade`, que mantém um `[]pesquisa.Item` paralelo a `Localidades()` construído no arranque e converte os índices devolvidos, com teste em `internal/dados/base_test.go`
- [ ] T042 [US3] Implementar `ResultadoPesquisa` e `ResultadosPesquisa(base, []dados.Localidade) Lista` em `internal/respostas/pesquisa.go`, conforme o schema `ResultadoPesquisa` de `contracts/openapi.yaml` (incl. `caminho`)
- [ ] T043 [US3] Implementar handler `GET /v1/pesquisa` em `internal/handlers/pesquisa.go`: validar `q` (obrigatório; `len([]rune(q)) <= 100`; vazio após `Normalizar` → `PARAMETRO_INVALIDO` com mensagem `"O parâmetro 'q' é obrigatório e deve ter entre 1 e 100 caracteres."`), chamar `base.Pesquisar(q)` e responder `respostas.ResultadosPesquisa` via `responder(...)`; registar em `NovoRouter`
- [ ] T044 [US3] Documentar `/v1/pesquisa` no `README.md` (regras de normalização e ordenação, exemplos)

**Checkpoint**: pesquisa funcional e independente da navegação

---

## Phase 6: User Story 4 - Informações sobre o projeto (Priority: P4)

**Goal**: `GET /v1/sobre` com autor, licença, versão e fontes

**Independent Test**: resposta contém `autor: "Francisco Naifa"`, `licenca: "MIT"`, `versao`,
`versao_api: "v1"`, `repositorio` e `fontes` iguais às do ficheiro de dados

### Tests for User Story 4 ⚠️

- [ ] T045 [P] [US4] Testes em `internal/handlers/sobre_test.go`: 200; envelope `{"dados":{...}}` sem `total` nem `fonte`; campos `nome`, `autor`, `licenca`, `versao`, `versao_api`, `repositorio`, `fontes` (cada fonte com `id`, `nome`, `referencia`, `consultado_em` e `url` só se presente)

### Implementation for User Story 4

- [ ] T046 [US4] Implementar `Sobre(base) Objeto` em `internal/respostas/sobre.go` (schema `Sobre` de `contracts/openapi.yaml`, valores de `internal/sobre`)
- [ ] T047 [US4] Implementar handler `GET /v1/sobre` em `internal/handlers/sobre.go` (via `responder(...)`) e registar em `NovoRouter`
- [ ] T048 [US4] Documentar `/v1/sobre` no `README.md`

**Checkpoint**: metadados do projeto disponíveis (Princípio IX)

---

## Phase 7: User Story 5 - Erros claros (Priority: P5)

**Goal**: todos os erros, em todas as rotas, seguem o mesmo formato com mensagem em português

**Independent Test**: pedidos inexistentes/inválidos em cada rota devolvem
`{"erro":{"codigo","mensagem"}}` com o código e o status corretos

### Tests for User Story 5 ⚠️

- [ ] T049 [P] [US5] Teste de conformidade em `internal/handlers/erros_test.go`: tabela com um caso de erro por rota (`/v1/regioes/xyz`, `/v1/regioes/xyz/setores`, `/v1/regioes/regiao-a/setores/xyz`, `/v1/regioes/regiao-a/setores/xyz/seccoes`, `/v1/pesquisa`, `/v1/nada`, `/nada`) verificando status, `Content-Type` JSON, que o corpo tem só a chave `erro` com `codigo` ∈ enum de `contracts/openapi.yaml` e `mensagem` não vazia; e que respostas de sucesso nunca contêm `erro`

### Implementation for User Story 5

- [ ] T050 [US5] Corrigir qualquer divergência encontrada por T049 em `internal/handlers/*.go` / `internal/respostas/erros.go` e garantir que Gin não emite respostas em texto simples (ex.: redirecionamento de barra final: `RedirectTrailingSlash` mantido, mas confirmar que `/v1/regioes/` resulta em 301 para `/v1/regioes` ou em erro JSON)
- [ ] T051 [US5] Documentar no `README.md` a secção "Erros" com a tabela de códigos (`REGIAO_NAO_ENCONTRADA` 404, `SETOR_NAO_ENCONTRADO` 404, `PARAMETRO_INVALIDO` 400, `ROTA_NAO_ENCONTRADA` 404, `ERRO_INTERNO` 500) e um exemplo de corpo

**Checkpoint**: todas as histórias funcionais

---

## Phase 8: Publicação estática (GitHub Pages)

**Purpose**: gerar JSON estático idêntico às respostas da API (contracts/json-estatico.md);
depende de US1, US2, US3 e US4

- [ ] T052 Escrever teste de paridade em `cmd/gerador/gerador_test.go`: com `valido.json`, gerar num `t.TempDir()` e comparar byte a byte cada ficheiro com o corpo da resposta `200` da rota equivalente servida por `handlers.NovoRouter` via `httptest`; verificar que existem `v1/regioes.json`, `v1/regioes/<r>.json`, `v1/regioes/<r>/setores.json`, `v1/regioes/<r>/setores/<s>.json`, `v1/regioes/<r>/setores/<s>/seccoes.json`, `v1/sobre.json`, `v1/pesquisa/indice.json`
- [ ] T053 Implementar `Gerar(base *dados.Base, dir string) error` em `cmd/gerador/gerador.go` usando apenas `internal/respostas` e escrevendo cada ficheiro com `respostas.Serializar` (a mesma função usada pelos handlers), e `v1/pesquisa/indice.json` = envelope de lista com todos os `ResultadoPesquisa` + campo `nome_normalizado`
- [ ] T054 Implementar `cmd/gerador/main.go` com flag `-saida` (padrão `publico`), que carrega `brutos.DivisaoAdministrativa`, apaga e recria a pasta de saída e chama `Gerar`
- [ ] T055 [P] Criar `.github/workflows/pages.yml` (push em `main` + `workflow_dispatch`; permissões `pages: write`, `id-token: write`): `go test ./...`, `go run ./cmd/gerador -saida publico`, `actions/upload-pages-artifact` com `path: publico`, `actions/deploy-pages`
- [ ] T056 [P] Documentar no `README.md` a secção "Versão estática (GitHub Pages)": tabela rota → ficheiro de `contracts/json-estatico.md`, uso do índice de pesquisa no cliente e limitações (sem erros JSON, URLs com `.json`)

---

## Phase 9: Polish & Cross-Cutting Concerns

**Purpose**: hospedagem, desempenho, consistência final e pendências

- [ ] T057 [P] Criar `Dockerfile` multi-stage: `golang:1.26` com `CGO_ENABLED=0 go build -ldflags "-s -w -X <módulo>/internal/sobre.Versao=${VERSAO}" -o /api ./cmd/api` → imagem final `gcr.io/distroless/static-debian12:nonroot`, `EXPOSE 8080`, `ENTRYPOINT ["/api"]`; e `.dockerignore` (`publico/`, `.git`, `.idea`, `specs/`)
- [ ] T058 [P] Documentar no `README.md` a secção "Hospedagem": build do binário com `-ldflags`, variável `PORT`, health check `/v1/saude`, notas para Render, Koyeb e Cloud Run (arranque a frio dos planos gratuitos)
- [ ] T059 [P] Adicionar benchmarks em `internal/handlers/bench_test.go` (`BenchmarkListaRegioes`, `BenchmarkPesquisa`) e registar no `README.md` que cada pedido fica bem abaixo de 50 ms
- [ ] T060 Medir memória e arranque do binário de produção (quickstart.md §5): RSS em repouso < 30 MB e arranque < 1 s; registar resultado na descrição do PR
- [ ] T061 [P] Copiar `specs/001-divisao-administrativa/contracts/openapi.yaml` para `docs/openapi.yaml` e referenciá-lo no `README.md` (Princípio X)
- [ ] T062 Rever todas as mensagens, campos e nomes de rotas: português, sem inglês nos contratos públicos (Princípio III)
- [ ] T063 Correr `gofmt -l .`, `go vet ./...` e `go test ./...` e corrigir falhas
- [ ] T064 Executar todos os cenários de `specs/001-divisao-administrativa/quickstart.md` contra o servidor local e o gerador, e corrigir divergências

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: sem dependências
- **Foundational (Phase 2)**: depende do Setup — BLOQUEIA todas as histórias
- **US1–US5 (Phases 3–7)**: dependem só da Phase 2; podem seguir em paralelo ou por prioridade
- **Publicação estática (Phase 8)**: depende de US1, US2, US3 e US4 (gera os seus ficheiros)
- **Recolha de dados (Phase 2b)**: depende de T015–T016 e do documento oficial; corre em paralelo com as histórias
- **Polish (Phase 9)**: depende das fases desejadas

### Dependências dentro da Phase 2

- T008 depende de T007; T012 depende de T008, T009, T011; T013 depende de T012
- T015 depende de T014; T016 depende de T013 e T015
- T020 depende de T013 e T019; T021 depende de T020
- T023 depende de T021 e T022; T024 depende de T016 e T023

### User Story Dependencies

- **US1 (P1)**: após a Phase 2 — sem dependências de outras histórias
- **US2 (P2)**: após a Phase 2 — implementada e testada com a fixture; a validação com dados reais (e SC-001) depende da Phase 2b; independente (reutiliza `Referencia` de US1; se US2 for feita antes, criar `Referencia` em `internal/respostas/respostas.go`)
- **US3 (P3)**: após a Phase 2 — independente; `caminho` não exige que as rotas de US1/US2 existam
- **US4 (P4)**: após a Phase 2 — independente
- **US5 (P5)**: o teste T049 cobre as rotas de US1–US3; executar depois delas para cobertura total

### Within Each User Story

- Testes primeiro, a falhar → construtores em `internal/respostas` → handlers → README
- Uma história concluída (checkpoint) antes de passar à seguinte prioridade

### Parallel Opportunities

- Setup: T003, T004, T005, T006 em paralelo
- Foundational: T007, T009, T010, T018, T019 em paralelo; T014 em paralelo com T011–T013
- Após a Phase 2: US1, US2, US3 e US4 em paralelo (ficheiros distintos; só `router.go` é partilhado ao registar rotas — registar em sequência)
- Em cada história, os dois testes [P] em paralelo
- Polish: T057, T058, T059, T061 em paralelo

---

## Parallel Example: User Story 1

```bash
# Testes de US1 em paralelo:
Task: "Testes dos construtores em internal/respostas/regioes_test.go"
Task: "Testes HTTP em internal/handlers/regioes_test.go"

# Depois, em sequência:
Task: "Construtores em internal/respostas/regioes.go"
Task: "Handlers em internal/handlers/regioes.go"
Task: "Documentar endpoints no README.md"
```

## Parallel Example: após a Phase 2

```bash
Task: "US1 — T029..T033"
Task: "US3 — T039..T044"
Task: "US4 — T045..T048"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Phase 1: Setup
2. Phase 2: Foundational (inclui as 9 entidades com fonte e `verificado: false`)
3. Phase 3: US1
4. **STOP and VALIDATE**: cenários 1–4 e 10 do quickstart.md
5. Publicar o servidor num plano gratuito

### Incremental Delivery

1. Setup + Foundational → base pronta
2. US1 → MVP (lista de regiões)
3. US2 + Phase 2b → hierarquia completa com dados reais (SC-001)
4. US3 → pesquisa
5. US4 + US5 → metadados e erros uniformes
6. Phase 8 → versão estática no GitHub Pages
7. Polish → hospedagem documentada

---

## Notes

- [P] = ficheiros diferentes, sem dependências pendentes
- [Story] liga cada tarefa à história da spec, para rastreabilidade
- Verificar que os testes falham antes de implementar
- Fazer commit após cada tarefa ou grupo lógico
- Nunca acrescentar dados territoriais sem fonte (Princípio I); em dúvida, perguntar ao autor
