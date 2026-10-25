# Chess Go

[English](README.md) | [Português](README-PT-BR.md)

Jogue xadrez no terminal com dois jogadores locais ou contra uma IA integrada.
Desenvolvido em Go, com interface de terminal construída com Bubble Tea e Lip Gloss.

## Recursos

- **Partidas locais e contra a IA:** compartilhe o terminal ou escolha entre três níveis de dificuldade.
- **Regras de xadrez:** movimentos legais, xeque, xeque-mate, afogamento, roque, en passant e promoção com escolha da peça.
- **Interface por teclado:** selecione peças, veja os destinos legais e inverta a orientação do tabuleiro.
- **Revisão do histórico:** consulte posições anteriores sem alterar a partida em andamento.
- **IA offline:** a busca de jogadas e a avaliação de posições executam localmente, sem serviços externos.

Consulte as [limitações](docs/pt-br/limitacoes.md) para conhecer as condições de empate suportadas e outras restrições.

## Início rápido

Requer **Go 1.26.3 ou superior**. Na raiz do projeto, abra o menu principal:

```sh
go run ./cmd/chess-tui
```

No menu, escolha uma partida local ou configure o adversário de IA. Para iniciar
uma partida contra a IA na dificuldade média sem passar pelo menu:

```sh
go run ./cmd/chess-tui -ai medium
```

Use as setas e Enter para navegar. Durante a partida, pressione `?` para abrir a ajuda.
Consulte [Execução](docs/pt-br/execucao.md) para opções da CLI e comandos do Task,
e [Controles](docs/pt-br/testes.md#controles) para a referência completa do teclado.

## Documentação

Acesse o [índice da documentação](docs/pt-br/README.md) para os guias detalhados e as referências técnicas.

| Guia | Conteúdo |
| --- | --- |
| [Execução](docs/pt-br/execucao.md) | Modos de jogo, opções da CLI e comandos do Task |
| [Regras](docs/pt-br/regras.md) | Regras de xadrez implementadas |
| [IA](docs/pt-br/ia.md) | Avaliação, busca e níveis de dificuldade |
| [Arquitetura](docs/pt-br/arquitetura.md) | Responsabilidades dos pacotes e limites de dependência |
| [Testes](docs/pt-br/testes.md) | Comandos de verificação, benchmarks e controles |
| [Limitações](docs/pt-br/limitacoes.md) | Restrições da IA, tratamento de empates e recursos indisponíveis |

## Desenvolvimento

O projeto separa as regras de xadrez, a busca da IA e a apresentação no terminal:

```text
cmd/chess-tui/          Entrada da CLI
internal/application/  Coordenação da IA
internal/domain/       Regras de xadrez, partida e histórico
internal/engine/       Avaliação e busca da IA
internal/delivery/tui/  Interface de terminal
internal/architecture/ Verificação de dependências
docs/                  Documentação em inglês e português
```

Execute as verificações na raiz do projeto:

```sh
go test ./... -count=1 -timeout 60s
go vet ./...
go build ./cmd/chess-tui
```

Com o Task instalado, `task check` executa testes, análise estática e compilação.
A CI executa essas verificações no Linux e no Windows. Consulte o
[guia de arquitetura](docs/pt-br/arquitetura.md) antes de alterar as dependências entre pacotes.

## Licença

[MIT](LICENSE).
