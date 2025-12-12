# ReviewLink — Backend

ReviewLink — backend-сервис для сбора отзывов через QR-коды.  
Пользователь сканирует QR-код, подтверждает визит и оставляет отзыв, а за оценки **4★ и 5★** получает бонусы.  
Проект включает пользовательскую и административную части.

---

##  Структура проекта

### Backend

```text
backend/
├── cmd/
│   └── main.go                  # Точка входа в приложение
│
├── configs/
│   └── config.go                # Загрузка конфигурации
│
├── docs/                        # Swagger документация
│   ├── swagger.json
│   └── swagger.yaml
│
├── internal/                    # Основная логика сервиса
│   ├── app/                     # Инициализация приложения
│   ├── controller/              # HTTP handlers (Gin)
│   ├── model/                   # Модели и JWT claims
│   ├── repository/              # Работа с PostgreSQL
│   ├── service/                 # Бизнес-логика
│   └── tests/                   # Интеграционные тесты
│
├── migrations/                  # SQL-миграции
│
├── pkg/
│   └── middleware/
│       └── jwt.go               # JWT middleware
│
├── go.mod
├── go.sum
└── .env.example                 # Пример окружения 
```
### Frontend (демо-UI)
```text
frontend/
├── admin.html
├── admin.js
│
├── dashboard.html
├── dashboard.js
│
├── login.html
├── login.js
│
├── review-form.html
├── review-form.js
│
└── style.css
```

Swagger
API-документация доступна после запуска сервера:
```text
/swagger/index.html
```