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
