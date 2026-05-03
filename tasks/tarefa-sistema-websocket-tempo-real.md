# 🧩 Tarefa: Sistema WebSocket para atualização em tempo real dos monitores

## 💡 Ideia
Implementar comunicação via WebSocket entre o Dashboard e os schedulers para que mudanças de estado dos monitores sejam refletidas instantaneamente na interface, sem necessidade de polling manual.

## 🎯 Objetivo
Atualizar automaticamente o estado dos monitores no Dashboard quando os schedulers detectarem mudanças, eliminando a necessidade de refresh manual ou polling.

## 📦 Escopo inicial
- Servidor WebSocket (backend)
- Cliente WebSocket no frontend (Dashboard)
- Evento de push quando scheduler altera estado de monitor
- Reconexão automática em caso de queda

## 🚧 Fora de escopo
- Autenticação WebSocket (usar existente)
- Histórico de estados ou logging
- Outras páginas além do Dashboard

## 🔍 Contexto
O Dashboard em `web/src/pages/Dashboard.jsx` usa `useMonitors` hook com `refetch()` para obter dados. Os schedulers (background jobs) alteram o estado dos monitors mas não comunicam ao frontend.

## ⚠️ Observações
- Definir formato das mensagens WebSocket
- Escolher biblioteca (socket.io, ws, etc)
- Definir porta/path do WebSocket

## 🏷️ Tipo
feature

## ⏱️ Prioridade
média