# Execução

[Documentação](README.md)

Com o [Task instalado](https://taskfile.dev/docs/installation/), os comandos
principais ficam disponíveis na raiz do projeto:

```sh
task                         # listar comandos
task run                     # abrir o menu principal
task ai                      # contra a IA em dificuldade média
task ai -- -ai hard -color black
task test                    # testes
task vet                     # análise estática
task fmt                     # formatação
task build                   # executável em bin/
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
