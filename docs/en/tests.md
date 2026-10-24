# Tests

[Documentation](README.md)

CI runs tests, `go vet`, and builds on Windows and Linux.

All `*_test.go` files, including benchmarks, live in
`internal/domain/tests`, `internal/engine/tests` or `internal/delivery/tui/tests`, grouped by subject.
Tests use the packages' public APIs; TUI flows are exercised through
keyboard events and rendered output.

```sh
go test ./...
go vet ./...
go test ./internal/engine/tests -run '^$' -bench BenchmarkChooseMove -benchmem
```

### Controls

- In menus, ↑/↓ and Enter navigate and select; ←/→ change settings.
- During a game, arrow keys or WASD move the cursor; Enter selects or confirms a move.
- Esc: cancel selection/promotion; with nothing selected, open the game menu.
- In the game menu, resume, resign, or return to the main menu. In local mode,
  you can also offer a draw; declining it keeps the game going.
- F: flip the board orientation.
- ?: open help; Esc returns to the previous screen.
- H: open or close history review.
- During review, ←/A and →/D step through positions; Esc returns to the game.
- During promotion, ↑/↓ and Enter choose a piece; Q/R/B/N are shortcuts for queen,
  rook, bishop, or knight. Esc cancels without making the move.
- Q opens the game menu; in the main menu, it quits. Ctrl+C always quits.
- In the main menu, A toggles animations for the current run. Enter skips
  the opening animation without starting a game; press it again to select.
- While the AI is thinking, you can open the menu, help, or history. The search
  is canceled and restarts when you return to the game.
- If the AI fails, the screen lets you retry or return to the menu. Technical
  details are written to `chess-go/chess-go.log` inside the user's cache directory
  (`os.UserCacheDir`). If the log cannot be opened, the game continues without logging.
