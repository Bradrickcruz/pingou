# Plano de Desenvolvimento: Aviso de manutenção começou/terminou

## Overview

Hoje o Pingou pausa os avisos durante uma janela de manutenção e não avisa ninguém: quem recebe os alertas só percebe pelo silêncio. Este plano manda um aviso nos canais (webhook, e-mail, ntfy) quando uma **janela agendada** começa e quando termina.

**O que já existe e será reaproveitado:**

- `Scheduler.applyMaintenance` ([maintenance.go](internal/scheduler/maintenance.go#L47)) roda a cada 10 s e na subida. Ele é o gatilho, e não precisa de checagem nova.
- `started_at` e `finalized_at` já garantem que cada efeito rode uma vez. Gravar o aviso na **mesma transação** dá o "uma vez só, mesmo reiniciando", sem coluna nova.
- O outbox (`notification_outbox`) já faz fila, retry, prazo e retoma pendentes na subida.

**Limites (decididos em `IDEAS.md`):**

- Pausa manual não avisa.
- Um interruptor global em Configurações, sem escolha por janela.
- O webhook sempre envia quando o interruptor está ligado.
- Um aviso só por grupo de monitores.

**Restrições:**

- Sem migration: `event_type` não tem `CHECK`, e `settings` é chave-valor.
- Sem dependência nova.
- SQLite com `MaxOpenConns=1`: dentro da transação, toda leitura usa o `tx`, nunca o `db` (usar o `db` trava esperando a própria conexão).

## Project Type

**FULLSTACK** — backend Go (scheduler, repositório, canais de aviso) e um interruptor e um aviso no painel React.

## Definition of Done (DoD)

Este plano está concluído quando:

- [x] Todos os Success Criteria estão marcados
- [x] Build, lint e type-check passam
- [x] Nenhum item pendente no Status de Execução
- [x] Evidências de verificação registradas (ver `docs/implementacoes.md`, item 18)

> Revisar o plano se: escopo mudar, nova dependência for descoberta, ou mais de 2 tasks forem bloqueadas.

## Success Criteria

- [x] Com o interruptor ligado, o início e o fim de uma janela agendada geram um aviso em cada canal ligado.
- [x] Com o interruptor desligado (padrão), nenhum aviso entra na fila e nada muda para quem já usa o Pingou.
- [x] Janela criada para N monitores gera **um** aviso de início e **um** de fim, listando os N.
- [x] Reiniciar no meio da janela não repete o aviso de início. Falha na transação não deixa aviso sem efeito, nem efeito sem aviso.
- [x] Janela cancelada antes de começar não avisa. Janela adiada avisa no horário novo. Janela cancelada em curso gera aviso de fim.
- [x] Pausa manual não gera aviso.
- [x] O aviso de fim diz só que a manutenção terminou, sem falar do estado do monitor.
- [x] O modal de agendar avisa que o motivo vai para os canais quando o interruptor está ligado.
- [x] O README documenta os dois eventos, o payload e o interruptor.

## Tech Stack

| Technology               | Choice Rationale                                                 |
| ------------------------ | ---------------------------------------------------------------- |
| Go + `database/sql`      | Stack atual; o tick e as transações já existem                   |
| `notification_outbox`    | Entrega com retry pronta; o aviso vira só mais uma linha         |
| `settings` (chave-valor) | Interruptor sem migration, como o fuso                           |
| React + tokens `pg-*`    | Padrão dos cards de Configurações; `check:tokens` barra cor crua |

## File Structure

Só dois arquivos novos de produção, marcados abaixo.

```text
internal/domain/
  outbox.go                  # EventMaintenanceStarted/Ended, SnapshotMaintenance, NewMaintenanceSnapshot
  maintenance.go             # MaintenanceGroupKey (agrupar janelas)
  repository.go              # StartGroup, FinalizeGroup no lugar de Start, Finalize
internal/repository/
  maintenance_repo.go        # grupo numa transação + outbox; guarda de starts_at
internal/scheduler/
  maintenance.go             # agrupa, lê o interruptor, acorda o worker
  scheduler.go               # WithMaintenanceNotices
internal/service/
  maintenance_notice.go      # (novo) monta as linhas do outbox; textos de e-mail/ntfy
  notifier.go                # webhook: maintenance_started/ended, campo monitors
  email_notifier.go          # case dos eventos novos
  ntfy_notifier.go           # case dos eventos novos
  settings_service.go        # maintenance_notices
internal/handler/settings.go # campo no PATCH
cmd/pingou/commands/serve.go # liga o interruptor no scheduler
web/src/components/settings/MaintenanceNoticesCard.jsx  # (novo)
web/src/pages/Settings.jsx
web/src/components/monitors/MaintenanceModal.jsx
README.md, IDEAS.md, docs/implementacoes.md
```

## Task Breakdown

| Task ID | Agent                   | Name                              | Priority | Complexity | Dependencies |
| ------- | ----------------------- | --------------------------------- | -------- | ---------- | ------------ |
| T1      | `go-backend-specialist` | Eventos e snapshot do grupo       | P0       | S          | -            |
| T2      | `go-backend-specialist` | Repositório: grupo numa transação | P0       | M          | T1           |
| T3      | `go-backend-specialist` | Interruptor nas settings          | P0       | S          | -            |
| T4      | `go-backend-specialist` | Scheduler: agrupar e enfileirar   | P1       | M          | T2, T3       |
| T5      | `go-backend-specialist` | Canais: webhook, e-mail e ntfy    | P1       | M          | T1           |
| T6      | `frontend-specialist`   | Interruptor e aviso no modal      | P2       | S          | T3           |
| T7      | `code-test`             | Testes de ponta a ponta do tick   | P3       | M          | T4, T5       |
| T8      | `document`              | README, IDEAS e histórico         | P3       | S          | T4, T5, T6   |
| T9      | `qa-review`             | QA com o binário real             | P3       | S          | T7, T8       |

### Task T1: Eventos e snapshot do grupo

- **Agent**: `go-backend-specialist`
- **Priority**: P0
- **Complexity**: S
- **Dependencies**: None
- **PRE-CONDITIONS**:
  - Nenhuma.
- **INPUT**: `EventSnapshot` descreve um monitor só; `MaintenanceWindow` tem uma linha por monitor e nenhum id de lote.
- **OUTPUT**:
  1. Constantes `EventMaintenanceStarted = "maintenance_started"` e `EventMaintenanceEnded = "maintenance_ended"`.
  2. `SnapshotMaintenance{Monitors []SnapshotMonitor, StartsAt, EndsAt, Reason *string}` e o campo `Maintenance *SnapshotMaintenance` (`omitempty`) em `EventSnapshot`.
  3. `NewMaintenanceSnapshot(windows, monitors, at)`. `OccurredAt` é o momento do tick, não `starts_at`: o prazo do outbox conta a partir dele, e uma janela notada depois de uma queda longa do Pingou expiraria antes de sair.
  4. `MaintenanceGroupKey(w)` = `starts_at|ends_at|reason`.
- **VERIFY**:
  - Teste: JSON do snapshot ida e volta. Snapshots antigos, sem o campo, continuam lendo.
  - Teste: motivos diferentes dão chaves diferentes; motivo nil é diferente de motivo preenchido.
- **ROLLBACK**: reverter o commit; nada usa os tipos ainda.

### Task T2: Repositório: grupo numa transação

- **Agent**: `go-backend-specialist`
- **Priority**: P0
- **Complexity**: M
- **Dependencies**: T1
- **PRE-CONDITIONS**:
  - T1 pronta.
- **INPUT**: `Start` e `Finalize` tratam uma janela por transação. Com N janelas, um crash no meio geraria avisos separados na subida.
- **OUTPUT**:
  1. `StartGroup(ctx, windows, now, build) ([]*MaintenanceWindow, error)` e `FinalizeGroup(ctx, windows, now, endsAt, build) ([]*MaintenanceWindow, error)` substituem `Start` e `Finalize`. Devolvem as janelas em que o efeito foi aplicado. Atualizar o fallback em `MaintenanceService.Cancel`.
  2. `build` tem o tipo `func(applied []*MaintenanceWindow, monitors []*Monitor) ([]OutboxEvent, error)`; `nil` = sem aviso. Só é chamado depois de aplicar, com as janelas que realmente mudaram: assim, uma janela apagada entre a leitura e a transação não entra na lista.
  3. Monitores lidos pelo `tx` com `scanMonitor`. Gravação com `OutboxRepoTx.EnqueueWithTx`, na mesma transação.
  4. **Aviso de fim só para janela que já tinha `started_at`** antes da transação (lido dentro dela).
  5. Guarda `AND starts_at <= ?` no `UPDATE` de `startWindowTx`. Hoje, adiar uma janela entre o `DueStart` e o `Start` aplica o início no horário velho, e "se o horário mudou, vale o novo".
- **VERIFY**:
  - Teste: `build` que devolve erro desfaz tudo (`started_at` nulo, nenhum incidente fechado, nenhuma linha no outbox).
  - Teste: segunda chamada não aplica nada e não enfileira nada.
  - Teste: janela adiada depois do `DueStart` não inicia.
  - Testes atuais de `maintenance_test.go` e `maintenance_service_test.go` verdes, com a chamada ajustada.
- **ROLLBACK**: reverter o commit; os dados não mudam de formato.

### Task T3: Interruptor nas settings

- **Agent**: `go-backend-specialist`
- **Priority**: P0
- **Complexity**: S
- **Dependencies**: None
- **PRE-CONDITIONS**:
  - Nenhuma.
- **INPUT**: `Settings` não tem o campo.
- **OUTPUT**:
  1. Chave `maintenance_notices`, valor `"true"` ou `"false"`; ausente = desligado.
  2. Campo `MaintenanceNotices bool` em `Settings` e `*bool` em `UpdateSettingsInput` e `updateSettingsRequest`.
- **VERIFY**:
  - Teste em `settings_test.go`: padrão `false`; `PATCH` liga e `GET` mostra; `PATCH` sem o campo não muda.
- **ROLLBACK**: reverter; a chave que sobrar em `settings` é ignorada.

### Task T4: Scheduler: agrupar e enfileirar

- **Agent**: `go-backend-specialist`
- **Priority**: P1
- **Complexity**: M
- **Dependencies**: T2, T3
- **PRE-CONDITIONS**:
  - T2 e T3 prontas.
- **INPUT**: `applyMaintenance` percorre as janelas uma a uma. `Scheduler` já tem `channels` e `wake`.
- **OUTPUT**:
  1. `service.NewMaintenanceEvents(eventType, applied, monitors, channels, now)`: um `event_id` por aviso, uma linha por canal, prazo de 60 min (`incidentNoticeTTL`). `monitor_id` da linha = primeiro monitor do grupo (a coluna é `NOT NULL`; o worker não depende dela para enviar).
  2. `WithMaintenanceNotices(enabled func(ctx) (bool, error))`, ligado em `serve.go` lendo `settingsRepository`. Lido só quando há janela vencida, para não somar uma consulta a cada 10 s.
  3. Em `applyMaintenance`: agrupar `DueStart` e `DueEnd` por `MaintenanceGroupKey`. Janelas de início com `ends_at <= now` ficam de fora do `StartGroup`, porque o `FinalizeGroup` já aplica o início que faltou. Por isso, janela que começou e terminou com o Pingou fora não gera nenhum aviso.
  4. `dbMu` cobre a transação do grupo. Depois do commit: `wake()` e `Reload` de cada janela finalizada.
  5. `EndMaintenance` (cancelamento em curso) vira um grupo de uma janela, com aviso de fim.
  6. Erro ao ler o interruptor: no tick, sai e tenta no próximo (o bloqueio dos checks vem do relógio, então nada se perde). No cancelamento, segue sem aviso e registra no log (cancelar não pode falhar por isso).
- **VERIFY**:
  - `go test -race ./internal/scheduler/...` verde.
  - Testes de T7.
- **ROLLBACK**: desligar o interruptor interrompe os avisos sem deploy. Reverter o commit volta ao comportamento atual.

### Task T5: Canais: webhook, e-mail e ntfy

- **Agent**: `go-backend-specialist`
- **Priority**: P1
- **Complexity**: M
- **Dependencies**: T1
- **PRE-CONDITIONS**:
  - T1 pronta.
- **INPUT**: Os três `Send` têm `switch` por `EventType`; tipo desconhecido vira erro. O payload do webhook tem `monitor` de um monitor só.
- **OUTPUT**:
  1. **Webhook:** `maintenance_started` e `maintenance_ended`, com o formato abaixo. `monitors` usa `toWebhookMonitor` (nunca leva `ping_token`). `monitor` sai do JSON desses eventos com `omitzero`. `dashboard_url` aponta para `/maintenance`. Os eventos existentes não mudam.
  2. **Textos** (em `maintenance_notice.go`, usados por e-mail e ntfy): assunto com um monitor `[Pingou] 🔧 Manutenção começou: <nome>`; com vários `[Pingou] 🔧 Manutenção começou em N monitores`. Corpo: monitores, início, fim e motivo (com `oneLine`), no fuso configurado. O aviso de fim diz só "a manutenção terminou".
  3. **Lista longa:** e-mail e ntfy mostram até 20 nomes e "e mais N". O ntfy transforma mensagem acima de 4 KB em anexo, e uma janela pode ter 100 monitores.
  4. **ntfy:** prioridade `default`, tag `wrench`, adicionar `🔧 ` em `ntfyTitle`, `click` para `/maintenance`.

  ```json
  {
    "id": "0193...",
    "event": "maintenance_started",
    "monitors": [
      {
        "id": "...",
        "name": "API",
        "type": "http",
        "target": "https://...",
        "url": "https://..."
      }
    ],
    "starts_at": "2026-11-08T02:00:00Z",
    "ends_at": "2026-11-08T04:00:00Z",
    "reason": "troca de disco",
    "timestamp": "2026-11-08T02:00:04Z",
    "dashboard_url": "https://pingou.exemplo.com/maintenance"
  }
  ```

- **VERIFY**:
  - Testes nos arquivos de notifier existentes: JSON exato dos dois eventos, sem `monitor` e sem `ping_token`; `monitor_down` idêntico ao de antes.
  - Teste: motivo com `\r\n` não vira cabeçalho no e-mail; lista de 100 monitores corta em 20; `reason` nil sai como `null` no webhook e sem a linha no texto.
- **ROLLBACK**: reverter o commit. Linhas pendentes dos eventos novos dão "unknown event type" no binário antigo, repetem e expiram em até 60 min, sem efeito colateral.

### Task T6: Interruptor e aviso no modal

- **Agent**: `frontend-specialist`
- **Priority**: P2
- **Complexity**: S
- **Dependencies**: T3
- **PRE-CONDITIONS**:
  - `GET /api/settings` devolve `maintenance_notices`.
- **INPUT**: `Settings.jsx` tem o formulário com barra de não salvos. O `MaintenanceModal` usa `ReasonField`, sem saber do interruptor.
- **OUTPUT**:
  1. `MaintenanceNoticesCard` no grupo "Notificações", com `<input type="checkbox" role="switch">` e `<label>`. Entra no `dirty`, no `discard` e no `update` do formulário.
  2. Descrição: "Avisa nos canais ligados quando uma janela agendada começa e termina. Pausa manual não avisa."
  3. `MaintenanceModal`, com o interruptor ligado: texto abaixo do motivo, "Este motivo vai junto no aviso de manutenção (webhook, e-mail e ntfy)." Valor vem de `useSettings`; enquanto carrega ou com erro, o texto não aparece e o modal funciona igual.
- **VERIFY**:
  - `npm run lint`, `npm run format:check`, `npm run check:tokens`, `npm run build` verdes.
  - No navegador: ligar, salvar, recarregar e continuar ligado; só com teclado; claro, escuro e mobile.
- **ROLLBACK**: reverter o commit; o backend continua aceitando o campo.

### Task T7: Testes de ponta a ponta do tick

- **Agent**: `code-test`
- **Priority**: P3
- **Complexity**: M
- **Dependencies**: T4, T5
- **PRE-CONDITIONS**:
  - T4 e T5 prontas.
- **INPUT**: `maintenance_test.go` tem `newMaintenanceEnv`, `addWindow`, `expire` e `countNotices`.
- **OUTPUT**: testes em `internal/scheduler/maintenance_test.go`:
  1. Interruptor desligado: nenhuma linha nos dois momentos.
  2. Três monitores, mesmo horário e motivo: um `event_id` de início e um de fim, com uma linha por canal.
  3. Dois ticks seguidos e nova `Start` do scheduler no meio da janela: sem repetir.
  4. Janela apagada antes do início: nenhuma linha. Janela adiada: nenhuma linha até o horário novo.
  5. Cancelamento em curso: um aviso de fim.
  6. Janela inteira no passado na subida: efeitos aplicados, nenhuma linha.
  7. Janela de um lote editada: vira aviso separado.
  8. Monitor ainda falhando ao fim: `maintenance_ended` sem texto de estado, e o `monitor_down` sai à parte pelo fluxo normal.
  9. Pausa manual: nenhuma linha.
  10. Grupo de 100 monitores: um aviso só, e o tempo da transação registrado em `docs/implementacoes.md`.
- **VERIFY**:
  - `go test -race ./...` verde.
- **ROLLBACK**: não se aplica (só testes).

### Task T8: README, IDEAS e histórico

- **Agent**: `document`
- **Priority**: P3
- **Complexity**: S
- **Dependencies**: T4, T5, T6
- **PRE-CONDITIONS**:
  - Payload e textos finais definidos.
- **INPUT**: README diz "ninguém é avisado" na janela e lista os eventos do webhook.
- **OUTPUT**:
  1. README, seção Manutenção: interruptor e as regras (agrupamento, cancelamento, janela com o Pingou fora).
  2. README, seção Webhook: tabela de eventos, payload dos dois e o aviso de que esses eventos usam `monitors` no lugar de `monitor`; evento desconhecido deve ser ignorado.
  3. README, Entrega das notificações: prazo de 60 min para os avisos de manutenção.
  4. Comentário de `event_type` em [0010_notification_outbox.sql](internal/database/migrations/0010_notification_outbox.sql#L14) fica como está (migration aplicada não muda); a lista passa a morar nas constantes.
  5. `IDEAS.md`: remover a ideia 6 e a linha do resumo. `docs/implementacoes.md`: novo item com decisões e validação.
- **VERIFY**:
  - Exemplo de JSON do README igual ao do teste de T5.
- **ROLLBACK**: reverter o commit.

### Task T9: QA com o binário real

- **Agent**: `qa-review`
- **Priority**: P3
- **Complexity**: S
- **Dependencies**: T7, T8
- **PRE-CONDITIONS**:
  - `make build` e um receptor de webhook local.
- **INPUT**: Binário com o interruptor ligado.
- **OUTPUT**:
  1. Janela de 2 min em 3 monitores: um aviso de início, um de fim, nos canais disponíveis.
  2. `kill -9` no meio da janela e subir: sem repetir o início.
  3. Receptor devolvendo 500 no início: retry com o mesmo `id`.
  4. Resultado registrado em `docs/implementacoes.md`; o que não der para testar (SMTP ou ntfy reais) vai para Pendências abertas.
- **VERIFY**:
  - Log do receptor e `notice sent` no log do Pingou.
- **ROLLBACK**: não se aplica.

## Risk Areas

| Risk                                          | Impact                             | Mitigation                                                                                    |
| --------------------------------------------- | ---------------------------------- | --------------------------------------------------------------------------------------------- |
| Integrador espera `monitor` em todo evento    | Fluxo quebra ao ler `monitor.name` | Interruptor começa desligado; README avisa; regra "ignore evento desconhecido" já documentada |
| Fluxo trata "não é down" como "voltou"        | Falso "voltou"                     | Renomeação para `monitor_down`/`monitor_up` já entregue; interruptor desligado por padrão     |
| Leitura pelo `db` dentro da transação         | Trava o tick, com `MaxOpenConns=1` | Tudo pelo `tx` em T2; `go test -race` com timeout                                             |
| Motivo com detalhe interno vai para terceiros | Vazamento em Slack ou e-mail       | Aviso no modal (T6); `reason` nunca é preenchido sozinho                                      |
| Grupo grande segura `dbMu` por mais tempo     | Atrasa checks por instantes        | Teste com 100 monitores em T7; só acontece nas viradas de janela                              |
| Cancelar um lote, um a um                     | N avisos de fim                    | Aceito: o painel cancela por janela; ver Notes                                                |

## Rollback Strategy

1. Desligar o interruptor em Configurações: para os avisos na hora, sem deploy.
2. Reverter os commits e publicar o binário anterior. Sem migration, nada a desfazer no banco.
3. Linhas pendentes dos eventos novos expiram em até 60 min. A chave `maintenance_notices` fica inerte em `settings`.

## Phase X: Verification

```bash
# Backend
go vet ./...
go test -race ./...

# Frontend
cd web && npm run lint && npm run format:check && npm run check:tokens && npm test && npm run build && cd ..

# Build completo
make build

# Grafo do projeto
graphify update .
```

## Cronograma Sugerido

| Fase   | Tarefas    | Complexity | Esforço | Target                                    | Status      |
| ------ | ---------- | ---------- | ------- | ----------------------------------------- | ----------- |
| Fase 0 | T1, T2, T3 | M          | ~1 h    | Base com grupo transacional e interruptor | ✅ Concluída |
| Fase 1 | T4, T5     | M          | ~1 h    | Avisos saindo nos três canais             | ✅ Concluída |
| Fase 2 | T6         | S          | ~15 min | Interruptor no painel                     | ✅ Concluída |
| Fase 3 | T7, T8, T9 | M          | ~1 h    | Testes, documentação e QA                 | ✅ Concluída |

## Status de Execução

### Concluído

- T1, T2, T3, T4, T5, T6, T7, T8, T9 (QA com SMTP e ntfy reais fica em Pendências abertas de `docs/implementacoes.md`)

### Em andamento

-

### Pendente

-

## Notes

**Decisões assumidas neste plano, a confirmar (não estavam em `IDEAS.md`):**

1. **Interruptor começa desligado.** Mesmo critério do aviso de certificado: quem atualiza não recebe evento novo sem pedir.
2. **O que é "um aviso só":** o banco tem uma janela por monitor e nenhum id de lote. O grupo é "mesmo início, fim e motivo", sem coluna nova. Dois lotes criados separadamente com valores iguais viram um aviso só, o que é até o desejado.
3. **Webhook com `monitors` (lista) e sem `monitor`,** mesmo com um monitor só. Um formato único evita dois casos para quem integra.
4. **Janela que começou e terminou com o Pingou desligado:** nenhum aviso. O início não saiu, então um "terminou" sozinho só confunde.
5. **Janela em curso cancelada** gera aviso de fim, porque o início já foi avisado.
6. **Prazo de 60 min,** igual a queda e recuperação. Passou disso, o aviso é descartado.
7. **Início notado com atraso** (Pingou voltou no meio da janela): o aviso sai, porque a manutenção ainda está valendo.

**Limitações aceitas:**

- Cancelar em curso as janelas de um lote, uma a uma, gera um aviso de fim por janela. Só resolve com cancelamento em lote, que não existe hoje.
- Editar o motivo de uma janela do lote em curso separa o aviso de fim dela. Texto correto, só um aviso a mais.
- Ligar o interruptor no meio de uma janela gera aviso de fim sem o de início.
- Página pública (ideia 3) e traduções (ideia 4) ainda não existem. Os textos nascem em português em `maintenance_notice.go`, junto dos demais, e a ideia 4 os leva com o resto.

### Alternativas consideradas e descartadas

- `Coluna batch_id em maintenance_windows` → descartada por exigir migration, e janela editada ainda precisaria sair do lote; a chave por valores resolve.
- `Colunas start_notified_at / end_notified_at` → descartada porque `started_at` e `finalized_at` na mesma transação já garantem o "uma vez só".
- `Enfileirar depois da transação do início` → descartada: um crash entre os dois perderia o aviso ou aplicaria o efeito sem ele.
- `Um ticker próprio a cada minuto` → descartada: `runMaintenance` já roda a cada 10 s.
- `Campo monitor com o primeiro do grupo no webhook` → descartada: em lote, afirma algo falso sobre os demais.
