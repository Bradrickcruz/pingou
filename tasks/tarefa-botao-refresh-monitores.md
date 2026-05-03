# 🧩 Tarefa: Botão de refresh na lista de monitores

## 💡 Ideia
Adicionar botão no header do Dashboard para permitir refresh manual da lista de monitores sem recarregar a página.

## 🎯 Objetivo
Permitir que o usuário atualize os dados dos monitores sob demanda.

## 📦 Escopo inicial
- Adicionar botão de refresh no header do Dashboard
- Utilizar função `refetch` já disponível no hook `useMonitors`
- Indicar visualmente quando o refresh está em andamento

## 🚧 Fora de escopo
- Auto-refresh automático
- Notificações em tempo real (WebSocket)

## 🔍 Contexto
Página Dashboard em web/src/pages/Dashboard.jsx - hook `useMonitors` já expõe `refetch`.

## ⚠️ Observações
Verificar se há ícone de refresh disponível ou usar emoji/texto.

## 🏷️ Tipo
melhoria

## ⏱️ Prioridade
baixa