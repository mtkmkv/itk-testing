# Wallet API

REST API для управления балансом кошельков.

## Стек

- Go 1.25.1
- PostgreSQL 18
- Docker

## Запуск

Запустить приложение:

    make up

Остановить:

    make down

Перезапустить:

    make restart


## API

### Проверка баланса

    curl http://localhost:8080/api/v1/wallets/550e8400-e29b-41d4-a716-446655440000

Пример ответа:

    {
      "walletId": "550e8400-e29b-41d4-a716-446655440000",
      "balance": 1000
    }

### Пополнение кошелька

    curl -X POST http://localhost:8080/api/v1/wallet \
      -H "Content-Type: application/json" \
      -d "{\"walletId\":\"550e8400-e29b-41d4-a716-446655440000\",\"operationType\":\"DEPOSIT\",\"amount\":1000}"

Ожидаемый статус: 200 OK

### Списание средств

    curl -X POST http://localhost:8080/api/v1/wallet \
      -H "Content-Type: application/json" \
      -d "{\"walletId\":\"550e8400-e29b-41d4-a716-446655440000\",\"operationType\":\"WITHDRAW\",\"amount\":1000}"

Ожидаемый статус: 200 OK

### Проверка результата

    curl http://localhost:8080/api/v1/wallets/550e8400-e29b-41d4-a716-446655440000


## Тесты

Запуск всех тестов:

    go test ./...

## Структура проекта

    cmd/server/          # запуск HTTP-сервера
    initial/core/        # domain, router и HTTP server
    initial/features/    # wallet repository, service и handler
    migrations/          # миграции PostgreSQL
    docker-compose.yml   # Docker Compose конфигурация
    Dockerfile           # сборка приложения
    Makefile             # команды запуска