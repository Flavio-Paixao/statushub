# StatusHub

Serviço em Go que verifica, de tempos em tempos, se uma lista de sites e APIs está no ar e mostra o resultado em uma página HTML e em JSON.

Projeto de estudo de **Docker** e **Kubernetes**, rodando 100% local em um cluster **Kind**, com custo zero.

## Situação

🚧 Fase 0: preparação do ambiente.

## Endpoints planejados

| Rota | O que faz |
|---|---|
| `GET /` | Página HTML com cada alvo: no ar ou fora, código HTTP e tempo de resposta |
| `GET /api/status` | O mesmo resultado em JSON |
| `GET /healthz` | Liveness: o processo está vivo |
| `GET /readyz` | Readiness: o serviço está pronto para receber tráfego |

A lista de alvos e o intervalo de verificação vêm de um ConfigMap (arquivo YAML), sem mudar o código.

## Stack

- Go (somente biblioteca padrão)
- Docker (build multi-stage, imagem distroless, usuário não-root)
- Kubernetes local com Kind e kubectl
- GitHub Actions (CI)

## Decisões principais

- **Kind** em vez de Minikube: leve, roda o cluster dentro do Docker e é comum em pipelines de CI.
- **Cada réplica verifica os alvos por conta própria**: mais simples, sem estado compartilhado. Custo: as checagens se multiplicam pelo número de réplicas.
- **Imagem base distroless**: pequena, sem shell, já traz certificados HTTPS e usuário não-root.

O detalhamento (por que, o que foi descartado, do que abri mão) vai ficar em `docs/DECISIONS.md`.

## Estrutura

```
cmd/statushub/   ponto de entrada (main.go)
internal/        código interno do serviço
deploy/k8s/      manifestos do Kubernetes
docs/            documentação
```
