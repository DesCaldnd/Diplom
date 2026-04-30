# ==============================================================================
# СТАДИЯ 1: Сборка (Builder)
# ==============================================================================
FROM golang:1.26.1-alpine AS builder

# Установка необходимых инструментов для сборки
RUN apk add --no-cache \
    git \
    protobuf \
    protobuf-dev

# Установка protoc-gen-go и protoc-gen-go-grpc
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.11 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

WORKDIR /app

# Копируем файлы зависимостей для кэширования
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код и генерируем gRPC код
COPY . .
RUN go generate ./...

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/diplom

# ==============================================================================
# СТАДИЯ 2: Релиз (Runtime) - легкий образ
# ==============================================================================
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем бинарник из стадии сборки
COPY --from=builder /app/main .

EXPOSE 9999

CMD ["./main"]
