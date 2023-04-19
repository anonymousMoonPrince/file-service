# File-service

Сервис распределенного хранения файлов

- **PUT** `http://localhost:8080/` - загрузка файла через multipart form
- **GET** `http://localhost:8080/<file_id>` - выгрузка файла по идентификатору

## Get Started

Команды:

- `make up` - поднять локальное окружение: minio, postgres
- `make migrate-up` - накатить миграции
- `make run` - запустить сервис
- `make run-checker` - запускает скрипт проверки загрузки и отгрузки файлов
