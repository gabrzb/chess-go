# How the AI works

[Documentation](README.md)

The AI uses move search and an evaluation function written for this project.
It uses no training, neural networks, external services, or internet connection.

### 1. Position evaluation

A score estimates the advantage on the board. The unit is a centipawn:
100 points represent the material value of a pawn.

| Piece | Material value |
| --- | ---: |
| Pawn | 100 |
| Knight | 320 |
| Bishop | 330 |
| Rook | 500 |
| Queen | 900 |

Evaluation adds White's material and subtracts Black's. Positive values favor
White; negative values favor Black. Small bonuses reward knights, bishops,
and queens near the center, as well as central and advanced pawns.
The calculation mirrors positions to treat both colors equally.

The king has no capture value: the rules prevent its capture. The search
recognizes checkmate separately and gives it more weight than any material
advantage. It prefers faster mates and delays forced losses.
Stalemate scores zero, including at the search boundary.

### 2. Minimax and Negamax

The search simulates a tree: each position generates its legal moves, and each
move generates another position. Minimax assumes the opponent will choose their
best response. It therefore compares alternatives considering the worst response for us.

A simplified example, with scores from the choosing player's perspective:

| Alternative | Scores after possible opponent responses | Worst response |
| --- | --- | ---: |
| A | +300, -100 | -100 |
| B | +80, +40 | +40 |

The choice is B: the opponent's strongest response still leaves +40.
This is the idea behind [Minimax, described in UC Berkeley's AI material](https://inst.eecs.berkeley.edu/~cs188/textbook/games/minimax.html).

In this project, the algorithm uses the **Negamax** form: the score is viewed
from the side to move, and its sign is reversed when passing to the opponent.
This lets both colors share the same logic. The search deepens from 1 to the
selected depth, with a budget of two seconds per move. When the budget runs out,
it uses the last completed depth; if none finishes, it uses an already generated
legal move. Interface cancellations discard the result, even after a completed iteration.

### 3. Alpha-Beta pruning and move ordering

Alpha-Beta pruning stops analyzing a sequence when previously calculated bounds
show that it cannot improve the choice. It preserves the Minimax result at the
same depth and with the same evaluation while visiting fewer positions.
See [UC Berkeley's Alpha-Beta explanation](https://inst.eecs.berkeley.edu/~cs188/textbook/games/minimax.html#321-alpha-beta-pruning).

To find these bounds earlier, the search examines promotions and captures first,
prioritizing valuable victims captured by cheaper pieces. Move ordering includes
en passant. All simulations use the engine's rules, including castling,
en passant, promotion, and protection of your own king.

### 4. Difficulty levels and tie-breaking

Depth counts individual turns, called *plies*. A White move and a Black response
make up two plies.

| Difficulty | Depth | What it considers from the current position |
| --- | ---: | --- |
| Easy (`easy`) | 1 | AI move |
| Normal (`medium`) | 2 | AI move and opponent response |
| Hard (`hard`) | 3 | AI move, opponent response, and next AI move |

All three difficulty levels use the same evaluator. Increasing depth requires
more computation; these names do not correspond to an Elo rating.

When several moves have exactly the best score found, the AI breaks ties using
`-seed`. The same position, difficulty, and seed produce the same choice when
the same depth is completed within the budget. Changing the seed can vary ties;
it does not randomly select a move with a lower score. For example:

```sh
go run ./cmd/chess-tui -ai medium -seed 42
```

### 5. Relationship with the game and history

The search receives a complete copy of the position, including castling rights,
en passant, and counters. Simulating a sequence does not change the game or
record positions in its history. Only the chosen move is sent to the game
controller, which validates and records it like any other move.

Computation runs in the background while the TUI continues receiving keystrokes.
Opening the menu, help, or review, leaving the game, or quitting the program
cancels the search. Old results are associated with the session, request, and
original position, preventing them from changing another game or rematch.
The AI also chooses its promotion piece during search without asking the human player.
