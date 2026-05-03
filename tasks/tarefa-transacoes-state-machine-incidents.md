# Tarefa: transacoes SQL na state machine de incidentes

## Status

A fazer.

## Objetivo

Garantir atomicidade no fluxo de state machine que persiste check, atualiza monitor e abre/fecha incidentes.

Hoje `StateMachine.Process` executa operacoes sequenciais via repositorios. Se uma etapa falhar depois de outra ja ter sido gravada, pode haver estado parcialmente persistido.

## Escopo inicial

- Envolver em transacao as operacoes relacionadas a um resultado de check:
  - inserir `checks`;
  - atualizar `monitors.current_state`, `last_checked_at`, `updated_at`;
  - abrir incidente quando transita para DOWN;
  - fechar incidente quando transita para UP.
- Manter notificacoes fora da transacao, apos commit bem-sucedido.
- Preservar comportamento atual de threshold e transicoes.

## Implementacao sugerida

1. Definir como repositorios recebem `*sql.Tx` sem duplicar toda logica.
2. Criar camada de unidade de trabalho no pacote repository ou service.
3. Alterar `StateMachine.Process` para iniciar transacao no começo do processamento.
4. Executar inserts/updates dentro da mesma transacao.
5. Fazer commit antes de chamar notifier.
6. Fazer rollback em qualquer erro antes do commit.

## Pontos de atencao

- Notificacao webhook nao deve acontecer antes do commit.
- Se webhook falhar apos commit, estado do banco deve continuar correto.
- Evitar transacao longa envolvendo chamada HTTP externa.
- Garantir que `countConsecutiveFails` leia os dados esperados dentro da transacao.
- Considerar impacto do SQLite com `SetMaxOpenConns(1)`.

## Fora de escopo

- Reescrever state machine como funcao pura.
- Migrar para event-driven por channels.
- Criar testes automatizados agora.

## Criterios de aceite

- Check + mudanca de estado + incidente sao persistidos atomicamente.
- Em erro no meio do fluxo, nenhuma escrita parcial fica no banco.
- Webhook so dispara apos commit.
- `go test ./...` passa.
