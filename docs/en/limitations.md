# Limitations

[Documentation](README.md)

- The AI has limited depth and a simple evaluation: it can make mistakes and
  does not guarantee the best move in the game. A tactical sequence may continue
  beyond the search boundary, making an evaluation overly optimistic.
- Quiescence extends captures, en passant, promotions, and responses to check
  by up to eight additional plies. At that extension's boundary, it applies
  static evaluation; it can still miss longer tactics.
- There is no opening book or position cache. The two-second budget applies
  to deepening after the initial moves are generated; cancellation is
  cooperative and does not provide a real-time guarantee.
- The search recognizes insufficient material and the automatic 75-move limit,
  but does not receive repetition history; the game tracks those repetitions.
- Threefold repetition and the 50-move rule allow a draw claim in the current
  position through the game menu. The AI claims as soon as it can. Claims based
  on an intended move that has not yet been played are not supported.
- Fivefold repetition and 75 moves automatically end the game, preserving
  checkmate's priority. Insufficient material covers bare kings, a single bishop
  or knight, and positions with only bishops on the same square color in addition
  to the kings; other dead positions, such as complex blockades, are not detected.
  Reference: [FIDE rules, Articles 5 and 9](https://handbook.fide.com/chapter/e012023).
- History review lets you inspect previous positions without undoing moves
  or changing the game.
- There is no game saving or chess clock. Returning to the menu discards the game
  after confirmation; quitting the program does not save it either.
- Draw offers against the AI are unavailable: the engine does not yet have
  a policy for accepting or declining them.
