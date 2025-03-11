## Legendas

`GERAL` - Tarefas que competem a todo sistema.<br/>

`CAMPEONATO` - Tarefas que competem a gestao, consulta de campeonatos e API Externa 
- > Desenvolver uma API REST em Go que gerencie campeonatos esportivos e suas partidas, garantindo segurança com autenticação JWT.

<br/>

`TORCEDORES` - Tarefas que competem a gestao e notificacao de torcedores
- > Objetivo: Permitir que os usuários se cadastrem para receber notificações sobre os jogos do seu time.

<br/>
<br/>

## Tarefas

- [ ] `GERAL` Planning/Discovery do projeto (Esse documento e os respectivos items como `Tasks Kamban`) [ISSUE-0001](https://github.com/jtonynet/go-soccer-fan/issues/1)


- [ ] `GERAL` Definir o ADR inicial e de arquitetura com os requisitos da [ISSUE-0001](https://github.com/jtonynet/go-soccer-fan/issues/1)
    - Router: Gin,
    - DBLib: Gorm
    - DB: Postgres
    - TDD
    - Docker

<br/>

- [ ] `GERAL` Definir changelog com v0.0.0 [ISSUE-0001](https://github.com/jtonynet/go-soccer-fan/issues/1)

<br/>

- [ ] `GERAL` - Definir `DER` para campeonatos/partidas e torcedores [ISSUE-0002](https://github.com/jtonynet/go-soccer-fan/issues/2)

<br/>

- [ ] `CAMPEONATO` - `Repository`, `Service` a partir do `DER` com `TDD` `happy-path`  [ISSUE-0003](https://github.com/jtonynet/go-soccer-fan/issues/3)

<br/>

- [ ] `CAMPEONATO` -  Endpoints `REST` de consulta do(s) campeonato(s) `TDD` `happy-path`  [ISSUE-0004](https://github.com/jtonynet/go-soccer-fan/issues/4)
  > **Endpoint:** `GET /campeonatos`
  > Resposta esperada:
  > ```json
  > [
  >   { "id": "campeonato_001", "nome": "Campeonato Brasileiro", "temporada": "2025" },
  >   { "id": "campeonato_002", "nome": "UEFA Champions League", "temporada": "2025" }
  > ]
  > ```
  >
  > ---
  >
  > **Endpoint:** `GET /campeonatos/{id}/partidas?equipe=Flamengo&rodada=3`
  > - **Filtro Obrigatório:** 
  >   - O filtro de campeonato.
  > 
  > - **Filtros opcionais:**
  >   - Equipe: se informado, exibe apenas jogos da equipe.
  >   - Rodada: se informado, exibe apenas os jogos da rodada específica.
  > 
  > Se ambos os filtros forem usados juntos, exibir apenas os jogos dessa equipe na rodada escolhida.
  > 
  > Resposta esperada:
  > ```json
  > {
  >   "rodada": 3,
  >   "partidas": [
  >     { "time_casa": "Flamengo", "time_fora": "Palmeiras", "placar": "2-1" }
  >   ]
  > }
  >```

<br/>

  - [ ] `TORCEDORES` - Endpoint `REST` de Cadastro dos torcedores `TDD` `happy-path`  [ISSUE-0005](https://github.com/jtonynet/go-soccer-fan/issues/5)
    > **Objetivo:** <br/>Permitir que os usuários se cadastrem para receber notificações sobre os jogos do seu time.<br/>
    > **Descrição:**<br/>
    > Esse endpoint receberá os dados do torcedor e o registrará como destinatário de mensagens. É necessário que sejam enviados, por exemplo, nome, e-mail e o time de interesse.<br/><br/>
    >
    > ---
    >
    > **Endpoint**: `POST /torcedores`<br/>
    > Exemplo de corpo da requisição:
    > ```json
    > {
    >   "nome": "João Silva",
    >   "email": "joao.silva@example.com",
    >   "time": "Flamengo"
    > }
    > ```
    >
    > Resposta esperada:
    > ```json
    > {
    >   "id": "torcedor_001",
    >   "nome": "João Silva",
    >   "email": "joao.silva@example.com",
    >   "time": "Flamengo",
    >   "mensagem": "Cadastro realizado com sucesso"
    > }
    >```
    > **Notas:**<br/>
    > Realizar a validação dos dados (campos obrigatórios, formato de e-mail, etc.).
    > Armazenar os dados em uma base de dados ou outro mecanismo de persistência (ou, para testes, em memória).

<br/>

- [ ] `CAMPEONATO` - `Service` de Consumo da API externa de Campeonatos (`pill 💊`) [ISSUE-0006](https://github.com/jtonynet/go-soccer-fan/issues/6)
  - > [URL da API](http://api.football-data.org/) - [Documentação](https://www.football-data.org/documentation/quickstart)
  - Para maiores informacoes, verificar o `Documento de Requisitos` em `6) API Externa`.

<br/>

- [ ] `CAMPEONATO` - `Scheduller` Consumindo os dados de Campeonato da `API Externa` alimentando o nosso `DB` [ISSUE-0007](https://github.com/jtonynet/go-soccer-fan/issues/7)

<br/>

- [ ] `CAMPEONATO` - Testes e validações de `corner-cases`  (testes que o `happy-path` não deve ter validado) [ISSUE-0008](https://github.com/jtonynet/go-soccer-fan/issues/8)

<br/>

  - [ ] `TORCEDORES` - Endpoint `REST` e `Service` para `Broadcast` de notificações dos torcedores. A princípio para estratégia de `websockets`, **OPCIONAL** (Não cobertos na tarefa, mas para discutir no futuro) abranger outras estratégias de notificações como `push`, `email`, `SMS`, etc... [ISSUE-0009](https://github.com/jtonynet/go-soccer-fan/issues/9)
    > Tipos de Mensagens a serem transmitidas:
    > 1. Notificação de Início do Jogo: Informar os torcedores que o jogo vai começar.
    > 2. Notificação de Fim do Jogo: Informar o placar final do jogo.<br/>
    >
    > **Implementação:**
    > Para essa funcionalidade, pode-se optar por uma das seguintes abordagens:<br/>
    > - Endpoint REST para Broadcast:
    >   - Criar um endpoint, por exemplo, POST /broadcast, que receba o tipo de evento, time e os detalhes da mensagem, e então distribua as notificações para os torcedores cadastrados que acompanham o time informado.
    > - Comunicação em Tempo Real (ex.: WebSockets):<br/>
    > Implementar uma conexão WebSocket para enviar as mensagens assim que o evento ocorrer.
    > - Outros Mecanismos:<br/>
    > Pode-se utilizar notificações por e-mail, SMS ou push notifications, dependendo do escopo e das tecnologias escolhidas.
    > <br/><br/>
    > **Endpoint:** `POST /broadcast`<br/>
    > Exemplo de Payload para Broadcast (Início do Jogo):
    > 
    > ```json
    > {
    >   "tipo": "inicio",
    >   "time": "Flamengo",
    >   "mensagem": "O jogo do Flamengo vai começar em breve!"
    > }
    > ```
    >
    > Exemplo de Payload para Broadcast (Fim do Jogo):
    >```json
    > {
    >   "tipo": "fim",
    >   "time": "Flamengo",
    >   "placar": "2-1",
    >   "mensagem": "O jogo terminou com placar 2-1"
    > }
    >```
    - [ ] `DOC Opcional` - Diagrama `Mermaid` no `README` do seguinte fluxo (que deve ser obedecido no desenvolvimento da task):
      > - **Fluxo do Processo:**
      >   - Ao receber uma solicitação de broadcast, o sistema identifica todos os torcedores cadastrados que estão interessados no time especificado.
      >   - Em seguida, distribui a mensagem de notificação para cada torcedor.
      >   - Opcionalmente, registra/loga o envio das mensagens para controle e auditoria.

<br/>

- [ ] `TORCEDORES` - Testes e validações de `corner-cases`  (testes que o `happy-path` não deve ter validado) [ISSUE-0010](https://github.com/jtonynet/go-soccer-fan/issues/10)
  > - Desenvolver testes unitários e testes de integração para:
  >   - Validar o correto cadastro dos torcedores.
  >   - Verificar se a distribuição das mensagens broadcast ocorre conforme o esperado.

<br/>

- [ ] `GERAL` - Rota de Autenticação `JWT` [ISSUE-0011](https://github.com/jtonynet/go-soccer-fan/issues/11)
  > **Autenticação JWT**<br/>
  > - Criar um endpoint REST para login que gere um token JWT.
  > - Criar um endpoint REST para registrar novos usuários **(OPCIONAL)**.
  > <br/>
  > <br/>
  > **Endpoint para Login:** `POST /auth/login`<br/>
  > Exemplo de corpo da requisição:<br/>
  >
  > ```json
  > {
  >   "usuario": "admin",
  >   "senha": "123456"
  > }
  > ```
  >
  > Resposta esperada:
  >```json
  > {
  >   "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  > }
  >```

<br/>

- [ ] `GERAL` - Autenticar Rotas `GET /campeonatos`, `GET /campeonatos/{id}` e `POST /broadcast`. Alterar testes se necessário [ISSUE-0012](https://github.com/jtonynet/go-soccer-fan/issues/12)

<br/>

- [ ] **`GERAL` `OPCIONAL`** - Estrutura de Logs  `[ISSUE-???]`