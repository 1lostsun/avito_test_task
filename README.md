# 🚀 Avito Test Task

Система управления Pull Request'ами с поддержкой команд, пользователей и ролевой модели доступа.

## 📋 Содержание

- [Описание](#описание)
- [Технологии](#технологии)
- [Архитектура](#архитектура)
- [Быстрый старт](#быстрый-старт)
- [API](#api)
- [Мониторинг](#мониторинг)
- [Разработка](#разработка)

## 📖 Описание

Приложение для управления Pull Request'ами в командах разработки. Система позволяет:

- Управление командами и пользователями
- Создание и управление Pull Request'ами
- Автоматическое назначение ревьюеров
- Ролевая модель доступа (Admin/User)
- Интеграция с HashiCorp Vault для хранения секретов
- Мониторинг через Prometheus и Grafana

## 🛠 Технологии

- **Backend**: Go 1.24
- **Framework**: Gin
- **База данных**: PostgreSQL
- **Секреты**: HashiCorp Vault
- **Мониторинг**: Prometheus + Grafana
- **Контейнеризация**: Docker + Docker Compose

## 🏗 Архитектура

```
┌─────────────┐
│   Client    │
└──────┬──────┘
       │
┌──────▼──────────────────┐
│   HTTP Server (Gin)      │
│   - Handlers             │
│   - Middleware           │
└──────┬───────────────────┘
       │
┌──────▼──────────────────┐
│   UseCase Layer         │
│   - Business Logic      │
└──────┬───────────────────┘
       │
┌──────▼──────────────────┐
│   Repository Layer       │
│   - PostgreSQL          │
└──────┬───────────────────┘
       │
┌──────▼──────────────────┐
│   PostgreSQL Database   │
└──────────────────────────┘
```

## 🚀 Быстрый старт

### Предварительные требования

- Docker и Docker Compose
- Make (опционально)

### Установка и запуск

1. **Клонируйте репозиторий** (если еще не сделано)

2. **Запустите все сервисы одной командой:**

```bash
make up
```

Или без Make:

```bash
cd infra
docker-compose up -d
```

3. **Проверьте статус сервисов:**

```bash
make status
```

### Доступ к сервисам

После запуска сервисы будут доступны по следующим адресам:

| Сервис | URL | Учетные данные |
|--------|-----|----------------|
| **Приложение** | http://localhost:8080 | - |
| **Vault UI** | http://localhost:8201 | Token: `root` |
| **Prometheus** | http://localhost:9090 | - |
| **Grafana** | http://localhost:3000 | `admin` / `admin` |
| **PostgreSQL** | localhost:5432 | См. Vault |

## 📡 API

### Teams (Команды)

#### Создать команду
```http
POST /avito-test-task/team/add
token: admin
Content-Type: application/json

{
  "name": "backend-team",
  "members": ["user1", "user2"]
}
```

#### Получить команду
```http
GET /avito-test-task/team/get?name=backend-team
token: admin/user
```

### Users (Пользователи)

#### Активировать/деактивировать пользователя
```http
POST /avito-test-task/users/setIsActive
token: admin
Content-Type: application/json

{
  "user_id": "user1",
  "is_active": true
}
```

#### Получить ревью пользователя
```http
GET /avito-test-task/users/getReview?user_id=user1
token: admin/user
```

### Pull Requests

#### Создать Pull Request
```http
POST /avito-test-task/pullRequest/create
token: admin
Content-Type: application/json

{
  "id": "pr-123",
  "author_id": "user1",
  "team_name": "backend-team"
}
```

#### Смержить Pull Request
```http
POST /avito-test-task/pullRequest/merge
token: admin
Content-Type: application/json

{
  "pr_id": "pr-123"
}
```

#### Переназначить ревьюера
```http
POST /avito-test-task/pullRequest/reassign
token: admin
Content-Type: application/json

{
  "pr_id": "pr-123",
  "old_reviewer_id": "user2",
  "new_reviewer_id": "user3"
}
```

### Metrics

```http
GET /metrics
```

## 📊 Мониторинг

### Prometheus

Метрики доступны на http://localhost:9090

### Grafana

Дашборды доступны на http://localhost:3000

**Учетные данные по умолчанию:**
- Username: `admin`
- Password: `admin`

## 🛠 Разработка

### Полезные команды Make

```bash
# Показать все доступные команды
make help

# Запустить сервисы
make up

# Остановить сервисы
make down

# Показать логи
make logs
make logs-app      # только приложение
make logs-db        # только база данных
make logs-vault     # только Vault


# Проверить статус сервисов
make status
make health

# Очистка
make clean          # временные файлы
make clean-all      # все включая Docker
```

### Структура проекта

```
.
├── cmd/
│   └── server/          # Точка входа приложения
├── internal/
│   ├── app/             # Инициализация приложения
│   ├── config/          # Конфигурация
│   ├── entity/          # Доменные сущности
│   ├── handler/         # HTTP обработчики
│   ├── middleware/      # Middleware (Prometheus, Auth)
│   ├── repository/      # Слой доступа к данным
│   └── usecase/         # Бизнес-логика
├── migrations/          # SQL миграции
├── infra/               # Docker и инфраструктура
│   ├── docker-compose.yml
│   ├── app.Dockerfile
│   └── prometheus/
└── Makefile
```

### Локальная разработка

1. **Запустите зависимости:**

```bash
make up
```

2. **Установите зависимости Go:**

```bash
go mod download
```

3. **Запустите приложение локально:**

```bash
go run cmd/server/main.go
```

### Миграции базы данных

Миграции находятся в директории `migrations/` и применяются автоматически при первом запуске PostgreSQL контейнера.

## 🔐 Безопасность

- Секреты хранятся в HashiCorp Vault
- Ролевая модель доступа (Admin/User)
- Middleware для проверки прав доступа

---

## Бенчмарки API

| Эндпоинт                                    | Метод | Описание           | Средняя задержка | RPS      | Передача данных |
| ------------------------------------------- | ----- | ------------------ | ---------------- | -------- | --------------- |
| /avito-test-task/team/add                   | POST  | Добавление команды | 34.87ms          | 10658.05 | 1.96MB/s        |
| /avito-test-task/team/get?team_name=backend | GET   | Получение команды  | 35.96ms          | 10373.21 | 4.34MB/s        |

### Подробные результаты

#### POST /avito-test-task/team/add
```wrk -t4 -c40 -d30s --latency -s post.lua http://localhost:8080/avito-test-task/team/add

Running 30s test @ http://localhost:8080/avito-test-task/team/add
  4 threads and 40 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    34.87ms   60.13ms 602.73ms   82.80%
    Req/Sec     2.69k     1.06k    6.05k    68.70%
  Latency Distribution
    50%    2.36ms
    75%   44.38ms
    90%  141.94ms
    99%  201.26ms
  320087 requests in 30.03s, 58.92MB read
  Non-2xx or 3xx responses: 320086
Requests/sec:  10658.05
Transfer/sec:      1.96MB
```

#### GET /avito-test-task/team/get?team_name=backend

```wrk -t4 -c40 -d30s --latency -s post.lua 'http://localhost:8080/avito-test-task/team/get?team_name=backend'

Running 30s test @ http://localhost:8080/avito-test-task/team/get?team_name=backend
  4 threads and 40 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    35.96ms   65.20ms 610.78ms   83.82%
    Req/Sec     2.61k     1.20k    6.03k    64.52%
  Latency Distribution
    50%    2.36ms
    75%   43.22ms
    90%  143.29ms
    99%  206.48ms
  311741 requests in 30.05s, 130.51MB read
Requests/sec:  10373.21
Transfer/sec:      4.34MB
```

### Примеры скриптов для wrk

#### post.lua (POST)

```wrk.method = "POST"
wrk.body   = '{"team_name":"payments","team_members":[{"user_id":"u1","username":"Alice","is_active":true},{"user_id":"u2","username":"Bob","is_active":true}]}'
wrk.headers["Content-Type"] = "application/json"
wrk.headers["token"] = "ADMIN"
```

#### post.lua (GET)

```wrk.method = "GET"
wrk.headers["token"] = "ADMIN"
```

---

| Эндпоинт                                       | Метод | Описание                       | Средняя задержка | RPS     | Передача данных |
| ---------------------------------------------- | ----- | ------------------------------ | ---------------- | ------- | --------------- |
| /avito-test-task/users/setIsActive             | POST  | Изменение статуса пользователя | 35.85ms          | 1199.14 | 337.13KB/s      |
| /avito-test-task/users/getReview?user_id=user2 | GET   | Получение ревью пользователя   | 41.81ms          | 8855.41 | 2.64MB/s        |

### Подробные результаты

#### POST /avito-test-task/users/setIsActive

```wrk -t4 -c40 -d30s --latency -s post.lua http://localhost:8080/avito-test-task/users/setIsActive

Running 30s test @ http://localhost:8080/avito-test-task/users/setIsActive
  4 threads and 40 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    35.85ms   23.29ms 338.33ms   92.31%
    Req/Sec   301.32     94.66   450.00     70.92%
  Latency Distribution
    50%   28.88ms
    75%   36.49ms
    90%   53.10ms
    99%  148.35ms
  36097 requests in 30.10s, 9.91MB read
Requests/sec:  1199.14
Transfer/sec:    337.13KB
```

#### GET /avito-test-task/users/getReview?user_id=user2

```wrk -t4 -c40 -d30s --latency -s post.lua 'http://localhost:8080/avito-test-task/users/getReview?user_id=user2'

Running 30s test @ http://localhost:8080/avito-test-task/users/getReview?user_id=user2
  4 threads and 40 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    41.81ms   79.00ms 789.61ms   85.22%
    Req/Sec     2.29k     1.06k    5.80k    67.67%
  Latency Distribution
    50%    2.75ms
    75%   52.00ms
    90%  153.44ms
    99%  379.11ms
  266558 requests in 30.10s, 79.57MB read
Requests/sec:  8855.41
Transfer/sec:      2.64MB
```

### Примеры скриптов для wrk

#### post.lua (POST)

```wrk.method = "POST"
wrk.body   = '{"user_id":"user4","is_active":false}'
wrk.headers["Content-Type"] = "application/json"
wrk.headers["token"] = "ADMIN"
```

#### post.lua (GET)

```wrk.method = "GET"
wrk.headers["token"] = "USER"
```

---

| Эндпоинт                              | Метод | Описание                    | Средняя задержка | RPS     | Передача данных |
| ------------------------------------- | ----- | --------------------------- | ---------------- | ------- | --------------- |
| /avito-test-task/pullRequest/create   | POST  | Создание пул-реквеста       | 38.20ms          | 8658.88 | 1.72MB/s        |
| /avito-test-task/pullRequest/merge    | POST  | Мерж пул-реквеста           | 33.17ms          | 1339.87 | 432.96KB/s      |
| /avito-test-task/pullRequest/reassign | POST  | Переназначение пул-реквеста | 28.68ms          | 4039.86 | 820.60KB/s      |

### Подробные результаты

#### POST /avito-test-task/pullRequest/create

```wrk -t4 -c40 -d30s --latency -s post.lua http://localhost:8080/avito-test-task/pullRequest/create

Running 30s test @ http://localhost:8080/avito-test-task/pullRequest/create
  4 threads and 40 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    38.20ms   70.04ms 705.67ms   84.22%
    Req/Sec     2.19k     1.02k    5.94k    66.81%
  Latency Distribution
    50%    2.82ms
    75%   47.13ms
    90%  146.55ms
    99%  233.34ms
  260404 requests in 30.07s, 51.65MB read
  Non-2xx or 3xx responses: 260404
Requests/sec:  8658.88
Transfer/sec:      1.72MB
```

#### POST /avito-test-task/pullRequest/merge

```wrk -t4 -c40 -d30s --latency -s post.lua http://localhost:8080/avito-test-task/pullRequest/merge

Running 30s test @ http://localhost:8080/avito-test-task/pullRequest/merge
  4 threads and 40 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    33.17ms   26.78ms 366.89ms   94.41%
    Req/Sec   338.89     92.64   515.00     81.70%
  Latency Distribution
    50%   26.25ms
    75%   32.20ms
    90%   45.21ms
    99%  170.09ms
  40323 requests in 30.09s, 12.72MB read
Requests/sec:  1339.87
Transfer/sec:    432.96KB
```

#### POST /avito-test-task/pullRequest/reassign

```wrk -t4 -c40 -d30s --latency -s post.lua http://localhost:8080/avito-test-task/pullRequest/reassign

Running 30s test @ http://localhost:8080/avito-test-task/pullRequest/reassign
  4 threads and 40 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    28.68ms   51.82ms 648.78ms   87.20%
    Req/Sec     1.02k   316.30     2.31k    76.73%
  Latency Distribution
    50%    8.18ms
    75%   12.73ms
    90%  109.84ms
    99%  206.66ms
  121428 requests in 30.06s, 24.09MB read
  Non-2xx or 3xx responses: 121428
Requests/sec:  4039.86
Transfer/sec:    820.60KB
```

### Примеры скриптов для wrk

#### post.lua (create)

```wrk.method = "POST"
wrk.body   = '{"pull_request_id":"pr-1001","pull_request_name":"Add search","author_id":"user1"}'
wrk.headers["Content-Type"] = "application/json"
wrk.headers["token"] = "ADMIN"
```

#### post.lua (merge)

```wrk.method = "POST"
wrk.body   = '{"pull_request_id":"pr-1001"}'
wrk.headers["Content-Type"] = "application/json"
wrk.headers["token"] = "ADMIN"
```

#### post.lua (reassign)

```wrk.method = "POST"
wrk.body   = '{"pull_request_id": "pr-1004","old_user_id": "user2"}'
wrk.headers["Content-Type"] = "application/json"
wrk.headers["token"] = "ADMIN"
```
---