# Xadrez

[English](README.md) | [Português](README-PT-BR.md)

Xadrez que funciona no terminal feito em Go para dois jogadores locais ou contra a IA.

## Estrutura do projeto

```text
cmd/chess-tui/          Entrada da CLI
internal/application/  Coordenação da IA
internal/domain/       Regras de xadrez, partida e histórico
internal/engine/       Avaliação e busca da IA
internal/delivery/tui/ Interface de terminal
internal/architecture/ Verificação de dependências
docs/                  Documentação em inglês e português
```

## Como rodar

Requer Go 1.26.3 ou superior. Na raiz do projeto:

```sh
go run ./cmd/chess-tui
```

Com Task instalado, use `task run`. Para iniciar diretamente contra a IA:

```sh
go run ./cmd/chess-tui -ai medium
```

[Documentação: controles, opções, arquitetura, regras, IA e testes](docs/pt-br/README.md).
