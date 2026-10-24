# Limitações

[Documentação](README.md)

- A IA tem profundidade limitada e avaliação simples: pode cometer erros e não
  garante a melhor jogada da partida. Uma sequência tática pode continuar além
  do limite pesquisado e tornar uma avaliação otimista demais.
- A quiescência estende capturas, en passant, promoções e respostas a xeque
  por até oito plies adicionais. No limite dessa extensão, aplica a avaliação
  estática; ainda pode perder táticas mais longas.
- Não há livro de aberturas nem cache de posições. O orçamento de dois segundos
  vale para o aprofundamento após gerar as jogadas iniciais; o cancelamento é
  cooperativo e não constitui uma garantia de tempo real.
- A busca reconhece material insuficiente e o limite automático de 75 movimentos,
  mas não recebe o histórico de repetições; a partida controla essas repetições.
- Repetição tripla e regra dos 50 movimentos permitem reivindicar empate na
  posição atual pelo menu da partida. A IA reivindica assim que pode. Não há
  reivindicação antecipada indicando um lance ainda não executado.
- Repetição quíntupla e 75 movimentos encerram automaticamente a partida,
  preservando a prioridade do xeque-mate. Material insuficiente cobre reis
  sozinhos, um único bispo ou cavalo e somente bispos da mesma cor de casas;
  outras posições mortas, como bloqueios complexos, não são detectadas.
  Referência: [regras da FIDE, artigos 5 e 9](https://handbook.fide.com/chapter/e012023).
- A revisão do histórico permite consultar posições anteriores sem
  desfazer jogadas ou modificar a partida.
- Não há salvamento de partidas nem relógio. Voltar ao menu descarta a partida
  após confirmação; encerrar o programa também não a salva.
- Oferta de empate contra a IA não está disponível: o motor ainda não possui
  uma política para aceitar ou recusar.
