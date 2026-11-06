# Plano de Desenvolvimento: {{Titulo}}

## Overview

{{Contexto + objetivo + boundary + restricoes}}

## Project Type

**{{BACKEND/FRONTEND/FULLSTACK}}** — {{descricao}}

## Definition of Done (DoD)

Este plano está concluído quando:

- [ ] Todos os itens do Success Criteria estão marcados
- [ ] Build, lint e type-check passam
- [ ] Nenhum item pendente no Status de Execução
- [ ] Evidências de verificação registradas

> Revisar o plano se: escopo mudar, nova dependência for descoberta, ou mais de 2 tasks forem bloqueadas.

## Success Criteria

- [ ] {{criterio 1}}
- [ ] {{criterio 2}}

## Tech Stack

| Technology | Choice Rationale |
| ---------- | ---------------- |
| {{tech}}   | {{rationale}}    |

## File Structure

```text
{{estrutura alvo}}
```

## Task Breakdown

| Task ID | Agent     | Name     | Priority | Complexity      | Dependencies |
| ------- | --------- | -------- | -------- | --------------- | ------------ |
| T1      | {{agent}} | {{name}} | P0       | {{XS/S/M/L/XL}} | -            |

### Task T1: {{Título}}

- **Agent**: {{agent}}
- **Priority**: P0
- **Complexity**: {{XS/S/M/L/XL}}
- **Dependencies**: None
- **PRE-CONDITIONS**:
  - {{pré-condição 1}}
- **INPUT**: {{estado atual}}
- **OUTPUT**:
  1. {{entrega 1}}
- **VERIFY**:
  - {{evidência 1}}

## Risk Areas

| Risk        | Impact      | Mitigation    |
| ----------- | ----------- | ------------- |
| {{risco 1}} | {{impacto}} | {{mitigação}} |

## Rollback Strategy

### Fase 1

1. {{rollback step}}

## Phase X: Verification

### Build Verification

```bash
{{build}}
```

### Lint & Type Check

```bash
{{lint}} && {{type-check}}
```

## Cronograma Sugerido

| Fase   | Tarefas  | Complexity | Esforço   | Target   | Status      |
| ------ | -------- | ---------- | --------- | -------- | ----------- |
| Fase 0 | {{T0.x}} | {{M}}      | {{tempo}} | {{meta}} | ⬜ Pendente |

## Status de Execução

### Concluído

-

### Em andamento

-

### Pendente

-

## Notes

- {{nota}}

### Alternativas consideradas e descartadas

- `{{opção}}` → descartada por `{{motivo}}`
