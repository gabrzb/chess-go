# Chess Go

[English](README.md) | [Português](README-PT-BR.md)

Play chess in your terminal with two local players or against a built-in AI.
Written in Go, with a terminal user interface built with Bubble Tea and Lip Gloss.

## Features

- **Local and AI play:** share a terminal or choose from three AI difficulty levels.
- **Chess rules:** legal moves, check, checkmate, stalemate, castling, en passant, and promotion with piece selection.
- **Keyboard interface:** select pieces, see legal destinations, and flip the board.
- **History review:** inspect earlier positions without changing the current game.
- **Offline AI:** move search and position evaluation run locally, without external services.

See the [limitations](docs/en/limitations.md) for supported draw conditions and other constraints.

## Quick start

Requires **Go 1.26.3 or later**. From the project root, open the main menu:

```sh
go run ./cmd/chess-tui
```

Choose a local game or configure an AI opponent in the menu. To start a game
against the AI on medium difficulty directly:

```sh
go run ./cmd/chess-tui -ai medium
```

Use the arrow keys and Enter to navigate. During a game, press `?` for help.
See [Running](docs/en/running.md) for CLI options and Task commands, and
[Controls](docs/en/tests.md#controls) for the full keyboard reference.

## Documentation

Use the [documentation index](docs/en/README.md) for detailed guides and technical references.

| Guide | Contents |
| --- | --- |
| [Running](docs/en/running.md) | Game modes, CLI options, and Task commands |
| [Rules](docs/en/rules.md) | Implemented chess rules |
| [AI](docs/en/ai.md) | Evaluation, search, and difficulty levels |
| [Architecture](docs/en/architecture.md) | Package responsibilities and dependency boundaries |
| [Tests](docs/en/tests.md) | Verification commands, benchmarks, and controls |
| [Limitations](docs/en/limitations.md) | AI constraints, draw handling, and unavailable features |

## Development

The project separates chess rules, AI search, and terminal presentation:

```text
cmd/chess-tui/          CLI entry point
internal/application/  AI coordination
internal/domain/       Chess rules, game, and history
internal/engine/       AI evaluation and search
internal/delivery/tui/  Terminal interface
internal/architecture/ Dependency checks
docs/                  English and Portuguese documentation
```

Run the checks from the project root:

```sh
go test ./... -count=1 -timeout 60s
go vet ./...
go build ./cmd/chess-tui
```

With Task installed, `task check` runs tests, static analysis, and a build.
CI runs these checks on Linux and Windows. See the
[architecture guide](docs/en/architecture.md) before changing package dependencies.

## License

[MIT](LICENSE).
