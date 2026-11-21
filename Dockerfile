# syntax=docker/dockerfile:1

# ---------- Estágio 1: build ----------
# Imagem com o compilador Go. Ela é grande, mas só serve para compilar.
FROM golang:1.27 AS build
WORKDIR /src

# Copiamos primeiro só go.mod e go.sum e baixamos as dependências.
# Como esses arquivos mudam pouco, o Docker reaproveita esta camada (cache).
COPY go.mod go.sum ./
RUN go mod download

# Agora o resto do código.
COPY . .

# CGO_ENABLED=0 gera um executável estático, que não depende de bibliotecas do sistema.
# -trimpath e -ldflags="-s -w" deixam o binário menor e sem caminhos da sua máquina.
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/statushub ./cmd/statushub

# ---------- Estágio 2: imagem final ----------
# distroless/static: sem shell, sem gerenciador de pacotes; tem certificados HTTPS e o usuário nonroot.
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build /out/statushub /statushub
COPY config/targets.yaml /config/targets.yaml

ENV CONFIG_PATH=/config/targets.yaml
EXPOSE 8080

# Usuário não-root (UID 65532), já criado na imagem distroless.
USER nonroot:nonroot

ENTRYPOINT ["/statushub"]
