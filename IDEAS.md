# 💡 Ideias de evolução

> Só o que ainda **não foi implementado**. O que já entrou está em `docs/` no histórico de implementações.
> Atualizado em 2026-11-04.

**Resumo**

| #   | Ideia                               | Estado                 | Prioridade |
| --- | ----------------------------------- | ---------------------- | ---------- |
| 1   | Backup sem configurações + importar | Aguardando decisão     | medium     |
| 2   | Status HTTP esperado                | Esperando um caso real | low        |
| 3   | Páginas públicas de status          | Para detalhar          | medium     |
| 4   | Modo escuro (dark mode)             | Para detalhar          | medium     |
| 5   | Suporte a inglês (e outros idiomas) | Para detalhar          | medium     |

Mais abaixo: ideias [candidatas](#-candidatas-ainda-não-decididas), [adiadas](#-adiadas) e [descartadas](#-descartadas).

---

# 🧩 Task: Backup sem configurações + importar

## 💡 Idea

Hoje o export (`GET /api/export`, `pingou export-db`) é uma cópia inteira do banco. Isso leva junto o `webhook_url` (que é segredo em Slack/Discord), o tópico do ntfy (que funciona como senha no ntfy.sh), os destinatários de e-mail e os tokens de ping dos heartbeats.

A ideia é o backup levar **só os dados** (monitores, incidentes e, opcionalmente, checks) e **não** as configurações. Para o backup servir de algo, entra junto o caminho de volta: **importar**.

## 🎯 Goal

Um arquivo de backup que possa ser guardado ou compartilhado sem vazar segredos de notificação, e que dê para restaurar numa instância nova.

## 📦 Initial Scope

**Exportação (o que sai no arquivo):**

- Monitores, incidentes e checks. **Fora:** `settings` inteira (webhook, destinatários, tópico do ntfy, fuso, retenção).
- Checks **opcionais**, desmarcados por padrão. Eles são quase todo o tamanho do arquivo (6 GB no perfil limite) e expiram pela retenção. Para migrar de servidor, monitores e incidentes bastam.
- **Token de ping dos heartbeats:** não vai no arquivo. É gerado de novo na importação. Custo: os cron jobs precisarão da URL nova depois de restaurar.
- Janelas de manutenção (`maintenance_windows`): não exportar (têm hora marcada e perdem o sentido num backup).
- A tela de Configurações avisa claramente: "Não inclui configurações de notificação".
- `pingou export-db` segue o mesmo formato.

**Importação:**

- Primeiro um comando `pingou import <arquivo>` no terminal, que **restaura numa instância vazia**. Não passa pelo navegador, então não tem limite de upload, e o scheduler não está rodando.
- Confere a versão das migrations do arquivo antes de importar e recusa arquivo de versão incompatível.
- Importa tudo numa transação só (tudo ou nada).
- Tela de importar em Configurações: **só depois**, e só se fizer falta.

## 🚧 Out of Scope

- Cifrar o arquivo (ver Notes).
- Mesclar com dados já existentes (importar em cima de uma instância com monitores).
- Importar pela tela, no primeiro momento.
- Agendar backups automáticos.

## 🔍 Context

Configurações → Backup. O `VACUUM INTO` copia o banco todo, então para tirar as configurações é preciso gerar um banco novo só com as tabelas certas (ou apagar e rodar `VACUUM` de novo, porque dado apagado continua no arquivo até o `VACUUM`). Em banco grande isso custa tempo e disco, e precisa ser medido.

## ⚠️ Notes

- **O que ainda vaza sem cifrar:** URLs de monitor com token na query (`?token=...`). O plano deve avisar isso na tela e no README. Só cifrar resolve de verdade.
- **Cifrar o arquivo (se um dia fizer falta):**
  - A chave tem de ser uma **senha digitada na hora**, nunca a API Key (trocar a API Key tornaria os backups antigos ilegíveis).
  - Cifrar protege o arquivo parado. **Não** protege contra vazamento da API Key: o atacante chama o export e escolhe a própria senha. Quem protege esse caso é tirar as configurações.
  - Senha perdida significa backup perdido.
  - Arquivo de vários GB exige criptografia em blocos. Escrever à mão é arriscado. O caminho seguro é uma biblioteca pronta (como a `age`), que é dependência nova.
- **Decisões que faltam** (por isso a ideia está "Aguardando decisão"):
  1. O uso principal é recuperar de um desastre ou migrar de servidor?
  2. Aceita regerar os tokens dos heartbeats ao importar? Manter os tokens exige cifrar.
- **Registrar na limitação conhecida** do ntfy e do e-mail que o export deixa de levar esses dados quando esta ideia entrar.

## 🏷️ Type

feature (segurança)

## ⏱️ Priority

medium (começar pela exportação sem configurações, que é pequena)

---

# 🧩 Task: Status HTTP esperado

## 💡 Idea

Hoje o monitor HTTP decide "no ar" ou "fora do ar" por uma regra fixa de status. Em alguns casos o resultado correto é outro: uma rota protegida que devolve `401` quando está saudável, um endpoint que devolve `204`, um redirect `301` que não deve ser seguido.

A ideia é cada monitor HTTP poder declarar **quais status contam como sucesso**.

## 🎯 Goal

Monitorar corretamente endpoints cujo "saudável" não é um `2xx` comum, sem criar falso alarme.

## 📦 Initial Scope

- Campo opcional "Status esperados" no formulário HTTP (seção avançada). Aceita lista e faixas, por exemplo `200-299,401`.
- Sem o campo preenchido, o comportamento atual não muda.
- Validar na API: só números de 100 a 599, tamanho máximo, faixa com início menor que o fim (`400` com `VALIDATION_ERROR`).
- A mensagem de erro (`unexpected status code: X`) passa a falar do esperado.
- Mostrar o esperado na página de detalhes do monitor.

## 🚧 Out of Scope

- Palavra-chave no corpo da resposta.
- Cabeçalhos e autenticação customizados.
- Qualquer coisa fora de monitores HTTP.

## 🔍 Context

Checker HTTP e formulário de monitor. Na migration `0003`, a coluna `config` em `monitors` guarda JSON por tipo. Verificar se o campo cabe lá, para **evitar uma migration nova**.

## ⚠️ Notes

- **Esperar um caso real antes de fazer.** Sem um caso concreto, a lista de faixas corre o risco de ficar mais complexa do que o necessário.
- Se aparecer o caso real, avaliar junto a **palavra-chave no corpo**: normalmente os dois pedidos nascem do mesmo problema ("o site responde 200 mas está quebrado").
- Verificar no checker o que é considerado sucesso hoje (o redirect segue até `PINGOU_MAX_REDIRECTS`), para a regra nova não mudar monitores existentes.

## 🏷️ Type

feature

## ⏱️ Priority

low

---

# 🧩 Task: Páginas públicas de status

## 💡 Idea

Criar páginas públicas, sem login, em que o dono escolhe **quais monitores** aparecem. Quem abre a página só **lê**: não altera nada.

## 🎯 Goal

Mostrar a clientes e usuários se os serviços estão no ar, sem dar acesso ao painel nem à API Key.

## 📦 Initial Scope

**Dados exibidos por monitor** (todos somente leitura):

- Status atual (no ar, fora do ar, desconhecido, em manutenção)
- Uptime (%) por período
- Incidentes (início, fim e duração)
- Tempo fora do ar
- Validação do SSL (só para HTTP com `https://`)
- Checks recentes: **só o status**, sem tempo de resposta
- Tipo do monitor (http, tcp, heartbeat)
- "Verifica a cada" (intervalo)

**Configuração (no painel, atrás do login):**

- Criar várias páginas, cada uma com nome, endereço próprio e a lista de monitores escolhidos.
- Ligar e desligar cada página sem apagar.
- Mostrar o link público para copiar.

**Páginas públicas:**

- Rota própria fora do login (ex.: `/status/<endereço>`), servida pelo mesmo binário.
- Endpoints públicos próprios (ex.: `/api/public/status/<endereço>`), separados da API protegida, e só de leitura (`GET`).
- Mobile e desktop, seguindo o `DESIGN.md`.

## 🚧 Out of Scope

- Assinar atualizações por e-mail ou SMS
- Domínio próprio e marca personalizada
- Comentários nos incidentes ou histórico escrito pelo dono
- Qualquer escrita feita por quem visita a página
- Gráfico de tempo de resposta (fica de fora de propósito)

## 🔍 Context

O README já lista "status page pública" em Pós-MVP. Hoje toda a API fica atrás da API Key. A página pública precisa de rotas que **não passem** pelo middleware de login, e é o primeiro dado do Pingou exposto a quem não está autenticado.

## ⚠️ Notes

- **Exposição é a parte delicada.** Montar a resposta pública com uma lista fixa de campos permitidos (lista de permissão), nunca reaproveitando os objetos da API protegida. Assim um campo novo no futuro não vaza sozinho.
- **Nunca expor:** URL ou alvo do monitor, token de ping do heartbeat, mensagem de erro (`last_error` costuma ter nome de host e detalhe interno) e as configurações.
- **Nome do monitor** é público. Avisar na tela de configuração, ou permitir um "nome público" separado.
- **Endereço da página:** decidir entre endereço legível (`/status/clientes`) ou aleatório (não adivinhável). Página "pública" com endereço aleatório não é privada, só é difícil de achar.
- **Carga:** a página pública é aberta por qualquer um, e o SQLite tem uma conexão só. Precisa de cache curto (por exemplo 30 a 60 s) e limite de requisições por IP, para uma onda de acessos não atrasar a gravação dos checks. O `/stats` de 30 dias é a consulta mais cara.
- **Validação do SSL:** decidir o que mostrar: só "válido" ou a data de expiração.
- **Manutenção:** decidir se uma janela aparece como "Em manutenção" na página pública (sugiro que sim, com o motivo oculto).
- **Decisões tomadas (2026-11-04):**
  1. **Várias páginas**, cada uma com seus monitores. Um monitor pode estar em mais de uma página.
  2. **Endereço legível**, digitado por quem cria a página, e **único**. Regras sugeridas: letras minúsculas, números e hífen, de 3 a 40 caracteres; recusar palavras reservadas (`api`, `admin`, `login`, `assets` etc.) e endereço repetido (`409`). Depois de criado, mudar o endereço quebra o link antigo, então avisar na tela.
  3. **Uptime em 24 h, 7 d e 30 d.** Respeitar a retenção: período acima de `retention_days` aparece desabilitado, como no detalhe do monitor.
  4. **Mostrar o motivo da manutenção.**
- **Cuidado com o motivo público:** o motivo é texto livre (até 200 caracteres) e pode conter detalhe interno ("migrando o servidor X"). Quando o monitor estiver numa página pública, o modal de pausa e o de agendar manutenção devem avisar: "este motivo aparece na página pública". O mesmo vale para o motivo de fechamento de incidente, se for exibido. 5. **SSL: só "Válido" e "Expira em breve".** A data exata não aparece na página pública. "Expira em breve" vale quando faltam 30 dias ou menos, o mesmo limite do aviso de SSL do painel. Só para monitor HTTP com `https://` e com data de certificado já conhecida; sem data, a linha não aparece.

## 🏷️ Type

feature (segurança)

## ⏱️ Priority

medium

---

# 🧩 Task: Modo escuro (dark mode)

## 💡 Idea

Oferecer um tema escuro no painel, além do tema claro atual.

## 🎯 Goal

Quem usa o painel à noite, ou prefere tema escuro, consegue trocar sem perder legibilidade nem a identidade visual do Pingou.

## 📦 Initial Scope

- Três opções: **Claro**, **Escuro** e **Automático** (segue o sistema, via `prefers-color-scheme`). Padrão: Automático.
- A escolha vale **por navegador** (guardada no `localStorage`), não por instalação. Diferente do fuso, que é global: aqui cada pessoa escolhe o seu.
- Seletor de tema nas Configurações e, se couber, um atalho discreto no cabeçalho.
- Aplicar o tema **antes** da primeira pintura, para não piscar o claro e depois virar escuro.
- Todas as telas: Painel, Incidentes, Detalhe, Manutenções, Configurações, login e modais, no desktop e no celular.
- Gráficos SVG (resposta, pings, faixas "fora do ar") e selos de estado (UP, DOWN, UNKNOWN, MAINTENANCE) com cores próprias do escuro.
- Registrar a paleta escura e as regras de uso no `DESIGN.md`.

## 🚧 Out of Scope

- Temas personalizados ou escolha de cor
- Tema diferente por monitor ou por página
- Mudar o tema claro atual

## 🔍 Context

O painel foi redesenhado em tema claro (v1.1). As cores já devem sair de tokens do `DESIGN.md`. O `DESIGN.md` §3.4 diz que cor nova precisa de aprovação, então a paleta escura é uma decisão de design, não só de código.

## ⚠️ Notes

- **Primeiro passo:** levantar quantas cores estão fixas no código (hexadecimais soltos, classes de cor do Tailwind fora dos tokens). Cada uma vira token antes de criar o tema, senão o escuro fica "furado".
- **Contraste:** conferir o mínimo de legibilidade (WCAG AA) nos dois temas, principalmente em texto cinza, selos e linhas de gráfico.
- Vermelho (DOWN) e verde (UP) precisam continuar distinguíveis no escuro, também para quem tem daltonismo; não depender só da cor (o texto do selo já ajuda).
- **Páginas públicas de status** (ideia 3) devem herdar o mesmo sistema de temas. Quem fizer uma antes da outra deve criar as telas novas já com tokens.
- **Decisões tomadas (2026-11-04):**
  1. **A paleta escura é aprovada pelo dono do projeto** (Bryan). O código só começa a usar o escuro depois da aprovação.
  2. **Sem atalho no cabeçalho do painel.** Para trocar o tema, vai em Configurações.
  3. **A página pública tem botão de tema no cabeçalho**, para o visitante escolher como quiser (Claro, Escuro, Automático). A escolha fica no navegador do visitante e não afeta o painel nem os outros visitantes.
- **Ajuste no escopo:** o "atalho no cabeçalho" do Initial Scope vale **só para a página pública**. No painel, o seletor fica apenas em Configurações.

## 🏷️ Type

improvement

## ⏱️ Priority

medium

---

# 🧩 Task: Suporte a inglês (e outros idiomas)

## 💡 Idea

Hoje o painel e as mensagens estão em português. A ideia é suportar **inglês**, já com a estrutura pronta para adicionar outros idiomas sem mexer em código de tela.

## 🎯 Goal

Quem não lê português consegue usar o Pingou, e adicionar um idioma novo vira só adicionar um arquivo de traduções.

## 📦 Initial Scope

- Idiomas iniciais: `pt-BR` (atual) e `en`.
- **Textos fora do código das telas:** um arquivo de traduções por idioma (por exemplo, `web/src/locales/pt-BR.json` e `en.json`), com chaves estáveis. Adicionar um idioma = criar um arquivo novo e registrá-lo numa lista.
- **Idioma padrão:** o do navegador na primeira abertura (como foi feito com o fuso), com fallback para `en` (ou `pt-BR`, a decidir).
- **Seletor de idioma** nas Configurações.
- **Datas, horas e números** pelo idioma escolhido (`Intl`), respeitando o fuso das Configurações. Rever as datas formatadas à mão, que foram feitas assim porque o `Intl` em pt-BR gerava "28 de out.".
- **Plurais e textos com variáveis** ("Fora do ar há 1 dia e 4 h", "3 monitores", "Expira em 12 dias") com regra de plural por idioma, sem juntar frases no código.
- **Notificações** (e-mail e ntfy) no idioma configurado. O JSON do webhook **não muda** (continua `event`, `monitor`, `timestamp`...).
- **Mensagens de erro da API:** hoje misturam inglês (validação) e português. Caminho sugerido: a API devolve **código estável** (`VALIDATION_ERROR` etc., já existe) e o painel traduz a mensagem pelo código.

## 🚧 Out of Scope

- Traduzir a documentação (README, `DESIGN.md`)
- Idiomas que escrevem da direita para a esquerda (árabe, hebraico)
- Tradução automática ou serviço externo de traduções
- Idioma por usuário na API (só uma API Key)

## 🔍 Context

Frontend React + Vite, backend Go. O Pingou segue a regra "sem dependência externa obrigatória", então a escolha entre uma biblioteca (como `react-i18next`) e uma solução mínima própria é uma decisão de projeto.

## ⚠️ Notes

- **Pesa mais no backend do que parece:** notificações (e-mail e ntfy), mensagens de erro da API, selos de motivo ("Fechado por manutenção") e motivos de erro do checker (`context deadline exceeded...`) chegam ao usuário. Decidir o que traduz e o que fica como está (erro técnico do checker, em inglês, é aceitável).
- **Onde guardar o idioma das notificações:** uma configuração global em `settings` (chave-valor, sem migration), como o fuso. O idioma do painel, por outro lado, pode ser por navegador. Isso significa duas coisas: **idioma do painel** (por navegador) e **idioma das notificações** (global). Decidir se são a mesma configuração ou duas.
- **Estratégia para não ficar para trás:** teste que compara as chaves de todos os idiomas e falha se faltar alguma. Fallback para o idioma padrão quando uma chave não existe.
- **Fazer cedo é mais barato:** cada tela nova (páginas públicas, importar etc.) escrita em português fixo precisa ser refeita. Se a ideia 3 for feita antes, nascerá com textos soltos. Recomendo planejar o i18n **antes** das páginas públicas, ou garantir que as telas novas já usem o arquivo de traduções.
- **Páginas públicas:** o idioma da página deve seguir o navegador do visitante, não a configuração do dono.
- **Decisões tomadas (2026-11-04):**
  1. **Biblioteca pronta** (como `react-i18next`). É uma dependência de frontend nova e aceita. Escolher a biblioteca no plano, olhando tamanho do bundle.
  2. **Uma configuração só de idioma**, não duas. Assumi que é **global** (uma por instalação, em `settings`, como o fuso), valendo para o painel e para as notificações. Assim painel e e-mail ficam no mesmo idioma. _(Confirmar essa leitura.)_
  3. **Fallback: `en`.** Chave que falta num idioma cai para o inglês. Idioma do navegador que não existe no Pingou também cai para `en`.
- **Consequência da decisão 2:** com idioma global, a primeira pessoa a abrir o painel define o idioma de todos (mesma limitação já conhecida do fuso). Mantém a regra "idioma do navegador na primeira abertura, uma vez" e depois fica salvo.

## 🏷️ Type

feature

## ⏱️ Priority

medium (fazer antes das páginas públicas)

---

# 🌱 Candidatas (ainda não decididas)

Surgiram no histórico de implementações. Ninguém decidiu fazer. Estão aqui para não se perder.

- **RAM acima da meta no cenário grande.** No perfil limite (100 monitores × 10 s, 25 milhões de checks), a primeira subida depois da migration usa 47 a 51 MB, e picos de 50 a 60 MB já foram vistos; a meta do README é menos de 50 MB. Opções: investigar e baixar, ou ajustar o texto do README para dizer em que cenário a meta vale. O `memtest` de 1 hora também não foi rodado inteiro.
- **Retry de notificações (`notification_outbox`).** Hoje um aviso perdido (e-mail, webhook ou ntfy falhou, ou o processo caiu) não é reenviado. Pesa mais agora, porque o aviso de SSL sai só duas vezes por certificado. Caminho: tabela gravada na transação do `Process` + um worker.
- **Baldes do `/stats` em fusos de meia hora** (ex.: `Asia/Kolkata`). Hoje alinhados em UTC, então os rótulos de 1 h não caem em hora cheia.
- **Janela de manutenção recorrente** ("todo domingo às 2h"). Ficou fora da v1.
- **Aviso de "a manutenção começou/terminou"** nos canais. Também ficou fora da v1.

# ⏸️ Adiadas

- **Proxies confiáveis** (`PINGOU_TRUSTED_PROXIES`): rate limit por IP real quando o Pingou roda atrás de proxy.
- **Importar e exportar só monitores** em formato legível (JSON), para compartilhar configuração entre instâncias. Parte disso deve cair na ideia 1; avaliar depois dela.

# 🗑️ Descartadas

- **Bloquear IPs privados no webhook e no teste** (descartada em 2026-11-04). Já tinha sido recusada antes, porque quebraria webhooks para o n8n em rede local.
- **Tags nos monitores.** A busca parcial resolveu.
- **Chave de API com rotação.**
