# Arquitetura

```mermaid
flowchart LR
    U[Navegador / curl] -->|port-forward 8081:80| S[Service statushub<br>ClusterIP :80]
    S --> P1[Pod 1<br>statushub]
    S --> P2[Pod 2<br>statushub]
    CM[(ConfigMap<br>targets.yaml)] -.montado em /etc/statushub.-> P1
    CM -.-> P2
    P1 -->|GET a cada 30s, timeout 5s| T[Sites e APIs alvo]
    P2 -->|GET a cada 30s, timeout 5s| T
```

## Componentes

| Peça | O que faz |
|---|---|
| `internal/config` | Lê e valida a lista de alvos em YAML (fail fast se estiver inválida) |
| `internal/checker` | Verifica todos os alvos em paralelo (goroutines), com timeout, e guarda o último resultado em memória protegido por `sync.RWMutex` |
| `internal/server` | Rotas HTTP: `/` (página HTML), `/api/status` (JSON), `/healthz` (liveness) e `/readyz` (readiness) |
| `cmd/statushub` | Liga tudo e faz graceful shutdown ao receber SIGTERM |

## Como funciona

- Cada réplica verifica os alvos por conta própria e guarda o último resultado em memória.
- O `/readyz` só responde 200 depois da primeira rodada de verificação; antes disso o Service não manda tráfego para o Pod.
- O `/healthz` nunca depende de nada externo; se ele parar de responder 3 vezes seguidas, o Kubernetes reinicia o container.
- A configuração vem de um ConfigMap montado como arquivo e é lida só na partida.
- Imagem distroless (`static-debian12:nonroot`), rodando como UID 65532, disco somente leitura e sem capabilities.
- Tudo roda num cluster Kind de um nó, dentro do Docker Desktop (WSL2).
