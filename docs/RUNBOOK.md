# Runbook

## Pré-requisitos
Docker Desktop (WSL2), Go 1.27, Kind e kubectl (mesma versão do cluster, com no máximo 1 de diferença).

## Subir do zero
```bash
docker build -t statushub:0.1.0 .
kind create cluster --name statushub
kind load docker-image statushub:0.1.0 --name statushub
kubectl config current-context        # deve ser kind-statushub
kubectl apply -f deploy/k8s/namespace.yaml
kubectl apply -f deploy/k8s/
kubectl rollout status deployment/statushub -n statushub
kubectl port-forward -n statushub svc/statushub 8081:80
```

## Testar
```bash
curl -s http://localhost:8081/readyz
curl -s http://localhost:8081/api/status
```
E no navegador: http://localhost:8081

## Mudar os alvos
Editar `deploy/k8s/configmap.yaml`, depois:
```bash
kubectl apply -f deploy/k8s/configmap.yaml
kubectl rollout restart deployment/statushub -n statushub
```

## Rollback
```bash
kubectl rollout history deployment/statushub -n statushub
kubectl rollout undo deployment/statushub -n statushub
```
Depois, deixe o YAML no Git igual ao que está rodando e aplique de novo.

## Problemas comuns
| Sintoma | Causa provável | O que fazer |
|---|---|---|
| `ErrImagePull` / `ImagePullBackOff` | imagem não carregada no Kind ou tag errada | `kind load docker-image ...` ou `kubectl rollout undo` |
| `CreateContainerConfigError` com `non-numeric user` | `runAsNonRoot` sem UID numérico | `runAsUser: 65532` no `securityContext` |
| Pod `0/1` READY por muito tempo | `/readyz` ainda 503 | `kubectl logs` e `kubectl describe pod` |
| `CrashLoopBackOff` | configuração inválida (fail fast) | `kubectl logs` mostra o motivo |
| `OOMKilled` | passou do limite de memória | revisar `resources.limits` |
| `Readiness/Liveness probe failed ... Client.Timeout` | máquina sobrecarregada ou resposta lenta | ver `kubectl get events`; uma falha isolada não reinicia (precisa de 3) |
| `404 Not Found` em HTML na porta 8080 | outro programa usando a porta 8080 | usar outra porta no port-forward (8081) |
| Aviso de versão do kubectl | kubectl com mais de 1 versão de diferença do cluster | atualizar o kubectl |

## Derrubar
```bash
kind delete cluster --name statushub
```
