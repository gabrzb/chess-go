# Running

[Documentation](README.md)

With [Task installed](https://taskfile.dev/docs/installation/), the main commands
are available from the project root:

```sh
task                         # list commands
task run                     # open the main menu
task ai                      # play against the AI on medium difficulty
task ai -- -ai hard -color black
task desktop                 # open the native desktop app
task test                    # run tests
task vet                     # static analysis
task fmt                     # format code
task build                   # build both executables in bin/
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

## Desktop app

`cmd/chess-desktop` opens the game in a native window drawn by
[MyGo](https://github.com/egoist/mygo), on macOS, Windows and Linux. As in the
TUI, it starts on a main menu (Player vs Player, Player vs AI, Quit); Player vs
AI leads to the setup of your pieces (white, black or random) and the
difficulty. The game screen scales the board to the window.

Drag a piece to a legal square, or click the piece and then its destination.
The move list beside the board reviews earlier positions (click a move, or use
← and →); `F` flips the board. `Esc` opens the game menu: resume, offer or
claim a draw, resign, or return to the main menu. When the game ends you can
play again, view the final board, or go back to the menu.

```sh
go run ./cmd/chess-desktop
go run ./cmd/chess-desktop -ai hard -color black
```

The desktop app accepts `-ai`, `-color`, and `-seed`. Passing `-ai easy`,
`-ai medium`, or `-ai hard` starts a game immediately. The default `-ai off`
opens the main menu.
