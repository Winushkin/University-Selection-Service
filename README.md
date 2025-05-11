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
- Golang
- gRPC + grpc-gateway
- PostgreSQL
- Docker + Docker Compose
- Zap logger
- pgx / SQL миграции
- Nginx balancer
---

## 🚀 Быстрый старт

### Требования

- Установленный Docker
- Установленный Docker Compose

### Шаги запуска

1. Клонируй репозиторий

```bash
$ git clone https://github.com/yourname/University-Selection-Service.git
$ cd University-Selection-Service
```

2. Настрой переменные окружения

```bash
$ cp .env.example .env
```

3. Запусти проект
```bash
$ make build
```

## Структура проекта
````
University-Selection-Service/
├──.github
│  └──workflows
│     └──ci.yaml
│
├──backend
│  ├── api
│  │   ├── analytic.proto
│  │   └── user.proto
│  ├── cmd
│  │    ├── analytic
│  │    │   ├── Dockerfile
│  │    │   └── main.go
│  │    ├── gateway
│  │    │   ├── Dockerfile
│  │    │   └── main.go
│  │    ├── nginx
│  │    │   ├── Dockerfile
│  │    │   └── nginx.conf
│  │    ├── university
│  │    │   └── main.go
│  │    └── user
│  │        ├── Dockerfile
│  │        └── main.go
│  ├── db
│  │   └── migrations
│  │       ├── universities
│  │       │   ├── 000001_create_universities_schema_down.sql
│  │       │   ├── 000001_create_universities_schema_up.sql
│  │       │   ├── 000002_create_universities_schema_up.sql
│  │       │   ├── 000002_create_regions_table_down.sql
│  │       │   ├── 000002_create_regions_table_up.sql
│  │       │   ├── 000003_create_universities_table_down.sql
│  │       │   ├── 000003_create_universities_table_up.sql
│  │       │   ├── 000004_create_specialties_table_down.sql
│  │       │   └── 000004_create_specialties_table_up.sql
│  │       └── users
│  │           ├── 000001_initialize_users_schema_down.sql
│  │           ├── 000001_initialize_users_schema_up.sql
│  │           ├── 000002_add_table_users_down.sql
│  │           ├── 000002_add_table_users_up.sql
│  │           ├── 000003_add_table_refresh_tokens_down.sql
│  │           └── 000003_add_table_refresh_tokens_up.sql
│  │
│  │
│  │
│
│
│
│
│
│
│
│
│
│
````

## Описание работы приложения
Приложение **University Selection Service** состоит из двух масштабных частей - _клиентской_ и _серверной_, которые взаимодействуют друг с другом с помощью _HTTP_-запросов.

## Описание клиентской части

## Описание серверной части
**Серверная часть** приложения состоит из нескольких микросервисов на языке _Go_ и баз данных _PostgreSQL_. Серверная часть обрабатывает HTTP-запросы от клиентской части, и распределяет и по соответсвующим микросервисам, которые, взаимодействуя между собой и с базами данных, собирают ответ и передают его в руки клиентской части.

- Выбор **микросервисной** модели для серверной части приложения обусловлен важностью преимуществ данной модели в рамках разработки приложения. Изоляция отдельных компонентов позволила вносить изменения безболезенно для остальных комнонетов, а также обеспечила гибкость и свободу выбора технологий. Также посредством данной модели была обеспечена масштабируемость серверной части.
- Выбор языка _Go_ обусловлен его высокой производительностью и скоростью разработки на нем, также поддержка gRPC в языке обеспечила высокую скорость обмена данных между сервисами.
- Выбор СУБД _PostgreSQL_ обусловлен ее открытостью, надежностью и кроссплатформенностью, а также поддержкой сложных запросов и больших объемов данных.

Серверная часть состоит из: 

Микросервисов:
1. User - производит аутентификацию и авторизацию пользователей в системе, а также отвечает за управление профилями пользователей
2. Analytic - анализирует данные пользователя и на их основе подбирает наиболее подходящие университеты с помощью "Метода анализа иерархий" из собранных в базе данных.
3. Universities - собирает данные о ВУЗах в открытых интернет-источниках и собирает их в базу данных.

Баз данных:
1. База данных университетов
2. База данных пользователей

### Сервис User
Сервис **User** содержит:

- _gRPC-сервер_ для взаимодействия с клиентом и другими сервисами  - `backend/internal/user/server.go`
- _Репозиторий_ для взаимодействия с базой данных пользователей - `backend/internal/repositories/user_repository.go`
- _Конфигурацию_ для настройки сервиса - `backend/internal/config/user_config.go`

Сервис **User** выставляет наружу для взаимодействия с клиентом следующие endpoint-ы:
1. /api/user/signup - производится регистрация пользователя в системе и создание профиля в базе данных пользователей.
2. /api/user/login - производится авторизация пользователя в системе и выдача прав на взаимодействие с остальным функционалом системы.
3. /api/user/refresh - обновление access-токена для авторизации запросов пользователя в системе.
4. /api/user/fill - обновление данных профиля в базе данных.
5. /api/user/logout - выход пользователя из системы и отзыв прав на взаимодействие с системой.

Для взаимодействия с сервисом **Analytic** выделен endpoint:
1. /api.UserService/ProfileDataForAnalytic - предоставляет данные профиля сервису **Analytic**.

### Сервис Analytic
Сервис **Analytic** содержит:

- _gRPC-сервер_ для взаимодействия с клиентом и другими сервисами - `backend/internal/analytic/server/server.go`
- _Репозиторий_ для взаимодействия с базой данных университетов - `backend/internal/repositories/analytic_repository.go`
- _Конфигурацию_ для настройки сервиса - backend/internal/config/analytic_config.go
- _Анализатор_ для выполнения `Метода анализа иерархий` с заданными параметрами - `backend/internal/analytic/analyze/analyze.go`

Сервис **Analytic** выставляет наружу для взаимодействия с клиентом следующий endpoint:
1. /api/analytic/analyze - выполняется `Метод анализа иерархий` с параметрами из запроса и данными профиля пользователя.

### Сервис University

