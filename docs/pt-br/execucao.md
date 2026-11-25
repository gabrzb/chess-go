# Execução

[Documentação](README.md)

Com o [Task instalado](https://taskfile.dev/docs/installation/), os comandos
principais ficam disponíveis na raiz do projeto:

```sh
task                         # listar comandos
task run                     # abrir o menu principal
task ai                      # contra a IA em dificuldade média
task ai -- -ai hard -color black
task desktop                 # abrir o app desktop nativo
task test                    # testes
task vet                     # análise estática
task fmt                     # formatação
task build                   # os dois executáveis em bin/
task check                   # testes, análise e compilação em sequência
```

Os comandos Go abaixo continuam disponíveis sem o Task.

Na raiz do projeto:

```sh
go run ./cmd/chess-tui
```

O menu permite iniciar uma partida local ou configurar a IA: escolha suas peças
(brancas, pretas ou aleatório) e a dificuldade. Após o resultado, é possível jogar
novamente, consultar o tabuleiro e o histórico ou voltar ao menu. A revanche usa
as mesmas configurações; a opção aleatória faz um novo sorteio.

Os comandos com `-ai` iniciam diretamente a partida. Para entrar diretamente no
modo local, use `go run ./cmd/chess-tui -ai off`.

Para jogar de brancas contra a IA:

```sh
go run ./cmd/chess-tui -ai medium
```

Para jogar de pretas, deixando a IA começar:

```sh
go run ./cmd/chess-tui -ai hard -color black
```

| Opção | Valores | Padrão |
| --- | --- | --- |
| `-ai` | `off`, `easy`, `medium`, `hard`; inicia sem passar pelo menu | Sem a flag: menu principal |
| `-color` | `white`, `black`: sua cor contra a IA | `white` |
| `-seed` | Inteiro usado no desempate da IA | `1` |
| `-no-animation` | Desativar animações | Animações ligadas |

## App desktop

`cmd/chess-desktop` abre o jogo em uma janela nativa desenhada pelo
[MyGo](https://github.com/egoist/mygo), no macOS, Windows e Linux. Como no TUI,
começa por um menu principal (Player vs Player, Player vs AI, Quit); Player vs AI
leva à escolha das suas peças (brancas, pretas ou aleatório) e da dificuldade. A
tela de jogo ajusta o tabuleiro ao tamanho da janela.

Arraste uma peça até uma casa legal, ou clique na peça e depois no destino. A
lista de lances ao lado do tabuleiro permite rever posições anteriores (clique em
um lance, ou use ← e →); `F` vira o tabuleiro. `Esc` abre o menu da partida:
continuar, oferecer ou reclamar empate, desistir ou voltar ao menu principal. Ao
fim da partida, é possível jogar de novo, ver o tabuleiro final ou voltar ao menu.

```sh
go run ./cmd/chess-desktop
go run ./cmd/chess-desktop -ai hard -color black
```

O app desktop aceita `-ai`, `-color` e `-seed`. Use `-ai easy`, `-ai medium`
ou `-ai hard` para começar uma partida direto. O padrão `-ai off` abre o menu
principal.
