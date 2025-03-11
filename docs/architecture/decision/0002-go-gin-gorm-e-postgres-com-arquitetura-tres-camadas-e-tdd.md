# 2. 0002: Go, Gin, Gorm e PostgreSQL com Arquitetura Três Camadas e TDD

Data: 10 de março de 2025

## Status

Aceito

## Contexto

Precisamos definir o melhor fluxo de trabalho e testes para a go-payments-api.

Existem muitas opções para a arquitetura base.

O desafio cita implementação de `unit tests` e `integration tests` em minhas `services`, e `routes` para auxiliar no fluxo de desenvolvimento em `TDD`.

<div align="center">
<img src="../../assets/images/layout/graphics/test_pyramid.jpg">
<br/><i><a href="https://www.headspin.io/blog/the-testing-pyramid-simplified-for-one-and-all">*Imagem retirada do artigo: The Testing Pyramid: Simplified for One and All</a></i>
</div>


## Decisão

Este documento determina o fluxo de trabalho __Branch Based com Feature Branch__, design estrutural e a abordagem de testes para garantir um padrão para a aplicação.

O projeto sera realizado em Golang, em arquitetura `Três Camadas`, por ter maior familiaridade, além de considerá-las altamente performáticas.

O desafio deixa claro no trecho:

### Justificativa

#### Arquitetura Três Camadas
A arquitetura de três camadas foi adotada para garantir uma separação clara de responsabilidades, facilitando a manutenção e a escalabilidade da aplicação.

1. Camada de Controller(Manipuladores de rotas): Responsável por receber as requisições HTTP, validar entradas e delegar a lógica de negócio para a camada de serviços.

2. Camada de Serviço (Service Layer): Contém a lógica de negócio, garantindo que as regras da aplicação sejam aplicadas corretamente, sem depender diretamente da infraestrutura.

3. Camada de Repositório (Repository Layer): Gerencia o acesso ao banco de dados, abstraindo a persistência dos dados e permitindo a troca de tecnologias sem impacto na lógica de negócio.

#### GIN
Sugestão do desafio, selecionado como o framework de API por sua alta performance e simplicidade. 

#### GORM
Escolhemos GORM pela sua flexibilidade e capacidade de integração com os principais bancos de dados. A abstração oferecida pelo GORM simplifica a manipulação de dados, ao mesmo tempo que permite um controle refinado sobre queries mais complexas, quando necessário. Essa escolha também nos prepara para adoção futura de ferramentas de observabilidade e tracking, garantindo que a aplicação possa evoluir sem grandes refactors.

#### Postgres
Sugestão do desafio, robusto e features modernas, Postgres também é conhecido pela sua confiabilidade em ambientes de alta carg.


#### TDD
A adoção de TDD garante que a aplicação seja desenvolvida com um foco claro na cobertura de testes, minimizando bugs e retrabalho ao longo do ciclo de vida do projeto. Isso também nos prepara para uma maior resiliência em produção, especialmente considerando o impacto de falhas em um sistema financeiro.

## Consequências

A escolha dessas tecnologias, aliada a uma abordagem iterativa e incremental, permite que o projeto seja escalável e flexível. A documentação clara, através de ADRs e diagramas, garantirá que as equipes futuras possam se integrar ao projeto com facilidade. Além disso, a arquitetura adotada nos prepara para futuras mudanças tecnológicas, sem comprometer o core business ou aumentar os custos operacionais.

