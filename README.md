# otus-user-service
User API

### Директория с чартом сервиса

`infra/user-service`

### Установка

`helm upgrade --install -f infra/user-service/values.yaml user-service infra/user-service/.`

### Тестирование

`newman run integration_tests/user_api_collection.json`

### Удаление

`helm uninstall user-service`
