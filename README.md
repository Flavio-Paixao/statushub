# StatusHub

![CI](https://github.com/Flavio-Paixao/statushub/actions/workflows/ci.yml/badge.svg)

Serviço pequeno em **Go** que verifica, de tempos em tempos, se uma lista de sites e APIs está no ar
e mostra o resultado numa página HTML e em JSON. Empacotado com **Docker** (multi-stage, distroless,
não-root) e rodando em **Kubernetes** local (Kind) com 2 réplicas, probes, limites e ConfigMap.

✅ Projeto concluído (Fases 0 a 6).

## O que ele faz

| Rota | O que mostra |
|---|---|
| `GET /` | Página HTML com cada alvo: no ar / fora do ar, código HTTP e tempo de resposta |
| `GET /api/status` | O mesmo em JSON |
| `GET /healthz` | Liveness: "estou vivo" |
| `GET /readyz` | Readiness: 503 até a primeira rodada de verificação, depois 200 |

A lista de alvos, o intervalo (30s) e o timeout (5s) vêm de um arquivo YAML, que no Kubernetes é um ConfigMap.

## Destaques

- Verificação em paralelo com goroutines, timeout por alvo e resultados protegidos por `sync.RWMutex`.
- Graceful shutdown: ao receber SIGTERM, termina os pedidos em andamento antes de sair.
- Imagem Docker de **4,88 MB compactada** (20,4 MB no disco), sem shell, rodando como UID 65532.
- Deployment com 2 réplicas, readiness e liveness probes, requests/limits e `securityContext` restrito (não-root, disco somente leitura, sem capabilities).
- Operação praticada: escalar, autocura, mudança de ConfigMap com `rollout restart` e rollback de um deploy quebrado sem derrubar o serviço.
- CI no GitHub Actions: `gofmt`, `go vet`, `go test -race` e `docker build` a cada push.

## Stack

Go (biblioteca padrão + `go.yaml.in/yaml/v3`) · Docker · Kind · kubectl · GitHub Actions

## Como rodar

Passo a passo completo no [RUNBOOK](docs/RUNBOOK.md). Resumo:

```bash
docker build -t statushub:0.1.0 .
kind create cluster --name statushub
kind load docker-image statushub:0.1.0 --name statushub
kubectl apply -f deploy/k8s/namespace.yaml
kubectl apply -f deploy/k8s/
kubectl port-forward -n statushub svc/statushub 8081:80
```

Depois abra http://localhost:8081.

Sem Kubernetes, só com Go: `go run ./cmd/statushub` e abra http://localhost:8080.

## Documentação

- [Arquitetura](docs/ARCHITECTURE.md): diagrama e componentes.
- [Decisões](docs/DECISIONS.md): o porquê de cada escolha, o que descartei e do que abri mão.
- [Runbook](docs/RUNBOOK.md): subir, testar, mudar alvos, rollback, problemas comuns e derrubar.

## Estrutura

```
cmd/statushub/       ponto de entrada (main.go)
internal/config/     leitura e validação do YAML
internal/checker/    verificador periódico
internal/server/     rotas HTTP
config/targets.yaml  alvos para rodar localmente
deploy/k8s/          namespace, configmap, deployment, service
docs/                arquitetura, decisões e runbook
Dockerfile           build multi-stage com distroless
.github/workflows/   CI
```

## Autor

Flávio Paixão · [flaviopaixao.dev](https://flaviopaixao.dev)
