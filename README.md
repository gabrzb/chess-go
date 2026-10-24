# Chess

[English](README.md) | [Português](README-PT-BR.md)

Terminal based chess made in Go for two local players or against the AI.

## Project structure

```text
cmd/chess-tui/          CLI entry point
internal/application/  AI coordination
internal/domain/       Chess rules, game, and history
internal/engine/       AI evaluation and search
internal/delivery/tui/ Terminal interface
internal/architecture/ Dependency checks
docs/                  Documentation in English and Portuguese
```

## How to run

Requires Go 1.26.3 or later. From the project root:

```sh
go run ./cmd/chess-tui
```

With Task installed, use `task run`. To start directly against the AI:

```sh
go run ./cmd/chess-tui -ai medium
```

[Documentation: controls, options, architecture, rules, AI, and tests](docs/en/README.md).
