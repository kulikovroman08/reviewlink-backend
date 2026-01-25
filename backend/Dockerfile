# syntax=docker/dockerfile:1

FROM golang:1.24-alpine3.21 AS build

RUN apk update && apk add --no-cache git ca-certificates

WORKDIR /build

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY . .

# Собираем один бинарь сервера (у тебя cmd/main.go)
RUN CGO_ENABLED=0 go build -p 1 -trimpath -ldflags="-w -s" -o ./server ./cmd

# Final image
FROM alpine:3.21

RUN adduser -D -H -u 10001 app && \
    mkdir -p /app /home/app && \
    chown -R app:app /app /home/app

ENV HOME=/home/app

COPY --from=build /build/server /app/server

USER app
WORKDIR /app

EXPOSE 8080
CMD ["/app/server"]

