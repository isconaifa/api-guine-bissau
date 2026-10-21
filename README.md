# API Guiné-Bissau

API aberta, gratuita e só de leitura com a **divisão administrativa da Guiné-Bissau**: regiões e
Setor Autónomo de Bissau → setores → secções.

Hoje não existe uma fonte digital simples, aberta e confiável com estes dados, e quem desenvolve
sistemas, sites e aplicativos para a Guiné-Bissau acaba por digitar listas à mão, com erros e
inconsistências. Este projeto quer resolver isso, com **todos os dados com fonte citada**.

> ⚠️ **Estado dos dados** (outubro de 2026):
> - **Regiões e setores**: nomes do Instituto Nacional de Estatística da Guiné-Bissau (INE),
>   *Guiné-Bissau em Números 2017*, Tabela 1. Os 38 setores estão `"verificado": true`.
> - **Capitais regionais**: vêm da Wikipédia (fonte secundária). Por isso as regiões estão
>   `"verificado": false` até confirmação oficial.
> - **Secções**: ainda não publicadas (as listas de secções vêm vazias). Contribuições com fonte
>   oficial são bem-vindas.
>
> Consulte `fonte_id`, `verificado` e `/v1/sobre`.

Repositório: <https://github.com/isconaifa/api-guine-bissau>

## Endpoints

Todas as rotas são `GET`, públicas e sem autenticação, e começam com `/v1`. O contrato completo
está em [`docs/openapi.yaml`](docs/openapi.yaml) (OpenAPI 3.1).

| Rota | Descrição |
|------|-----------|
| `/v1/regioes` | Lista as 8 regiões e o Sector Autónomo de Bissau |
| `/v1/regioes/{regiao}` | Detalhe de uma região, com a lista dos seus setores |
| `/v1/regioes/{regiao}/setores` | Setores de uma região |
| `/v1/regioes/{regiao}/setores/{setor}` | Detalhe de um setor, com a região a que pertence |
| `/v1/regioes/{regiao}/setores/{setor}/seccoes` | Secções de um setor (lista vazia se ainda não houver) |
| `/v1/pesquisa?q=texto` | Pesquisa regiões, setores e secções pelo nome |
| `/v1/sobre` | Autor, licença, versão e fontes dos dados |
| `/v1/saude` | Verificação de saúde (para plataformas de hospedagem) |

### Formato das respostas

Todas as respostas usam o mesmo envelope JSON:

- **Lista**: `{"dados": [...], "total": n, "fonte": "..."}`
- **Objeto**: `{"dados": {...}, "fonte": "..."}`
- **Erro**: `{"erro": {"codigo": "...", "mensagem": "..."}}`

`fonte` é o nome legível das fontes dos dados devolvidos. Cada entidade traz ainda:

- `fonte_id`: id da fonte (detalhes em `/v1/sobre`);
- `verificado`: `false` enquanto o dado estiver pendente de verificação com a fonte oficial;
- `capital_fonte_id` (só regiões): fonte da capital, quando diferente de `fonte_id`.

Cada entidade tem um `slug` estável, sem acentos, em minúsculas e com hífens (ex.: `bafata`,
`sector-autonomo-de-bissau`). O nome oficial, com acentos, está em `nome`.

### Regiões

```bash
curl https://<servidor>/v1/regioes
```

```json
{"dados":[{"slug":"bafata","nome":"Bafatá","tipo":"regiao","capital":"Bafatá","capital_fonte_id":"wikipedia-pt-subdivisoes","fonte_id":"ine-gb-em-numeros-2017","verificado":false}, ...],"total":9,"fonte":"INE Guiné-Bissau; Wikipédia (pt)"}
```

`tipo` é `regiao` ou `setor-autonomo` (Sector Autónomo de Bissau).

```bash
curl https://<servidor>/v1/regioes/bafata
```

```json
{"dados":{"slug":"bafata","nome":"Bafatá","tipo":"regiao","capital":"Bafatá","capital_fonte_id":"wikipedia-pt-subdivisoes","fonte_id":"ine-gb-em-numeros-2017","verificado":false,"setores":[{"slug":"bafata","nome":"Bafatá"},{"slug":"bambadinca","nome":"Bambadinca"}, ...]},"fonte":"INE Guiné-Bissau; Wikipédia (pt)"}
```

```bash
curl https://<servidor>/v1/regioes/bafata/setores
```

```json
{"dados":[{"slug":"bafata","nome":"Bafatá","fonte_id":"ine-gb-em-numeros-2017","verificado":true}, ...],"total":6,"fonte":"INE Guiné-Bissau"}
```

### Setores e secções

```bash
curl https://<servidor>/v1/regioes/bafata/setores/bambadinca
```

```json
{"dados":{"slug":"bambadinca","nome":"Bambadinca","fonte_id":"ine-gb-em-numeros-2017","verificado":true,"regiao":{"slug":"bafata","nome":"Bafatá"},"total_seccoes":0},"fonte":"INE Guiné-Bissau"}
```

```bash
curl https://<servidor>/v1/regioes/bafata/setores/bambadinca/seccoes
```

```json
{"dados":[],"total":0,"fonte":"INE Guiné-Bissau"}
```

Um setor só é encontrado dentro da sua região: `/v1/regioes/oio/setores/bambadinca` devolve
`404 SETOR_NAO_ENCONTRADO`.

### Pesquisa

```bash
curl "https://<servidor>/v1/pesquisa?q=gabu"     # igual a q=GABÚ ou q=Gabú
```

```json
{"dados":[{"tipo":"regiao","slug":"gabu","nome":"Gabú","caminho":"/v1/regioes/gabu","fonte_id":"ine-gb-em-numeros-2017","verificado":false},{"tipo":"setor","slug":"gabu","nome":"Gabú","regiao":{"slug":"gabu","nome":"Gabú"},"caminho":"/v1/regioes/gabu/setores/gabu","fonte_id":"ine-gb-em-numeros-2017","verificado":true}],"total":2,"fonte":"INE Guiné-Bissau"}
```

- **Normalização**: acentos removidos, minúsculas, pontuação vira espaço
  (`"Bolama/Bijagós"` → `"bolama bijagos"`).
- **Correspondência**: o nome contém o texto pesquisado (ex.: `q=bafa` encontra "Bafatá").
- **Ordenação**: correspondência exata, depois prefixo, depois "contém"; em empate, regiões antes
  de setores antes de secções; depois por nome.
- `q` é obrigatório, com 1 a 100 caracteres; caso contrário: `400 PARAMETRO_INVALIDO`.
- Sem resultados: `200` com `"dados": []` e `"total": 0`.
- `caminho` é a rota da entidade (para secções, a lista de secções do setor).

### Sobre o projeto

```bash
curl https://<servidor>/v1/sobre
```

```json
{"dados":{"nome":"API Guiné-Bissau","autor":"Francisco Naifa","licenca":"MIT","versao":"v1.0.0","versao_api":"v1","repositorio":"https://github.com/isconaifa/api-guine-bissau","fontes":[{"id":"ine-gb-em-numeros-2017","nome":"INE Guiné-Bissau","referencia":"...","url":"...","consultado_em":"2026-10-21"}, ...]}}
```

### Erros

| Código | HTTP | Quando |
|--------|------|--------|
| `REGIAO_NAO_ENCONTRADA` | 404 | Não existe região com o slug pedido |
| `SETOR_NAO_ENCONTRADO` | 404 | O setor não existe ou não pertence à região do caminho |
| `PARAMETRO_INVALIDO` | 400 | `q` ausente, vazio ou com mais de 100 caracteres |
| `ROTA_NAO_ENCONTRADA` | 404 | Rota ou método inexistente |
| `ERRO_INTERNO` | 500 | Falha inesperada no servidor |

```json
{"erro":{"codigo":"REGIAO_NAO_ENCONTRADA","mensagem":"Região 'xyz' não encontrada."}}
```

Se o slug pedido tiver maiúsculas ou acentos (ex.: `/v1/regioes/Bafatá`), a mensagem indica o
formato esperado. Uma barra final (`/v1/regioes/`) redireciona com `301` para a rota sem barra.

## Versão estática (GitHub Pages)

Os mesmos dados são publicados como ficheiros JSON estáticos, sem servidor, em
<https://isconaifa.github.io/api-guine-bissau/>. Cada ficheiro tem exatamente o mesmo conteúdo da
resposta da API (um teste garante isso).

| Rota da API | Ficheiro estático |
|-------------|-------------------|
| `/v1/regioes` | `/v1/regioes.json` |
| `/v1/regioes/{r}` | `/v1/regioes/{r}.json` |
| `/v1/regioes/{r}/setores` | `/v1/regioes/{r}/setores.json` |
| `/v1/regioes/{r}/setores/{s}` | `/v1/regioes/{r}/setores/{s}.json` |
| `/v1/regioes/{r}/setores/{s}/seccoes` | `/v1/regioes/{r}/setores/{s}/seccoes.json` |
| `/v1/sobre` | `/v1/sobre.json` |
| `/v1/pesquisa?q=` | sem equivalente: use `/v1/pesquisa/indice.json` |

Para pesquisar sem servidor, descarregue `/v1/pesquisa/indice.json`, normalize o texto como
descrito em [Pesquisa](#pesquisa) e filtre os itens cujo `nome_normalizado` contém esse texto.

Limitações: os URLs terminam em `.json` e um caminho inexistente devolve a página 404 do GitHub
Pages (não um erro JSON).

O workflow [`pages.yml`](.github/workflows/pages.yml) gera e publica os ficheiros a cada push em
`main`. Para gerar localmente: `go run ./cmd/gerador -saida publico`.

## Executar localmente

Requer Go 1.26 ou mais recente.

```bash
go run ./cmd/api              # servidor em http://localhost:8080
PORT=9000 go run ./cmd/api    # outra porta
go test ./...                 # testes (endpoints, dados, pesquisa e paridade do gerador)
go test ./internal/handlers -run XXX -bench .   # benchmarks
```

Desempenho medido (binário de produção, Linux): arranque até à primeira resposta em ~40 ms,
~12 MB de memória residente e ~10 µs por pedido em `/v1/regioes` e `/v1/pesquisa` — bem abaixo de
50 ms por pedido.

## Hospedagem

A API é um binário único com os dados embutidos, pensado para planos gratuitos.

```bash
CGO_ENABLED=0 go build -trimpath \
  -ldflags "-s -w -X github.com/isconaifa/api-guine-bissau/internal/sobre.Versao=v1.0.0" \
  -o api ./cmd/api
PORT=8080 ./api
```

Ou com Docker (imagem final distroless, só com o binário):

```bash
docker build --build-arg VERSAO=v1.0.0 -t api-guine-bissau .
docker run -p 8080:8080 api-guine-bissau
```

- **Porta**: variável de ambiente `PORT` (padrão `8080`), definida automaticamente pelas
  plataformas.
- **Health check**: `GET /v1/saude`.
- **Render** (Web Service gratuito): build `go build -o api ./cmd/api`, start `./api`.
- **Koyeb** ou **Google Cloud Run**: usar o `Dockerfile`.
- Nos planos gratuitos, o serviço adormece sem tráfego e o primeiro pedido pode demorar alguns
  segundos (arranque a frio da plataforma, não da API). Para respostas sempre imediatas, use a
  [versão estática](#versão-estática-github-pages).

## Fontes dos dados

Os dados vivem em [`dados/divisao-administrativa.json`](dados/divisao-administrativa.json). Cada
entidade indica a sua fonte (`fonte_id`) e se já foi verificada (`verificado`). A lista completa
de fontes, com referência e data de consulta, está no próprio ficheiro e em `GET /v1/sobre`.

| Id | Fonte | Usada para |
|----|-------|------------|
| `ine-gb-em-numeros-2017` | INE Guiné-Bissau, *Guiné-Bissau em Números 2017*, Tabela 1 – Divisão administrativa (p. 9) | Nomes das regiões e dos setores |
| `wikipedia-pt-subdivisoes` | Wikipédia (pt), «Subdivisões da Guiné-Bissau», revisão 72530219 | Capitais regionais (provisório) |

O documento do INE publica os nomes em maiúsculas; aqui foram convertidos para a capitalização
normal, mantendo a grafia e os acentos (ex.: `GÃMAMUDO` → "Gãmamudo", `S. DOMINGOS` →
"S. Domingos").

Nenhum dado é inventado: só entram dados com fonte identificada.

## Contribuir

Contribuições são bem-vindas via pull request.

1. Faça um fork e crie um branch.
2. Para alterações de dados, **cite a fonte** (documento, ano, página/URL) na descrição do PR.
3. Garanta que tudo passa: `gofmt -l .` (sem saída), `go vet ./...` e `go test ./...`.
4. Slugs já publicados não mudam dentro de uma versão da API.

As regras do projeto estão na [constituição](.specify/memory/constitution.md).

## Autor

**Francisco Naifa** — Engenheiro de Computação pela Universidade Federal de Mato Grosso (UFMT).
GitHub: [@isconaifa](https://github.com/isconaifa)

## Licença

[MIT](LICENSE) © 2026 Francisco Naifa
