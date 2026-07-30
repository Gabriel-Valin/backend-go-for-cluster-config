# --- estágio de build ---
FROM golang:1.23-alpine AS build
WORKDIR /src

# cache de dependências: copia go.mod primeiro
COPY go.mod ./
RUN go mod download

COPY . .

# a versão é injetada via ldflags no momento do build
ARG VERSION=dev
RUN CGO_ENABLED=0 go build -ldflags "-X main.version=${VERSION}" -o /app .

# --- estágio final ---
FROM gcr.io/distroless/static-debian12
COPY --from=build /app /app
USER nonroot:nonroot          # roda como não-root (satisfaz o PSS 'restricted' da Parte 5)
EXPOSE 8080
ENTRYPOINT ["/app"]
