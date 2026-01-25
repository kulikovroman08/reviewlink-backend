FROM golang:1.24-alpine3.21 AS build

RUN apk update && apk add --no-cache git ca-certificates

# Сборка бэка
WORKDIR /build/backend

# go.mod лежит в backend/
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# копируем исходники backend
COPY backend/ ./

# копируем фронт отдельно в /build/frontend
COPY frontend/ /build/frontend/

# Собираем бинарь сервера
RUN CGO_ENABLED=0 go build -p 1 -trimpath -ldflags="-w -s" -o /build/server ./cmd

# Final image
FROM alpine:3.21

RUN adduser -D -H -u 10001 app && \
    mkdir -p /app /home/app && \
    chown -R app:app /app /home/app

ENV HOME=/home/app

COPY --from=build /build/server /app/server
COPY --from=build /build/frontend /app/frontend

USER app
WORKDIR /app

EXPOSE 8080
CMD ["/app/server"]

