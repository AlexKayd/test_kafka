### Проект демонстрирует взаимодействие микросервисов с использованием Kafka на языке Go. Пример основан на системе авторизации пользователя с возможностью передачи событий (например, логин, регистрация) в другие сервисы через Kafka.

- **auth_service** — HTTP API (Go), реализует регистрацию и логин. Отправляет события в Kafka.
- **track_service** — потребитель Kafka, получает события пользователей и логирует их.
- **Kafka + Zookeeper** — запускаются через Docker Compose.

---

### Запуск

```bash
docker-compose up -d

go run track_service.go events.go
go run auth_service.go events.go

```
Пример запроса: curl "http://localhost:8082/login?username=alex&password=123"