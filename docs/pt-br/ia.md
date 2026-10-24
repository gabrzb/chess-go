# Como a IA funciona

[Documentação](README.md)

A IA usa busca de jogadas e uma função de avaliação escrita no projeto.
Não utiliza treinamento, redes neurais, serviços externos ou conexão com a internet.

### 1. Avaliação da posição

Uma pontuação estima a vantagem no tabuleiro. A unidade é o centésimo de peão:
100 pontos representam o valor material de um peão.

| Peça | Valor material |
| --- | ---: |
| Peão | 100 |
| Cavalo | 320 |
| Bispo | 330 |
| Torre | 500 |
| Dama | 900 |

A avaliação soma o material das brancas e subtrai o das pretas. Valores positivos
favorecem as brancas; negativos favorecem as pretas. Pequenos bônus valorizam
cavalos, bispos e damas próximos do centro, além de peões centrais e avançados.
O cálculo espelha as posições para tratar as duas cores da mesma maneira.

O rei não recebe um preço de captura: as regras impedem sua captura. A busca
reconhece xeque-mate separadamente e atribui a ele uma importância maior que
qualquer vantagem material. Ela prefere mates mais rápidos e adia uma derrota
forçada. Afogamento recebe pontuação zero, inclusive no limite da busca.

### 2. Minimax e Negamax

A busca simula uma árvore: cada posição gera suas jogadas legais e cada jogada
gera outra posição. O Minimax supõe que o adversário escolherá a melhor resposta
para ele. Por isso, compara as alternativas considerando a pior resposta para nós.

Exemplo simplificado, com notas vistas pelo jogador que está escolhendo:

| Alternativa | Notas após possíveis respostas adversárias | Pior resposta |
| --- | --- | ---: |
| A | +300, -100 | -100 |
| B | +80, +40 | +40 |

A escolha é B: a resposta adversária mais forte ainda deixa +40.
Essa é a ideia do [Minimax, descrita no material de IA da UC Berkeley](https://inst.eecs.berkeley.edu/~cs188/textbook/games/minimax.html).

Neste projeto, o algoritmo é escrito na forma **Negamax**: a nota é vista por quem
tem a vez, e seu sinal é invertido ao passar para o adversário. Isso permite usar
a mesma lógica para as duas cores. A busca aprofunda de 1 até a profundidade
escolhida, com orçamento de dois segundos por jogada. Ao esgotar o orçamento,
usa a última profundidade concluída;
se nenhuma terminar, usa uma jogada legal já gerada. Cancelamentos da interface
descartam o resultado, inclusive após uma iteração concluída.

### 3. Poda Alpha-Beta e ordenação

A poda Alpha-Beta interrompe a análise de uma sequência quando os limites já
calculados mostram que ela não pode melhorar a escolha. Ela preserva o resultado
do Minimax na mesma profundidade e com a mesma avaliação, visitando menos posições.
Veja a [explicação de Alpha-Beta da UC Berkeley](https://inst.eecs.berkeley.edu/~cs188/textbook/games/minimax.html#321-alpha-beta-pruning).

Para encontrar esses limites mais cedo, a busca examina primeiro promoções e
capturas, dando prioridade a vítimas valiosas capturadas por peças mais baratas.
A ordenação inclui en passant. Todas as simulações usam as regras do motor,
incluindo roque, en passant, promoção e proteção do próprio rei.

### 4. Dificuldades e desempate

A profundidade conta turnos individuais, chamados *plies*. Um movimento das
brancas e uma resposta das pretas correspondem a dois plies.

| Dificuldade | Profundidade | O que considera a partir da posição atual |
| --- | ---: | --- |
| Fácil (`easy`) | 1 | Jogada da IA |
| Normal (`medium`) | 2 | Jogada da IA e resposta adversária |
| Difícil (`hard`) | 3 | Jogada da IA, resposta adversária e próxima jogada da IA |

As três dificuldades usam o mesmo avaliador. Aumentar a profundidade exige mais
cálculo; esses nomes não correspondem a uma classificação Elo.

Quando várias jogadas têm exatamente a melhor nota encontrada, a IA desempata
usando `-seed`. A mesma posição, dificuldade e seed produzem a mesma escolha quando a mesma
profundidade é concluída dentro do orçamento.
Mudar a seed pode variar os empates; ela não escolhe aleatoriamente uma jogada
com nota inferior. Por exemplo:

```sh
go run ./cmd/chess-tui -ai medium -seed 42
```

### 5. Relação com a partida e o histórico

A busca recebe uma cópia completa da posição, incluindo direitos de roque,
en passant e contadores. Simular uma sequência não modifica a partida nem grava
posições no histórico. Somente a jogada escolhida é enviada ao controle da partida,
que a valida e a registra como qualquer outra jogada.

O cálculo roda em segundo plano enquanto a TUI continua recebendo teclas.
Ao abrir menu, ajuda ou revisão, abandonar a partida ou encerrar o programa,
a busca recebe um cancelamento. Resultados antigos são associados à sessão,
ao pedido e à posição original, impedindo que alterem outra partida ou revanche.
A IA também escolhe a peça de
promoção durante a busca, sem solicitar uma escolha ao jogador humano.
