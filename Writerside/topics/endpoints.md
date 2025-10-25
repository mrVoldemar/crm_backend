# API Endpoints

API Gateway предоставляет несколько типов endpoints для взаимодействия с внутренними сервисами и мониторинга состояния системы.

## Маршруты сервисов

Для каждого сервиса, определенного в `config/services.json`, автоматически создаются следующие маршруты:

### Список ресурсов

```
GET /api/{service}s
```

Получение списка ресурсов сервиса.

### Конкретный ресурс

```
GET /api/{service}s/{id}
PUT /api/{service}s/{id}
DELETE /api/{service}s/{id}
```

Операции с конкретным ресурсом сервиса по идентификатору.

### Создание ресурса

```
POST /api/{service}s
```

Создание нового ресурса в сервисе.

## Общие параметры

Все маршруты сервисов требуют:
- Аутентификации через заголовок `Authorization`
- Соблюдения лимита запросов (100 в минуту на клиента)

## Health Check

```
GET /health
```

Endpoint для проверки состояния API Gateway. Не требует аутентификации и ограничения по rate limit.

Ответ:
- Код: 200 OK
- Тело: "API Gateway is running"

## Обработка ошибок

API Gateway стандартизирует обработку ошибок для всех сервисов. При возникновении ошибок возвращается JSON следующего формата:

```json
{
  "error": "Описание ошибки",
  "code": 400
}
```

### Возможные коды ошибок

- `400 Bad Request` - некорректный запрос
- `401 Unauthorized` - отсутствует или недействительный токен авторизации
- `404 Not Found` - запрашиваемый ресурс не найден
- `429 Too Many Requests` - превышен лимит запросов
- `500 Internal Server Error` - внутренняя ошибка сервера
- `502 Bad Gateway` - ошибка при обращении к внутреннему сервису

## Примеры запросов

### Получение списка сотрудников

```bash
curl -H "Authorization: Bearer your_token" \
     http://api-gateway:8080/api/employees
```

### Создание нового клиента

```bash
curl -X POST \
     -H "Authorization: Bearer your_token" \
     -H "Content-Type: application/json" \
     -d '{"name": "John Doe", "email": "john@example.com"}' \
     http://api-gateway:8080/api/clients
```

### Получение информации о конкретном сотруднике

```bash
curl -H "Authorization: Bearer your_token" \
     http://api-gateway:8080/api/employees/123
```