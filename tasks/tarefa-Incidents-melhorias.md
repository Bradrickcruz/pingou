# 🧩 Tarefa: Melhorias na página de Incidentes

## 💡 Ideia

A página Incidents.jsx atualmente exibe uma lista simples de incidentes com filtro de "open only". O usuário solicitou ideias para melhorar essa página, tornando-a mais útil e amigável.

## 🎯 Objetivo

Transformar a página de incidentes em uma ferramenta mais completa para monitoramento e gestão de problemas nos monitores.

## 📦 Escopo inicial

### Funcionalidades de visualização

- Paginação para lidar com muitos incidentes
- Busca por nome do monitor ou mensagem de erro
- Filtros avançados (por data, status e nome do monitor)

### Detalhamento

- Ao clicar em um incidente, abrir detalhes em modal ou página dedicada
- Mostrar histórico de incidentes do mesmo monitor
- Exibir timeline de eventos do incidente

### Informações adicionais

- Duração do incidente (quanto tempo está aberto)
- Contagem de ocorrências do mesmo monitor
- Indicador de criticidade baseado em tempo

## 🚧 Fora de escopo

- Integração com sistemas de notificação externa
- Regras automáticas de escalação
- Integração com ferramentas de terceiros (PagerDuty, etc.)

## 🔍 Contexto

Página web/src/pages/Incidents.jsx do sistema Pingou Health Checker

## ⚠️ Observações

- O backend já suporta filtragem por "open" - verificar se suporta outros filtros
- Necesário verificar a estrutura atual da API de incidents

## 🏷️ Tipo

melhoria

## ⏱️ Prioridade

média
