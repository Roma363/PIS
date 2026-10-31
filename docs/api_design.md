# GoUptime — API Design (Mock API)

## 1. Призначення

Документ визначає декомпозицію REST-ресурсів, URI, HTTP-методів, схем даних і
єдиної моделі помилок для Mock API, що покриває сутності `Resource` та
`CheckResult`.

## 2. Ресурсна декомпозиція

### 2.1. `Resource`

Описує налаштування endpoint, який підлягає моніторингу.

```json
{
  "id": "api-main",
  "url": "https://example.com/health",
  "interval_seconds": 10,
  "timeout_ms": 2000,
  "warning_threshold_ms": 800,
  "created_at": "2026-10-31T18:00:00Z",
  "updated_at": "2026-10-31T18:00:00Z"
}
```

### 2.2. `CheckResult`

Описує результат однієї виконаної перевірки ресурсу.

```json
{
  "id": "chk-0001",
  "resource_id": "api-main",
  "timestamp": "2026-10-31T18:01:00Z",
  "status": "UP",
  "http_code": 200,
  "latency_ms": 123,
  "error": ""
}
```

## 3. URI та методи

Базовий префікс: `/api/v1`.

| URI | Метод | Призначення | Успішний статус |
|---|---|---|---|
| `/resources` | `GET` | Отримати список ресурсів | `200 OK` |
| `/resources` | `POST` | Створити ресурс | `201 Created` |
| `/resources/{resource_id}` | `GET` | Отримати ресурс за ID | `200 OK` |
| `/resources/{resource_id}` | `PUT` | Повністю оновити ресурс | `200 OK` |
| `/resources/{resource_id}` | `PATCH` | Частково оновити ресурс | `200 OK` |
| `/resources/{resource_id}` | `DELETE` | Видалити ресурс | `204 No Content` |
| `/resources/{resource_id}/check` | `POST` | Запустити одиничну перевірку | `201 Created` |
| `/check-results` | `GET` | Отримати історію перевірок | `200 OK` |
| `/check-results/{check_result_id}` | `GET` | Отримати результат за ID | `200 OK` |

Параметри `GET /check-results`:

- `resource_id` — фільтр за ресурсом.
- `limit` — обмеження кількості записів (`1..100`).

## 4. Єдина модель помилки — `ErrorResponse`

Усі помилки повертаються в однаковому форматі:

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "warning_threshold_ms must be less than timeout_ms",
    "details": [
      {
        "field": "warning_threshold_ms",
        "issue": "must be < timeout_ms"
      }
    ],
    "request_id": "req-9f1b8d1d"
  }
}
```

Семантика полів:

- `code` — стабільний машинно-читаний код (`VALIDATION_ERROR`,
  `RESOURCE_NOT_FOUND`, `CONFLICT`, `INTERNAL_ERROR`).
- `message` — людино-читаний опис помилки.
- `details` — масив уточнень для валідаційних збоїв.
- `request_id` — ідентифікатор трасування запиту.

## 5. Сценарії верифікації Mock API

### 5.1. Позитивні сценарії

| ID | Назва | Передумови | Кроки | Очікуваний результат |
|---|---|---|---|---|
| `POS-001` | Створення `Resource` | Валідний payload | `POST /resources` | `201`, повернено створений `Resource` |
| `POS-002` | Отримання списку `Resource` | Існує >=1 ресурс | `GET /resources` | `200`, масив ресурсів |
| `POS-003` | `UP`-класифікація перевірки (BDD-001) | `warning_threshold_ms=800`, код `200`, `latency_ms=650` | `POST /resources/{id}/check` | `201`, `CheckResult.status=UP`, `http_code=200` |
| `POS-004` | `DEGRADED` після редиректу (BDD-002) | Кінцева відповідь `200` після редиректу, `latency_ms=1200` | `POST /resources/{id}/check` | `201`, `status=DEGRADED`, `error=""` |
| `POS-005` | `DOWN` через timeout/redirect failure (BDD-003) | Timeout або >3 редиректи/цикл | `POST /resources/{id}/check` | `201`, `status=DOWN`, `http_code=0`, непорожній `error` |
| `POS-006` | Отримання історії `CheckResult` | Є результати перевірок | `GET /check-results?resource_id={id}&limit=10` | `200`, відфільтрований список `CheckResult` |

### 5.2. Негативні сценарії

| ID | Назва | Кроки | Очікуваний результат |
|---|---|---|---|
| `NEG-001` | Порожній `id` під час створення ресурсу | `POST /resources` з `id=""` | `400`, `ErrorResponse.code=VALIDATION_ERROR` |
| `NEG-002` | Некоректний URL | `POST /resources` з `url="ftp://example.com"` | `400`, `VALIDATION_ERROR` |
| `NEG-003` | Некоректні пороги часу | `POST /resources` з `warning_threshold_ms >= timeout_ms` | `400`, `VALIDATION_ERROR` |
| `NEG-004` | Дубльований `resource_id` | Повторний `POST /resources` з існуючим `id` | `409`, `ErrorResponse.code=CONFLICT` |
| `NEG-005` | Ресурс не знайдено | `GET /resources/{unknown_id}` або `POST /resources/{unknown_id}/check` | `404`, `ErrorResponse.code=RESOURCE_NOT_FOUND` |
| `NEG-006` | Некоректні query-параметри історії | `GET /check-results?limit=0` | `400`, `ErrorResponse.code=VALIDATION_ERROR` |
