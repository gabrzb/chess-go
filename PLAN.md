# Plano de adaptação da TUI

Referência visual e funcional: [UI.md](UI.md).

## Escopo e decisões

- Reaproveitar o tabuleiro atual, cursor, seleção, regras, promoção e revisão do histórico.
- Centralizar a navegação em um modelo raiz com tela ativa, dimensões, preferências,
  configurações da partida e identificação da sessão.
- Manter regras e resultados no motor; a interface apresenta os estados e envia ações.
- Preservar setas/WASD para movimento e H para histórico. A ajuda deve refletir esses
  atalhos, evitando o conflito de H com a alternativa hjkl do guia.
- Usar um componente de resultado para vitória, derrota e empate, parametrizado por
  modo, vencedor e motivo. No modo local, informar o lado vencedor.
- Manter a escolha por flags como atalho para iniciar uma partida configurada.
- Recursos opcionais só aparecem quando funcionarem. Oferta de empate contra IA
  fica fora do escopo inicial, pois ainda não existe política de aceitação.
- Cada bloco abaixo deve resultar em um commit coerente, com sua verificação.

## 1. Navegação, tema e menu principal

- [ ] Criar o controle de telas e transições no modelo raiz.
- [ ] Criar tema, faixas, títulos, opções com marcador e rodapé compartilhados.
- [ ] Implementar menu principal com modo local, modo IA e saída.
- [ ] Centralizar os conjuntos de conteúdo; preservar a coluna do marcador.

Aceite: navegar por setas/Enter, abrir cada modo e encerrar restaurando o terminal.
A abertura animada será acrescentada no bloco 7.

## 2. Configuração e criação de partidas

- [ ] Implementar configuração da IA: cor branca, preta ou aleatória e os três
  níveis existentes (easy, medium, hard), com rótulos em português.
- [ ] Revelar a cor sorteada na partida e armazenar as configurações da sessão.
- [ ] Permitir iniciar o modo local diretamente; sua confirmação é opcional no guia.
- [ ] Gerar uma identificação nova para cada partida.
- [ ] Associar respostas da IA à sessão, à posição e ao pedido de busca.
- [ ] Cancelar buscas ao abandonar a partida e descartar respostas antigas.

Aceite: iniciar ambos os modos pela interface, escolher qualquer cor e impedir que
uma busca da partida anterior altere a nova partida, mesmo em posições iguais.

## 3. Menu da partida, ajuda e confirmações

- [ ] Implementar menu com continuar, desistir e voltar ao menu principal.
- [ ] Implementar confirmações com a opção de continuar selecionada inicialmente.
- [ ] Implementar ajuda contextual com retorno à tela anterior.
- [ ] Esc cancela seleção/promoção primeiro; sem seleção, abre o menu da partida.
- [ ] Permitir menu e ajuda durante o cálculo da IA.
- [ ] Acrescentar inversão do tabuleiro por F, reaproveitando o renderizador atual.

Aceite: cancelar ações preserva a partida; confirmar abandono descarta a sessão
e cancela a IA. Textos de descarte devem corresponder à ausência de salvamento.

## 4. Desistência, acordo e resultado do motor

- [ ] Representar explicitamente vencedor e motivo do término no motor/snapshot.
- [ ] Implementar desistência com indicação de quem desistiu; não presumir que é
  sempre o lado da vez, especialmente enquanto a IA calcula.
- [ ] Implementar oferta, recusa e aceitação de empate no modo local.
- [ ] Bloquear novas jogadas após qualquer término.
- [ ] Preservar mate e afogamento já implementados.

Aceite: desistência atribui o vencedor correto; recusa mantém a partida; acordo
encerra sem vencedor. Ocultar oferta de empate contra IA.

## 5. Resultado, revanche e consulta final

- [ ] Apresentar vitória/derrota contra IA, vencedor local e empate com seu motivo.
- [ ] Implementar jogar novamente, ver tabuleiro e voltar ao menu principal.
- [ ] Iniciar revanche em posição nova, preservando as configurações escolhidas.
- [ ] Reaproveitar tabuleiro e histórico em modo de consulta sem executar jogadas.
- [ ] Esc na consulta final retorna ao resultado.

Aceite: resultado coincide com o motor; consulta preserva a posição final; revanche
não carrega seleção, mensagens, histórico ou respostas da IA da sessão anterior.

## 6. Painéis, promoção, erros e layout adaptável

- [ ] Exibir modo, jogador da vez, cor humana, dificuldade e último lance.
- [ ] Exibir lista de lances agrupada por brancas/pretas usando o histórico existente.
  A notação curta do mockup exige um formatador; não apresentar UCI como se fosse SAN.
- [ ] Mostrar rodapé contextual da seleção e mensagens curtas em português.
- [ ] Apresentar xeque textualmente sem interromper a partida.
- [ ] Transformar promoção em painel com setas/Enter, peças da cor correta e Esc.
  A partida só é alterada após a escolha definitiva.
- [ ] Apresentar falha da IA com tentar novamente ou voltar ao menu; registrar
  detalhes técnicos sem misturá-los ao conteúdo da TUI.
- [ ] Tratar dimensões e recompor layout largo/compacto sem perder posição e foco.
- [ ] Ocultar o painel de histórico antes de comprometer o tabuleiro; medir o mínimo
  real e mostrar aviso quando a janela não comportar o conteúdo.

Aceite: promoção cancelada e erro preservam a partida; redimensionamento mantém
estado e controles acessíveis. O tabuleiro existente continua sendo a base visual.

## 7. Animações e preferências visuais

- [ ] Acrescentar abertura curta com faixas e brilho único no logo.
- [ ] Consumir Enter ao pular a abertura, sem iniciar uma partida acidentalmente.
- [ ] Não repetir a abertura completa ao retornar ao menu.
- [ ] Acrescentar spinner somente durante a busca e transição curta no resultado.
- [ ] Permitir desativar animações; atualizar quadros por mensagens de tempo.
- [ ] Não bloquear Update nem animar continuamente o tabuleiro.

Aceite: navegação imediata, layout estável durante as animações e todos os estados
compreensíveis com os efeitos desativados. Partículas de vitória são opcionais.

## Política de testes durante a adaptação

- Manter uma suíte essencial com menos de 50 funções Test e menos de 50 casos
  enumerados pelo Go, incluindo os subtestes executáveis.
- Reduzir também linhas e duplicação; não atingir a meta apenas escondendo testes
  antigos em funções auxiliares ou criando um único teste gigante.
- Priorizar regras especiais, perft, preservação do estado, busca e transições da TUI.
- Evitar repetir no teste da interface todas as regras já verificadas pelo motor.
- Cobrir geometria e geração de movimentos pelo perft e pelas regras públicas,
  sem manter testes individuais de cada helper ou de cada peça.
- No perft, manter a posição inicial na profundidade 4 e as posições de referência
  na profundidade 3; dispensar as repetições das profundidades menores.
- Na IA, priorizar mate nas duas cores, comparação com Minimax completo,
  legalidade, seed reproduzível, cancelamento e ausência de jogadas legais.
- Novos testes devem substituir ou ampliar cenários relevantes da suíte compacta.
- Verificar cada bloco com testes pertinentes e task check antes de finalizar.
