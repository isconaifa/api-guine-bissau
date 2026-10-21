# Feature Specification: Divisão Administrativa da Guiné-Bissau

**Feature Branch**: `001-divisao-administrativa`

**Created**: 2026-10-21

**Status**: Draft

**Input**: User description: "Fonte digital simples, aberta e confiável com a divisão administrativa
da Guiné-Bissau (regiões e Setor Autónomo de Bissau → setores → secções): listar e detalhar regiões,
listar e detalhar setores, listar secções, pesquisar localidades pelo nome, informações sobre o
projeto e erros claros para itens inexistentes."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Consultar regiões e os seus setores (Priority: P1)

Uma desenvolvedora está a criar um formulário de cadastro para uma ONG e precisa de uma lista
fiável das regiões do país. Consulta a lista de todas as regiões (incluindo o Setor Autónomo de
Bissau), escolhe uma e obtém os seus detalhes: nome oficial, identificador amigável (slug),
capital e a lista de setores que a compõem.

**Why this priority**: é o nível mais alto da hierarquia e o caso de uso mais comum (menus de
seleção de região). Sozinho já elimina a necessidade de digitar listas à mão.

**Independent Test**: consultar a lista de regiões e verificar que contém exatamente as entidades
de primeiro nível da fonte oficial, cada uma com nome, slug, capital e fonte; depois consultar uma
região pelo slug e verificar os seus setores.

**Acceptance Scenarios**:

1. **Given** os dados carregados, **When** o utilizador pede a lista de regiões, **Then** recebe
   todas as entidades de primeiro nível (8 regiões + Setor Autónomo de Bissau, sujeito a
   confirmação na fonte oficial), cada uma com nome oficial, slug e indicação da fonte.
2. **Given** a região "Bafatá" existe, **When** o utilizador pede os detalhes da região `bafata`,
   **Then** recebe o nome "Bafatá", o slug `bafata`, a capital e a lista dos seus setores.
3. **Given** a região "Bafatá" existe, **When** o utilizador pede apenas os setores de `bafata`,
   **Then** recebe a lista dos setores dessa região, cada um com nome e slug.

---

### User Story 2 - Consultar setores e as suas secções (Priority: P2)

Um investigador precisa de descer até ao nível mais detalhado. Consulta um setor e obtém os seus
detalhes, incluindo a região a que pertence, e depois a lista das secções desse setor.

**Why this priority**: completa a hierarquia região → setor → secção, que é o critério de sucesso
principal, mas depende de a navegação por regiões já existir para ser descoberta.

**Independent Test**: consultar um setor conhecido pelo slug e verificar que indica a região a que
pertence; consultar as secções desse setor e verificar que todas pertencem a ele.

**Acceptance Scenarios**:

1. **Given** um setor existe, **When** o utilizador pede os detalhes desse setor, **Then** recebe
   nome, slug, a região a que pertence (nome e slug) e a fonte dos dados.
2. **Given** um setor existe, **When** o utilizador pede as secções desse setor, **Then** recebe a
   lista das secções, cada uma com nome, slug e fonte.
3. **Given** um setor sem secções confirmadas na fonte, **When** o utilizador pede as suas
   secções, **Then** recebe uma lista vazia (não um erro).

---

### User Story 3 - Pesquisar localidades pelo nome (Priority: P3)

Um utilizador sabe o nome de uma localidade, mas não o seu nível nem a região. Escreve o nome (por
exemplo "gabu", "GABÚ" ou "Gabú") e recebe as localidades correspondentes, de qualquer nível, com
indicação do tipo e da sua posição na hierarquia.

**Why this priority**: muito útil para autocompletar e validação, mas a navegação hierárquica
(P1/P2) já permite obter todos os dados sem ela.

**Independent Test**: pesquisar o mesmo nome com e sem acentos e em maiúsculas/minúsculas e
verificar que os resultados são idênticos.

**Acceptance Scenarios**:

1. **Given** existe a região "Gabú", **When** o utilizador pesquisa "gabu", **Then** os
   resultados incluem a região "Gabú", identificada como região.
2. **Given** a mesma pesquisa, **When** o utilizador escreve "GABÚ" ou "Gabú", **Then** recebe
   exatamente os mesmos resultados.
3. **Given** uma pesquisa por parte do nome (ex.: "bafa"), **When** executada, **Then** os
   resultados incluem as localidades cujo nome contém esse texto.
4. **Given** um setor ou secção encontrado, **When** apresentado nos resultados, **Then** indica a
   que setor e/ou região pertence.
5. **Given** um texto sem correspondência, **When** pesquisado, **Then** recebe uma lista vazia.

---

### User Story 4 - Informações sobre o projeto (Priority: P4)

Uma organização quer saber se pode usar e citar os dados. Consulta as informações do projeto:
autor, licença, versão e fontes dos dados.

**Why this priority**: necessário para confiança e reutilização, mas não bloqueia o uso dos dados.

**Independent Test**: consultar as informações do projeto e verificar que autor, licença (MIT),
versão e fontes estão presentes e coincidem com o que está registado na documentação.

**Acceptance Scenarios**:

1. **Given** o serviço disponível, **When** o utilizador pede as informações do projeto,
   **Then** recebe autor, licença, versão e a lista das fontes dos dados.

---

### User Story 5 - Erros claros (Priority: P5)

Um desenvolvedor comete um erro de escrita no slug (ex.: `bafta`). Em vez de uma resposta vazia
ou confusa, recebe um erro claro, em português, a dizer que a região não existe.

**Why this priority**: transversal a todas as histórias; melhora a experiência mas não acrescenta
dados.

**Independent Test**: pedir uma região, um setor e as secções de um setor inexistentes e
verificar que cada pedido devolve um erro com o mesmo formato e uma mensagem compreensível.

**Acceptance Scenarios**:

1. **Given** não existe região com slug `bafta`, **When** o utilizador pede essa região,
   **Then** recebe um erro de "não encontrado" com mensagem em português que identifica o recurso
   em falta.
2. **Given** não existe setor com o slug pedido, **When** o utilizador pede o setor ou as suas
   secções, **Then** recebe o mesmo tipo de erro, no mesmo formato.
3. **Given** uma pesquisa sem texto, **When** executada, **Then** recebe um erro de pedido
   inválido, no mesmo formato.

---

### Edge Cases

- Slug com maiúsculas ou acentos (ex.: `Bafatá`): tratado como inexistente, com erro claro que
  indica o formato esperado; os slugs válidos são apenas os publicados.
- Nomes iguais em níveis diferentes (ex.: região e setor com o mesmo nome, como Bafatá): a
  pesquisa devolve ambos, distinguidos pelo tipo.
- Nomes iguais no mesmo nível em regiões diferentes: cada entidade mantém um slug único entre as
  entidades do mesmo tipo em todo o conjunto de dados (ex.: o setor desambigua-se com um sufixo).
  Entidades de tipos diferentes podem partilhar o slug (região `bafata` e setor `bafata`).
- Setor Autónomo de Bissau: aparece na lista de primeiro nível ao lado das regiões, distinguido
  pelo tipo.
- Região "Bolama/Bijagós" (nome com barra): o slug não contém a barra (ex.: `bolama-bijagos`).
- O nome oficial pode usar grafia anterior ao Acordo Ortográfico (ex.: "Sector"): o nome é
  mantido como na fonte e o slug acompanha o nome oficial (ex.: `sector-autonomo-de-bissau`).
- Dados ainda não confirmados pela fonte oficial: são devolvidos, mas claramente marcados como
  pendentes de verificação.
- Pesquisa com texto muito curto (1 caractere) ou só com espaços: espaços são ignorados; texto
  vazio após limpeza resulta em erro de pedido inválido.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: O sistema MUST disponibilizar a lista completa das entidades de primeiro nível
  (regiões e Setor Autónomo de Bissau), cada uma com nome oficial, slug, tipo e fonte.
- **FR-002**: O sistema MUST permitir consultar os detalhes de uma região pelo slug: nome, slug,
  tipo, capital, lista de setores (nome e slug) e fonte.
- **FR-003**: O sistema MUST permitir listar apenas os setores de uma região indicada pelo slug.
- **FR-004**: O sistema MUST permitir consultar os detalhes de um setor pelo slug, incluindo a
  região a que pertence (nome e slug) e a fonte.
- **FR-005**: O sistema MUST permitir listar as secções de um setor indicado pelo slug, cada uma
  com nome, slug e fonte.
- **FR-006**: O sistema MUST permitir pesquisar localidades (regiões, setores e secções) por
  nome ou parte do nome, ignorando acentos e diferenças entre maiúsculas e minúsculas.
- **FR-007**: Cada resultado da pesquisa MUST indicar o tipo da localidade e a sua posição na
  hierarquia (região e, quando aplicável, setor).
- **FR-008**: O sistema MUST disponibilizar as informações do projeto: autor, licença, versão e
  fontes dos dados.
- **FR-009**: Quando um recurso pedido não existe, o sistema MUST responder com um erro de "não
  encontrado" com mensagem em português que identifique o que não foi encontrado.
- **FR-010**: Todas as respostas, de sucesso ou de erro, MUST seguir o mesmo formato geral.
- **FR-011**: Cada entidade MUST ter um slug único entre as entidades do mesmo tipo em todo o
  conjunto de dados, sem acentos, em minúsculas e com hífens; o nome oficial com acentos MUST ser mantido separadamente.
- **FR-012**: Cada entidade MUST indicar a fonte dos seus dados; entidades não confirmadas pela
  fonte oficial MUST estar marcadas como pendentes de verificação.
- **FR-013**: Nenhum setor pode existir sem região, e nenhuma secção sem setor.
- **FR-014**: O acesso a todas as consultas MUST ser público, sem necessidade de registo ou
  autenticação.
- **FR-015**: Os nomes de recursos, campos e mensagens MUST estar em português.
- **FR-016**: Todas as consultas MUST pertencer à primeira versão identificada do serviço, de
  forma que mudanças incompatíveis futuras não afetem quem a utiliza.

### Key Entities *(include if feature involves data)*

- **Região**: entidade de primeiro nível da divisão administrativa. Atributos: nome oficial,
  slug, tipo (região ou setor autónomo), capital, fonte, estado de verificação. Contém setores.
- **Setor Autónomo de Bissau**: entidade de primeiro nível com estatuto próprio, apresentada ao
  lado das regiões e distinguida pelo tipo. A sua subdivisão segue a fonte oficial.
- **Setor**: subdivisão de uma região. Atributos: nome oficial, slug, região a que pertence,
  fonte, estado de verificação. Contém secções.
- **Secção**: subdivisão de um setor. Atributos: nome oficial, slug, setor a que pertence, fonte,
  estado de verificação.
- **Fonte**: origem identificada de um dado (ex.: Instituto Nacional de Estatística da
  Guiné-Bissau – INE). Atributos: nome, referência/documento, data de consulta.
- **Informações do projeto**: autor, licença, versão e lista de fontes utilizadas.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Qualquer pessoa consegue obter a hierarquia completa região → setor → secção de
  qualquer região percorrendo no máximo 3 consultas por ramo, sem documentação adicional além do
  README.
- **SC-002**: 100% das entidades publicadas têm fonte citada ou marcação de pendente de
  verificação.
- **SC-003**: 100% das entidades têm slug único e no formato definido; nenhum setor ou secção
  fica sem entidade superior.
- **SC-004**: 95% das consultas respondem em menos de 1 segundo, mesmo em hospedagem gratuita.
- **SC-005**: Uma pesquisa por um mesmo nome escrita com ou sem acentos e em qualquer combinação
  de maiúsculas/minúsculas devolve sempre resultados idênticos.
- **SC-006**: 100% dos pedidos por recursos inexistentes recebem um erro no formato comum, com
  mensagem em português.

## Assumptions

- A lista de regiões indicada (Bafatá, Biombo, Bolama/Bijagós, Cacheu, Gabú, Oio, Quinara,
  Tombali e Setor Autónomo de Bissau) é provisória e será confirmada com a fonte oficial (INE);
  até lá, os itens ficam marcados como pendentes de verificação.
- O Setor Autónomo de Bissau é tratado como entidade de primeiro nível, listado junto às regiões
  e distinguido pelo tipo; os seus setores/secções seguem o que a fonte oficial indicar.
- Os dados são estáticos e mudam raramente; atualizações acontecem com novas versões do projeto.
- O volume de dados é pequeno (dezenas de setores, centenas de secções), pelo que as listas são
  devolvidas completas, sem paginação.
- Não há consulta individual de detalhes de uma secção na primeira versão; as secções são
  obtidas pela lista do setor ou pela pesquisa.
- O acesso é só de leitura; não há criação, edição ou remoção de dados pelos utilizadores.
- A pesquisa encontra correspondências por "contém" (parte do nome), não por semelhança
  aproximada (erros de ortografia).
- Fora do âmbito desta versão: coordenadas e limites geográficos, feriados nacionais, validação
  de números de telefone (+245), dicionário Kriol ↔ Português e estatísticas de população.
