# 🎓 University Selection Service

**University Selection Service** — микросервисное веб-приложение для подбора университетов по критериям будущих абитуриентов

---

## 🧱 Архитектура
	•	gateway — HTTP вход, преобразование HTTP → gRPC (grpc-gateway)
	•	user_service — регистрация, логин, JWT
	•	analytic_service — хранение аналитических событий
	•	integration_service — загрузка данных о вузах с внешних источников (например, ege.hse.ru)
	•	PostgreSQL — централизованное хранилище данных, разбитое по схемам

## 📦 Стек технологий
	•	Golang
	•	gRPC + grpc-gateway
	•	PostgreSQL
	•	Docker + Docker Compose
	•	Zap logger
	•	pgx / SQL миграции
---

## 🚀 Быстрый старт

### 1. Клонируй репозиторий

```bash
$ git clone https://github.com/yourname/University-Selection-Service.git
$ cd University-Selection-Service
```

### 2. Настрой переменные окружения

```bash
$ cp .env.example .env
```

### 3. Запусти проект
```bash
$ make build
```






