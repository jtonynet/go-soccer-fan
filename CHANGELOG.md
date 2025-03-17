# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [Unreleased]
### Added

## [0.0.9] - 2025-17-03
### Added
- [Issue 9](https://github.com/jtonynet/go-soccer-fan/issues/9)
  - Adicionado `RabbitMQ`
  - Adicionados `Mailhog` e `stunnel`
  - Adicionado service de `Broadcast`
  - Broadcast de email escalável concluido

## [0.0.8] - 2025-15-03
### Added
- [Issue 7](https://github.com/jtonynet/go-soccer-fan/issues/7)
  - Adicionada `CLI` para consumo de `API externa`

## [0.0.7] - 2025-15-03
### Added
- [Issue 6](https://github.com/jtonynet/go-soccer-fan/issues/6)
  - Adicionada `service` para consumo de `API externa`
  - Utilizando [`exponentialBackoff`](github.com/cenkalti/backoff) permitir consumo da `API` e seu `rate-limit`
  - Melhorias nas `repositories`

## [0.0.6] - 2025-12-03
### Added
- [Issue 18](https://github.com/jtonynet/go-soccer-fan/issues/18)
  - Adicionados repositorios `GORM`
  - Melhorias no `REST`

## [0.0.5] - 2025-11-03
### Added
- [Issue 17](https://github.com/jtonynet/go-soccer-fan/issues/17)
  - Dockerizado e Postgres adiconado

## [0.0.4] - 2025-11-03
### Added
- [Issue 5](https://github.com/jtonynet/go-soccer-fan/issues/5)
  - Endpoint criação de torcedores 

## [0.0.3] - 2025-11-03
### Added
- [Issue 4](https://github.com/jtonynet/go-soccer-fan/issues/4)
  - Endpoint de consulta de campeonatos

## [0.0.2] - 2025-11-03
### Added
- [Issue 3](https://github.com/jtonynet/go-soccer-fan/issues/3)
  - Criada Service com `TDD` com repository fake sendo injetada. Performando consultas que obedecem ao `Diagrama Entidade Relacionamento` anteriormente criado.

## [0.0.1] - 2025-10-03
### Added
- [Issue 2](https://github.com/jtonynet/go-soccer-fan/issues/2)
  - Adicionado `Diagrama Entidade Relacionamento` base ao `README`

## [0.0.0] - 2025-10-03
### Added

- [Issue 1](https://github.com/jtonynet/go-soccer-fan/issues/1)
  - [Kanban Project View Iniciado](https://github.com/users/jtonynet/projects/8/views/1) com o commit inicial. 
  - Documentação base: Readme Rico, ADRs: [0001: Registro de Decisões de Arquitetura (ADR)](./docs/architecture/decisions/registro-de-decisoes-de-arquitetura.md) e [0002: Go, Gin, Gorm e PostgreSQL com Arquitetura Três Camadas e TDD](./docs/architecture/decisions/0002-go-gin-gorm-e-postgres-com-arquitetura-tres-camadas-e-tdd.md).
  - Sabemos o que fazer, graças às definições do arquivo __README.md__. Sabemos como fazer graças aos __ADRs__ e documentações vinculadas. Devemos nos organizar em estrutura __Kanban__, guiados pelo modelo Agile, em nosso __Github Project__, e dar o devido prosseguimento às tarefas.

[0.0.9]: https://github.com/jtonynet/go-soccer-fan/compare/v0.0.8...v0.0.9
[0.0.8]: https://github.com/jtonynet/go-soccer-fan/compare/v0.0.7...v0.0.8
[0.0.7]: https://github.com/jtonynet/go-soccer-fan/compare/v0.0.6...v0.0.7
[0.0.6]: https://github.com/jtonynet/go-soccer-fan/compare/v0.0.5...v0.0.6
[0.0.5]: https://github.com/jtonynet/go-soccer-fan/compare/v0.0.4...v0.0.5
[0.0.4]: https://github.com/jtonynet/go-soccer-fan/compare/v0.0.3...v0.0.4
[0.0.3]: https://github.com/jtonynet/go-soccer-fan/compare/v0.0.2...v0.0.3
[0.0.2]: https://github.com/jtonynet/go-soccer-fan/compare/v0.0.1...v0.0.2
[0.0.1]: https://github.com/jtonynet/go-soccer-fan/compare/v0.0.0...v0.0.1
[0.0.0]: https://github.com/jtonynet/go-soccer-fan/releases/tag/v0.0.0
