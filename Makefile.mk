.PHONY: help build run clean test docker-build docker-up docker-down docker-logs migrate lint swagger install

# متغیرها
APP_NAME=pra-exchange
DOCKER_IMAGE=$(APP_NAME):latest
DOCKER_CONTAINER=$(APP_NAME)

# رنگ‌ها برای خروجی
GREEN=\033[0;32m
RED=\033[0;31m
YELLOW=\033[0;33m
NC=\033[0m # No Color

help: ## نمایش راهنما
	@echo "${GREEN}Available commands:${NC}"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "${YELLOW}%-20s${NC} %s\n", $$1, $$2}'

install: ## نصب وابستگی‌ها
	@echo "${GREEN}Installing dependencies...${NC}"
	go mod download
	go mod tidy

build: ## ساخت باینری
	@echo "${GREEN}Building binary...${NC}"
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $(APP_NAME) .

run: ## اجرای برنامه به صورت محلی
	@echo "${GREEN}Running application...${NC}"
	go run main.go

test: ## اجرای تست‌ها
	@echo "${GREEN}Running tests...${NC}"
	go test -v -cover ./...

lint: ## اجرای linter
	@echo "${GREEN}Running linter...${NC}"
	golangci-lint run

clean: ## پاک کردن فایل‌های build
	@echo "${GREEN}Cleaning...${NC}"
	rm -f $(APP_NAME)
	rm -rf ./tmp

docker-build: ## ساخت Docker image
	@echo "${GREEN}Building Docker image...${NC}"
	docker build -t $(DOCKER_IMAGE) .

docker-up: ## راه‌اندازی همه سرویس‌ها با Docker Compose
	@echo "${GREEN}Starting all services...${NC}"
	docker-compose up -d
	@echo "${YELLOW}Services:${NC}"
	@echo "  API: http://localhost:8080"
	@echo "  PgAdmin: http://localhost:5050 (admin@pra-exchange.com / admin123)"
	@echo "  PostgreSQL: localhost:5432"
	@echo "  Redis: localhost:6379"

docker-down: ## توقف همه سرویس‌ها
	@echo "${GREEN}Stopping all services...${NC}"
	docker-compose down

docker-down-v: ## توقف و حذف والیوم‌ها (پاک کردن دیتابیس)
	@echo "${RED}Stopping and removing volumes...${NC}"
	docker-compose down -v

docker-logs: ## مشاهده لاگ‌ها
	docker-compose logs -f

docker-logs-api: ## مشاهده لاگ‌های API
	docker-compose logs -f pra-backend

docker-restart: ## ری‌استارت همه سرویس‌ها
	@echo "${GREEN}Restarting services...${NC}"
	docker-compose restart

migrate: ## اجرای مایگریشن دیتابیس (اگر نیاز باشد)
	@echo "${GREEN}Running migrations...${NC}"
	# go run scripts/migrate.go

swagger: ## تولید مستندات Swagger
	@echo "${GREEN}Generating Swagger docs...${NC}"
	# swag init -g api/server.go

status: ## نمایش وضعیت سرویس‌ها
	@echo "${GREEN}Service status:${NC}"
	docker-compose ps
	@echo "\n${GREEN}API Health:${NC}"
	curl -s http://localhost:8080/health | jq . || echo "API not responding"

logs-all: ## مشاهده لاگ‌های همه سرویس‌ها
	docker-compose logs

psql: ## اتصال به دیتابیس PostgreSQL
	docker exec -it pra_postgres psql -U pra_user -d pra_exchange

redis-cli: ## اتصال به Redis
	docker exec -it pra_redis redis-cli

dev: ## اجرای توسعه با بازسازی خودکار (نیاز به air)
	@echo "${GREEN}Starting dev mode with hot reload...${NC}"
	air

.PHONY: help build run clean test docker-build docker-up docker-down docker-logs migrate swagger status dev