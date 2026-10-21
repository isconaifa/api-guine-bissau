# API Guiné-Bissau Constitution

## Core Principles

### I. Precisão dos Dados (INEGOCIÁVEL)

- Todo dado publicado pela API MUST ter uma fonte identificada (ex.: Instituto Nacional de
  Estatística da Guiné-Bissau – INE), registada junto aos dados.
- Nenhum dado é inventado, estimado ou completado por suposição.
- Dados cuja fonte ainda não foi confirmada MUST ficar explicitamente marcados como
  pendentes de verificação, e essa marcação MUST ser visível no JSON devolvido.
- Em caso de conflito entre simplicidade e precisão, a precisão prevalece.

**Racional**: a utilidade da API depende inteiramente da confiança nos seus dados; um dado
errado sem fonte é pior do que um dado ausente.

### II. Simplicidade

- Preferir sempre a solução mais simples que funcione (YAGNI).
- Sem banco de dados enquanto os dados forem estáticos; os dados vivem em ficheiros
  versionados no repositório.
- Cada dependência externa MUST ser justificada; preferir a biblioteca padrão de Go.

**Racional**: menos peças significa menos falhas, manutenção mais fácil e contribuições mais
acessíveis.

### III. Português como Língua do Domínio

- Endpoints, campos JSON, mensagens de erro e documentação MUST estar em português
  (ex.: `regioes`, `setores`, `seccoes`).
- O código Go segue as convenções da linguagem (gofmt, nomes exportados em MixedCaps);
  identificadores MAY estar em português quando forem termos do domínio.

**Racional**: a API descreve a Guiné-Bissau para quem usa o português; o vocabulário do
domínio deve ser o oficial do país.

### IV. Identificadores Estáveis e Amigáveis

- Cada entidade MUST ter um `slug` sem acentos, em minúsculas e com hífens
  (ex.: `bafata`, `sector-autonomo-de-bissau`).
- O nome oficial, com acentos, MUST ser mantido no campo `nome`.
- Slugs são estáveis: uma vez publicados numa versão da API, não mudam nessa versão.

**Racional**: URLs previsíveis e permanentes permitem que clientes guardem e partilhem
referências sem quebrar.

### V. Baixo Consumo de Recursos

- O projeto MUST correr em planos de hospedagem gratuitos: pouca memória, arranque rápido.
- A entrega é um binário único, com os dados embutidos (ex.: `embed`), sem serviços
  auxiliares obrigatórios.

**Racional**: manter o custo de hospedagem próximo de zero garante a sustentabilidade de um
projeto aberto e comunitário.

### VI. Testes Obrigatórios (INEGOCIÁVEL)

- Todo endpoint MUST ter testes automatizados cobrindo respostas de sucesso e de erro.
- Os dados MUST ser validados por testes (ex.: nenhum setor sem região, nenhum slug
  duplicado, slugs no formato correto, toda entrada com fonte ou marcação de pendente).
- Nenhuma alteração é integrada com testes a falhar.

**Racional**: os testes são a garantia executável dos Princípios I e IV.

### VII. API Versionada

- Todas as rotas MUST começar com `/v1` (ou a versão vigente).
- Mudanças incompatíveis (remover/renomear campos, rotas ou slugs, mudar tipos) só podem
  ser feitas numa nova versão da API.
- Adições compatíveis (novos campos, novos endpoints) MAY entrar na versão atual.

**Racional**: clientes existentes não devem quebrar por causa da evolução da API.

### VIII. Respostas Consistentes

- Todos os endpoints MUST devolver o mesmo formato JSON de envelope, incluindo erros.
- Erros usam códigos HTTP adequados e uma estrutura de erro única, com mensagem em
  português.

**Racional**: um formato único simplifica clientes e reduz ambiguidade.

### IX. Código Aberto e Autoria Clara

- O projeto é licenciado sob MIT.
- A autoria MUST estar registada no `LICENSE`, no `README` e no endpoint `/v1/sobre`.
- Contribuições são bem-vindas via pull request e seguem esta constituição.

**Racional**: transparência sobre origem e licença permite reutilização segura e reconhece
quem contribui.

### X. Documentação Acompanha o Código

- Todo endpoint novo ou alterado MUST ser documentado no `README` na mesma alteração que o
  introduz.
- Futuramente, a documentação será também publicada em OpenAPI/Swagger, mantida em sincronia
  com o código.

**Racional**: uma API sem documentação atualizada não é utilizável por terceiros.

## Restrições Técnicas

- Linguagem: Go, seguindo `gofmt` e `go vet`.
- Dados: ficheiros estáticos versionados no repositório, embutidos no binário; cada registo
  traz a sua fonte ou a marcação de pendente de verificação.
- Persistência: nenhum banco de dados enquanto os dados forem estáticos; introduzir um exige
  emenda ou justificação explícita no plano da funcionalidade.
- Implantação: binário único, configurável por variáveis de ambiente (ex.: porta), compatível
  com planos de hospedagem gratuitos.

## Fluxo de Desenvolvimento e Critérios de Qualidade

- Toda alteração entra via pull request.
- Antes de integrar, MUST passar: `go test ./...` (incluindo os testes de validação de dados),
  `go vet` e formatação `gofmt`.
- A revisão verifica: fonte dos dados novos, estabilidade de slugs, prefixo `/v1`, formato de
  resposta consistente e atualização do `README`.
- Alterações de dados MUST citar a fonte na descrição do pull request.

## Governance

- Esta constituição prevalece sobre quaisquer outras práticas do projeto.
- Emendas são feitas por pull request que altere este ficheiro, com justificação e, quando
  necessário, plano de migração.
- Versionamento da constituição segue SemVer:
  - MAJOR: remoção ou redefinição incompatível de princípios ou regras de governança.
  - MINOR: novo princípio ou secção, ou expansão material de orientação existente.
  - PATCH: clarificações, redação e correções sem mudança semântica.
- Todo plano de funcionalidade (`/speckit-plan`) MUST verificar conformidade com estes
  princípios; qualquer exceção MUST ser justificada na secção de complexidade do plano.

**Version**: 1.0.0 | **Ratified**: 2026-10-21 | **Last Amended**: 2026-10-21
