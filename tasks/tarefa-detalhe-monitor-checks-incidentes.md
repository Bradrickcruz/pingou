# 🧩 Tarefa: Tela de detalhe do monitor com checks e incidentes

## 💡 Ideia

Adicionar uma página de detalhe do monitor para exibir últimos checks e incidentes relacionados ao monitor selecionado.

## 🎯 Objetivo

Permitir navegação do dashboard para uma visão detalhada de cada monitor, com histórico recente e contexto de incidentes.

## 📦 Escopo inicial

- Criar rota de detalhe do monitor no frontend
- Exibir últimos checks do monitor
- Exibir incidentes associados ao monitor
- Adicionar link de navegação a partir da lista de monitores

## 🚧 Fora de escopo

- Alterar o fluxo atual das páginas `/`, `/incidents` e `/settings`
- Reestruturar o backend para outras entidades
- Implementar gráficos ou métricas avançadas

## 🔍 Contexto

Hoje o frontend real usa apenas as rotas `/`, `/incidents` e `/settings`. O projeto não possui rota de detalhe por monitor, embora o PRD original tenha previsto essa visão.

## ⚠️ Observações

- Definir se a rota será `/monitors/:id` ou outro padrão
- Reutilizar dados já expostos por `/api/monitors/:id/checks` e `/api/monitors/:id/incidents` quando disponíveis
- Manter a UI consistente com o restante do dashboard

## 🏷️ Tipo

feature

## ⏱️ Prioridade

baixa
