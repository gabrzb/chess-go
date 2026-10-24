# Testes

[Documentação](README.md)

O CI executa testes, `go vet` e compilação no Windows e no Linux.

Todos os arquivos `*_test.go`, incluindo benchmarks, ficam em
`internal/domain/tests`, `internal/engine/tests` ou `internal/delivery/tui/tests`, separados por assunto.
Os testes usam as APIs públicas dos pacotes; os fluxos da TUI são exercitados
por eventos de teclado e pela saída renderizada.

```sh
go test ./...
go vet ./...
go test ./internal/engine/tests -run '^$' -bench BenchmarkChooseMove -benchmem
```

### Controles

- Nos menus, ↑/↓ e Enter navegam e selecionam; ←/→ alteram as configurações.
- Na partida, setas ou WASD movem o cursor; Enter seleciona ou confirma a jogada.
- Esc: cancelar seleção/promoção; sem seleção, abrir o menu da partida.
- No menu da partida, continuar, desistir ou voltar ao menu principal. No modo
  local, também oferecer empate; recusar mantém a partida.
- F: inverter a orientação do tabuleiro.
- ?: abrir a ajuda; Esc retorna à tela anterior.
- H: abrir ou fechar a revisão do histórico.
- Na revisão, ←/A e →/D percorrem as posições; Esc retorna à partida.
- Na promoção, ↑/↓ e Enter escolhem a peça; Q/R/B/N são atalhos para dama,
  torre, bispo ou cavalo. Esc cancela sem executar o movimento.
- Q abre o menu da partida; no menu principal, encerra. Ctrl+C encerra sempre.
- No menu principal, A liga ou desliga animações nesta execução. Enter pula
  a abertura sem iniciar uma partida; pressione novamente para selecionar.
- Durante o cálculo da IA, é possível abrir menu, ajuda ou histórico. A busca
  é cancelada e recomeça ao retornar à partida.
- Se a IA falhar, a tela permite tentar novamente ou voltar ao menu. Detalhes
  técnicos ficam em `chess-go/chess-go.log`, dentro da pasta de cache do usuário
  (`os.UserCacheDir`). Se o log não puder ser aberto, o jogo continua sem gravá-lo.
