# TODO API (Fiber + Postgres) тестовое задание

## Установка

### 1. Установите Go (если не установлен)  
[Скачать Go](https://go.dev/dl/)

### 2. Установите PostgreSQL (если не установлен)  
[Скачать PostgreSQL](https://www.postgresql.org/download/)

### 3. Создайте базу данных  

```sql
CREATE DATABASE todolist;
```
### 4. Скопируйте .env.example → configs/local.env и укажите данные подключения

### 5. Установите зависимости и запустите сервер
```sh
go mod tidy
go run cmd/main.go
```
API запустится на http://localhost:8083.

### 6. Документация Swagger

Открой http://localhost:8083/swagger/index.html.