.PHONY: help build up down restart logs clean test migrate-up migrate-down docker-build docker-up docker-down

# Переменные
COMPOSE_FILE = docker-compose.yml
APP_NAME = avito_test_task
DOCKER_COMPOSE = docker-compose -f $(COMPOSE_FILE)

# Цвета для вывода
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
RESET  := $(shell tput -Txterm sgr0)

##@ Основные команды

help: ## Показать эту справку
	@echo '$(GREEN)Доступные команды:$(RESET)'
	@awk 'BEGIN {FS = ":.*##"; printf "\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  $(YELLOW)%-15s$(RESET) %s\n", $$1, $$2 } /^##@/ { printf "\n$(GREEN)%s$(RESET)\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

init: ## Инициализировать Vault и секреты через скрипт
	@echo "$(GREEN)Запуск Vault...$(RESET)"
	@echo "$(GREEN)Запуск скрипта инициализации Vault...$(RESET)"
	bash init-vault.sh
	@echo "$(GREEN)Секреты инициализированы. Готово!$(RESET)"

build: ## Собрать Docker образы
	@echo "$(GREEN)Сборка Docker образов...$(RESET)"
	$(DOCKER_COMPOSE) build

up: ## Запустить все сервисы (после init!)
	@echo "$(GREEN)Запуск всех сервисов...$(RESET)"
	$(DOCKER_COMPOSE) up -d
	@echo "$(GREEN)Сервисы запущены!$(RESET)"
	@echo "$(YELLOW)Приложение: http://localhost:8080$(RESET)"
	@echo "$(YELLOW)Vault UI: http://localhost:8201$(RESET)"
	@echo "$(YELLOW)Prometheus: http://localhost:9090$(RESET)"
	@echo "$(YELLOW)Grafana: http://localhost:3000 (admin/admin)$(RESET)"

down: ## Остановить все сервисы
	@echo "$(GREEN)Остановка всех сервисов...$(RESET)"
	$(DOCKER_COMPOSE) down

restart: down up ## Перезапустить все сервисы

logs: ## Показать логи всех сервисов
	$(DOCKER_COMPOSE) logs -f

logs-app: ## Показать логи приложения
	$(DOCKER_COMPOSE) logs -f app

logs-db: ## Показать логи базы данных
	$(DOCKER_COMPOSE) logs -f postgres

logs-vault: ## Показать логи Vault
	$(DOCKER_COMPOSE) logs -f vault

##@ Разработка

dev: ## Запустить в режиме разработки (с логами)
	$(DOCKER_COMPOSE) up

##@ База данных

migrate-up: ## Применить миграции
	@echo "$(GREEN)Применение миграций...$(RESET)"
	@if [ -d "migrations" ]; then \
		for file in migrations/*.up.sql; do \
			if [ -f "$$file" ]; then \
				echo "Применение: $$file"; \
			fi \
		done \
	fi

migrate-down: ## Откатить миграции
	@echo "$(GREEN)Откат миграций...$(RESET)"
	@if [ -d "migrations" ]; then \
		for file in migrations/*.down.sql; do \
			if [ -f "$$file" ]; then \
				echo "Откат: $$file"; \
			fi \
		done \
	fi

db-shell: ## Подключиться к базе данных
	$(DOCKER_COMPOSE) exec postgres psql -U postgres -d avito-test-task-db

##@ Docker

docker-build: build ## Собрать Docker образы (алиас)

docker-up: up ## Запустить сервисы (алиас)

docker-down: down ## Остановить сервисы (алиас)

docker-clean: ## Удалить все контейнеры, образы и volumes
	@echo "$(YELLOW)Внимание! Это удалит все контейнеры, образы и volumes!$(RESET)"
	$(DOCKER_COMPOSE) down -v --rmi all
	docker system prune -f

docker-ps: ## Показать запущенные контейнеры
	$(DOCKER_COMPOSE) ps

docker-stats: ## Показать статистику контейнеров
	docker stats

##@ Очистка

clean: ## Очистить временные файлы
	@echo "$(GREEN)Очистка временных файлов...$(RESET)"
	go clean
	rm -f coverage.out coverage.html

clean-all: clean docker-clean ## Полная очистка (включая Docker)

##@ Утилиты

status: ## Показать статус всех сервисов
	@echo "$(GREEN)Статус сервисов:$(RESET)"
	$(DOCKER_COMPOSE) ps
	@echo "\n$(GREEN)Проверка здоровья:$(RESET)"
	@echo "Приложение: $$(curl -s -o /dev/null -w '%{http_code}' http://localhost:8080/health || echo 'недоступно')"
	@echo "Vault: $$(curl -s -o /dev/null -w '%{http_code}' http://localhost:8201/v1/sys/health || echo 'недоступно')"

##@ По умолчанию
.DEFAULT_GOAL := help
