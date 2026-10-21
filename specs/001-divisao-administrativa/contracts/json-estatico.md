# Contrato: JSON estático (GitHub Pages)

Gerado por `go run ./cmd/gerador -saida publico` e publicado pelo CI em cada push para `main`.
O corpo de cada ficheiro é **byte a byte igual** à resposta `200` da rota equivalente da API
(ver [openapi.yaml](openapi.yaml)); um teste garante essa paridade.

| Rota da API | Ficheiro estático |
|-------------|-------------------|
| `GET /v1/regioes` | `/v1/regioes.json` |
| `GET /v1/regioes/{r}` | `/v1/regioes/{r}.json` |
| `GET /v1/regioes/{r}/setores` | `/v1/regioes/{r}/setores.json` |
| `GET /v1/regioes/{r}/setores/{s}` | `/v1/regioes/{r}/setores/{s}.json` |
| `GET /v1/regioes/{r}/setores/{s}/seccoes` | `/v1/regioes/{r}/setores/{s}/seccoes.json` |
| `GET /v1/sobre` | `/v1/sobre.json` |
| `GET /v1/pesquisa?q=` | sem equivalente — usar `/v1/pesquisa/indice.json` |

## Índice de pesquisa

`/v1/pesquisa/indice.json` usa o envelope de lista; cada item é um `ResultadoPesquisa` com o campo
extra `nome_normalizado`. Para pesquisar no cliente, aplicar a mesma normalização ao texto
(remover acentos → minúsculas → não alfanuméricos viram espaço → colapsar espaços) e filtrar os
itens cujo `nome_normalizado` contém o texto.

## Diferenças em relação ao servidor

- Não há respostas de erro em JSON: um caminho inexistente devolve o 404 HTML do GitHub Pages.
- Os URLs terminam em `.json`.
