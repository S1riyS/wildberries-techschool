# Используем более конкретный alpine-образ для меньшего размера
FROM golang:1.24-alpine AS builder

# Устанавливаем только необходимые зависимости
RUN apk add --no-cache git

# Устанавливаем goose с явным указанием версии (как у вас)
RUN go install github.com/pressly/goose/v3/cmd/goose@v3.24.3

# Финальный образ
FROM alpine:3.22 AS final

# Устанавливаем только bash (без git, так как он больше не нужен)
RUN apk add --no-cache bash

WORKDIR /job

# Копируем goose из builder-этапа
COPY --from=builder /go/bin/goose /usr/local/bin/

# Копируем скрипты и миграции
COPY scripts/migration.sh .
COPY migrations/*.sql migrations/

# Устанавливаем права и указываем entrypoint
RUN chmod +x migration.sh

ENTRYPOINT ["bash", "migration.sh"]