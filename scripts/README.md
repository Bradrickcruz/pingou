# scripts/

## memtest: prova da promessa de memória

O [PRD](../PRD.md) promete que o Pingou roda com **menos de 50 MB de RAM** com 100 monitores checando a cada 10 s (e sem vazamento em execuções longas). O `memtest` mede isso.

### Como usar

```sh
make build                  # ou: CGO_ENABLED=1 go build -o bin/pingou ./cmd/pingou
scripts/memtest.sh          # 1 h (padrão)
scripts/memtest.sh 86400    # 24 h, para olhar tendência de vazamento
HTTPS_N=50 scripts/memtest.sh 3600
```

Requisitos: `bin/pingou` compilado **na versão atual** (um binário velho sobe na porta errada), porta `8956` livre (o Pingou usa essa porta fixa), porta `9000` livre, `go`, `curl` e Linux (lê `/proc`).

Saída de exemplo:

```
monitores google.com: 30
101º monitor -> HTTP 409 (esperado 4xx)
14:54:12 rss_kb=23740 peak_kb=24704 threads=16
...
pico: 24 MB (limite 50)
PASSOU
```

O script sai com código `1` (`FALHOU`) se o pico passar de 50 MB. O log completo fica em `/tmp/tmp.*/mem.log`.

### O que ele faz

1. Compila e sobe o alvo local ([memtest/target.go](memtest/target.go), porta `9000`) com três rotas: `/ok/*` (200), `/slow/*` (200 após 2 s) e `/fail/*` (500). A mesma porta serve de alvo TCP.
2. Sobe o Pingou com um **banco temporário** e uma API Key de teste. O seu `pingou.db` não é tocado.
3. Cria 100 monitores de 10 s:

   | Tipo           | Quantidade     | Alvo                             |
   | -------------- | -------------- | -------------------------------- |
   | HTTPS          | 30 (`HTTPS_N`) | `https://google.com`             |
   | HTTP ok        | 40             | `/ok/*` local                    |
   | HTTP lento     | 10             | `/slow/*` (2 s)                  |
   | HTTP com falha | 10             | `/fail/*` (500, abre incidentes) |
   | TCP            | 10             | porta 9000 local                 |

4. Confirma que o 101º monitor é recusado (limite do produto).
5. A cada 10 s, simula um usuário com o dashboard aberto (`GET /api/monitors` e `/api/incidents`) e registra `VmRSS`, `VmHWM` (pico) e número de threads de `/proc/<pid>/status`.
6. No fim, compara o **pico** com 50 MB.

### Como interpretar

- **Passa** se o pico fica abaixo de 50 MB e o `rss_kb` forma um platô. Em testes longos, uma linha de RSS que sobe sem parar é vazamento, mesmo abaixo de 50 MB.
- O RSS já inclui o cache de páginas do SQLite e o front embutido, que é o que o PRD quer limitar.

### Limitações

- **HTTPS só contra o Google.** Alvos locais com certificado autoassinado seriam recusados pelo checker, então o TLS é exercitado só pelos monitores `google.com`. Eles geram cerca de 3 requisições por segundo para um serviço externo. Em testes longos o Google pode devolver 429 ou bloquear o IP, o que vira monitor DOWN e incidente falso. Reduza `HTTPS_N` se acontecer.
- **Goroutines não são contadas.** O projeto não expõe `pprof`; `threads` é só um indicador indireto. A promessa de "sem vazamento de goroutines" ainda não é coberta.
- **Sem webhook configurado.** Os incidentes abrem, mas nenhuma notificação é enviada, então o custo de memória dos envios não entra na medição.
- **Sem limite de cgroup.** O teste mede o processo, não impõe 50 MB. Para a prova mais rigorosa, rode o container com `docker run --memory=50m` e veja se ele é morto por OOM.
- Rodar enquanto outro Pingou usa `~/.pingou/pingou.lock` pode falhar, porque o Pingou impede duas instâncias na mesma máquina.

## dbbench: desempenho do banco com volume real

Mede as consultas do painel, a gravação dos checks e a limpeza de retenção num banco com 30 dias de dados. Serve para comparar **antes e depois** de uma mudança no banco (índice, consulta, migration), sempre na mesma máquina.

### Como usar

```sh
scripts/dbbench.sh            # perfil típico: 10 monitores a cada 60 s (~400 mil checks, ~100 MB)
scripts/dbbench.sh limite     # perfil limite: 100 monitores a cada 10 s (~25 milhões de checks, ~6 GB)
```

O script:

1. Gera o banco em `.bench/<perfil>.db` se ele não existir (`go run ./scripts/seed`). O perfil limite leva uns 5 minutos.
2. Imprime o plano de cada consulta quente (`EXPLAIN QUERY PLAN`).
3. Roda os benchmarks numa **cópia** do banco. O `.bench/<perfil>.db` nunca é alterado. A cópia aplica as migrations pendentes, como o app ao subir, e o relatório mostra quanto tempo isso levou.
4. Salva tudo em `.bench/results-<perfil>-<data>.md`.

Requisitos: `go` e espaço livre de cerca de **2× o tamanho do banco** (o original e a cópia). O perfil limite precisa de ~15 GB livres. A pasta `.bench/` fica fora do git e da imagem Docker.

Variáveis opcionais:

| Variável                   | Padrão | O que faz                                   |
| -------------------------- | -----: | ------------------------------------------- |
| `PINGOU_BENCH_REPS`        |     20 | Repetições de cada consulta                 |
| `PINGOU_BENCH_WRITES`      |   5000 | Gravações no teste de custo de escrita      |
| `PINGOU_BENCH_CONCURRENCY` |     30 | Segundos do cenário leitura + gravação      |

Para rodar só uma parte:

```sh
go run ./scripts/seed --profile tipico --out .bench/tipico.db
PINGOU_BENCH_DB=.bench/tipico.db go test ./internal/repository -run TestPrintQueryPlans -v
PINGOU_BENCH_DB=.bench/tipico.db go test ./internal/repository -run TestDBBench -v -timeout 2h
```

Sem `PINGOU_BENCH_DB` esses testes são pulados, e o `make test` continua rápido.

### O que é medido

- **Tamanho:** arquivo e cada índice.
- **`/stats`:** os 4 períodos, em um monitor HTTP, um TCP e um Heartbeat.
- **Lista do painel:** `GET /api/monitors`, com histórico, uptime 24 h e última latência.
- **Checks paginados**, **incidentes**, **busca do heartbeat** e do **token de ping**.
- **Gravação** pelo caminho real (state machine): check com sucesso e com falha.
- **Concorrência:** `/stats` 30d sem parar e uma gravação a cada 100 ms. Como o SQLite usa uma conexão só, uma leitura lenta atrasa a gravação.
- **Retenção:** apaga o dia mais antigo com gravações ao lado. A espera máxima de uma gravação mostra por quanto tempo a limpeza segura a conexão.

Cada linha mostra a primeira execução (cache frio para aquela consulta), a mediana (p50), o p95 e o máximo.

### Os dados

O gerador cria os dados pelas migrations reais, em **ordem cronológica e com os monitores intercalados**, como o scheduler grava. Dados agrupados por monitor deixariam as consultas artificialmente rápidas.

- Latência de ~20 a ~400 ms; ~0,3 % de falhas isoladas.
- Uma queda de 5 a 30 min por monitor por semana, com incidente aberto e fechado.
- Incidentes de até 1 ano atrás (incidentes nunca expiram).
- Monitores HTTP, TCP e Heartbeat (ping a cada 5 min).

É determinístico: o mesmo `--end` gera o mesmo banco. O `--end` padrão é hoje às 00:00 UTC, para os dados terminarem perto de "agora".

### Como interpretar

Os números dependem do hardware. Compare só rodadas feitas na mesma máquina, sem outra carga pesada ao lado. Rode mais de uma vez: uma diferença menor que ~15 % entre rodadas iguais é ruído.
