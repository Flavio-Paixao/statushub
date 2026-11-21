# Decisões

## 1. Kind em vez de Minikube
- **Por quê:** os nós do Kind são containers dentro do Docker Desktop que eu já usava, então ele é leve para um notebook de 8 GB, sobe e apaga um cluster em cerca de 1 minuto, e o `kind load docker-image` coloca a imagem local no cluster sem precisar de um registry.
- **O que descartei:** Minikube, que pede um driver ou VM a mais e traz addons que eu não usaria.
- **Do que abri mão:** não tem LoadBalancer nem Ingress prontos; para acessar o serviço uso `kubectl port-forward`.

## 2. Cada réplica verifica os alvos por conta própria
- **Por quê:** é o desenho mais simples. Não tem estado compartilhado nem componente extra, cada Pod é independente e a readiness de cada um reflete os próprios dados.
- **O que descartei:** uma verificadora única com estado compartilhado (Redis).
- **Do que abri mão:** com N réplicas, cada alvo recebe N verificações por intervalo, e dois Pods podem mostrar resultados de momentos um pouco diferentes.

## 3. Imagem base distroless (static, nonroot)
- **Por quê:** traz só o necessário (certificados HTTPS, fuso horário e o usuário nonroot), não tem shell nem gerenciador de pacotes e deixou a imagem com 4,88 MB compactada (20,4 MB no disco).
- **O que descartei:** scratch (sem certificados HTTPS) e alpine (shell e pacotes a mais).
- **Do que abri mão:** não dá para abrir shell no container para depurar (`docker exec ... sh` falha). E como o usuário `nonroot` tem nome e não número, o Deployment precisa de `runAsUser: 65532` para o `runAsNonRoot` funcionar.

## 4. YAML com go.yaml.in/yaml/v3 (exceção ao "só biblioteca padrão")
- **Por quê:** YAML aceita comentários, é fácil de ler e é o mesmo formato dos manifestos e do ConfigMap do Kubernetes.
- **O que descartei:** JSON com `encoding/json` (sem comentários e menos legível para editar à mão).
- **Do que abri mão:** uma dependência externa, travada no `go.sum` (v3.0.6).

## 5. Configuração lida só na partida
- **Por quê:** é simples, não precisa observar o arquivo nem sincronizar troca de configuração em tempo real, e configuração inválida derruba o serviço logo na partida (fail fast).
- **Do que abri mão:** mudar o ConfigMap exige `kubectl rollout restart` (na Fase 5 o log continuou dizendo 3 alvos até o restart).
