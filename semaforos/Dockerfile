# --- Estágio 1: Compilação ---
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copia o arquivo go.mod que criamos na máquina física
COPY go.mod ./

# Copia o seu arquivo main.go
COPY main.go .

# ATUALIZAÇÃO AQUI: Força o Go a criar o arquivo go.sum e baixar o pacote usando a rede interna modificada do build
RUN go mod tidy

# Compila o binário de forma estática
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# --- Estágio 2: Execução ---
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .

CMD ["./main"]