# Auth Service

Микросервис аутентификации на Go с использованием PostgreSQL.

## Описание

Этот сервис реализует регистрацию, вход, выход, обновление токенов и получение информации о пользователе. Использует JWT для токенов, хранит сессии и пользователей в PostgreSQL.

## Стек

- Go
- PostgreSQL
- Docker/Docker Compose

## Запуск

1. Клонируйте репозиторий:
    ``sh
    git clone https://github.com/FooxyS/auth-service.git
    cd auth-service
    ``

2. Создайте файл .env на основе .env.example и заполните переменные окружения.


3. Соберите и запустите сервис через Docker Compose:
    docker-compose up --build

Сервис будет доступен на http://localhost:8088.

## API
    POST /register — регистрация пользователя
    POST /login — вход пользователя
    POST /logout — выход пользователя
    POST /refresh — обновление токенов
    GET /me — информация о текущем пользователе

## Тесты
Для запуска unit-тестов:
    go test ./internal/usecase/

## Структура проекта
- cmd/main.go — точка входа
- internal/adapters/http — HTTP-обработчики и роутер
- internal/usecase — бизнес-логика
- internal/domain — доменные сущности и интерфейсы
- internal/infrastructure — работа с БД, хэширование, токены
- pkg — вспомогательные пакеты (ошибки, константы)
