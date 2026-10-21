# Quickstart: validar a funcionalidade

**Feature**: `001-divisao-administrativa`

Guia de validação ponta a ponta. Formatos em [contracts/openapi.yaml](contracts/openapi.yaml);
regras dos dados em [data-model.md](data-model.md).

## Pré-requisitos

- Go 1.26+
- `curl` (e opcionalmente `jq`)

## 1. Qualidade e testes

```bash
gofmt -l .          # não deve listar ficheiros
go vet ./...
go test ./...       # endpoints, validação de dados, pesquisa e paridade do gerador
```

Esperado: tudo passa. Se `dados/divisao-administrativa.json` violar uma regra (slug duplicado,
fonte inexistente, campo em falta), os testes de dados falham com a entidade identificada.

## 2. Arrancar o servidor

```bash
go run ./cmd/api              # porta 8080
PORT=9000 go run ./cmd/api    # porta alternativa
```

Esperado: arranque em menos de 1 s. Com dados inválidos, o processo termina com erro e não abre
a porta.

## 3. Cenários (com o servidor em `localhost:8080`)

| # | Pedido | Esperado | Spec |
|---|--------|----------|------|
| 1 | `curl -s localhost:8080/v1/saude` | `200`, `{"dados":{"estado":"ok"}}` | — |
| 2 | `curl -s localhost:8080/v1/regioes` | `200`, `total` = nº de entidades de 1.º nível (9 previstas), cada uma com `slug`, `nome`, `tipo`, `capital`, `fonte_id`, `verificado` | US1-1 |
| 3 | `curl -s localhost:8080/v1/regioes/bafata` | `200`, `nome` = "Bafatá", lista `setores` | US1-2 |
| 4 | `curl -s localhost:8080/v1/regioes/bafata/setores` | `200`, lista de setores com `total` | US1-3 |
| 5 | `.../v1/regioes/<r>/setores/<s>` (um setor existente) | `200`, inclui `regiao.slug` = `<r>` | US2-1 |
| 6 | `.../v1/regioes/<r>/setores/<s>/seccoes` | `200`, lista (possivelmente vazia) | US2-2/3 |
| 7 | `curl -s 'localhost:8080/v1/pesquisa?q=gabu'` vs `?q=GAB%C3%9A` | resultados idênticos, incluem a região Gabú | US3-1/2 |
| 8 | `curl -s 'localhost:8080/v1/pesquisa?q=zzzz'` | `200`, `dados: []`, `total: 0` | US3-5 |
| 9 | `curl -s 'localhost:8080/v1/pesquisa?q='` | `400`, `erro.codigo` = `PARAMETRO_INVALIDO` | US5-3 |
| 10 | `curl -s localhost:8080/v1/regioes/bafta` | `404`, `REGIAO_NAO_ENCONTRADA`, mensagem em português | US5-1 |
| 11 | `.../v1/regioes/bafata/setores/nao-existe` | `404`, `SETOR_NAO_ENCONTRADO` | US5-2 |
| 12 | setor de outra região, ex. `.../v1/regioes/<outra>/setores/<s>` | `404`, `SETOR_NAO_ENCONTRADO` | US5-2 |
| 13 | `curl -s localhost:8080/v1/inexistente` | `404`, `ROTA_NAO_ENCONTRADA` | FR-010 |
| 14 | `curl -s localhost:8080/v1/sobre` | `200`, `autor`, `licenca: "MIT"`, `versao`, `fontes` | US4-1 |
| 15 | `curl -si localhost:8080/v1/regioes \| grep -i access-control` | `Access-Control-Allow-Origin: *` | FR-014 |

Hierarquia completa (SC-001): repetir 3 → 5 → 6 para cada região devolvida em 2.

## 4. Gerar e inspecionar o JSON estático

```bash
go run ./cmd/gerador -saida publico
ls publico/v1                         # regioes.json, regioes/, sobre.json, pesquisa/
diff <(curl -s localhost:8080/v1/regioes) publico/v1/regioes.json   # sem diferenças
```

Mapeamento completo em [contracts/json-estatico.md](contracts/json-estatico.md).

## 5. Binário de produção

```bash
CGO_ENABLED=0 go build -ldflags "-s -w -X <módulo>/internal/sobre.Versao=v1.0.0" -o api ./cmd/api
./api
```

Esperado: um único ficheiro executável (sem ficheiros de dados ao lado), memória residente
< 30 MB em repouso, e `/v1/sobre` mostra `versao: "v1.0.0"`.
