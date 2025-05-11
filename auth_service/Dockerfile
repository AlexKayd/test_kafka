# Сборочный контейнер
FROM golang:1.21 AS builder

WORKDIR /app

COPY . .
ENV GOPROXY=https://proxy.golang.com.cn,direct

RUN go mod tidy
RUN go build -o app .

# Финальный контейнер на Alpine
FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/app .
COPY --from=builder /app/app.exe .

EXPOSE 8090

CMD ["./app"]