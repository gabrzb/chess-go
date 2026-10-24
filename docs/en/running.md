# Running

[Documentation](README.md)

With [Task installed](https://taskfile.dev/docs/installation/), the main commands
are available from the project root:

```sh
task                         # list commands
task run                     # open the main menu
task ai                      # play against the AI on medium difficulty
task ai -- -ai hard -color black
task test                    # run tests
task vet                     # static analysis
task fmt                     # format code
task build                   # build executable in bin/
task check                   # run tests, analysis, and build in sequence
```

The Go commands below are also available without Task.

From the project root:

```sh
go run ./cmd/chess-tui
```

The menu lets you start a local game or configure the AI: choose your pieces
(white, black, or random) and the difficulty. After the result, you can play
again, inspect the board and history, or return to the menu. A rematch uses
the same settings; the random option draws a new color.

Commands with `-ai` start a game directly. To enter local mode directly,
use `go run ./cmd/chess-tui -ai off`.

To play as White against the AI:

```sh
go run ./cmd/chess-tui -ai medium
```

To play as Black, letting the AI move first:

```sh
go run ./cmd/chess-tui -ai hard -color black
```

| Option | Values | Default |
| --- | --- | --- |
| `-ai` | `off`, `easy`, `medium`, `hard`; starts without going through the menu | Without the flag: main menu |
| `-color` | `white`, `black`: your color against the AI | `white` |
| `-seed` | Integer used to break ties between AI moves | `1` |
| `-no-animation` | Disable animations | Animations enabled |
