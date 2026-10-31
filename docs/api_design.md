# GoUptime — Проєктування REST API за специфікацією (API-First Design)

Документ визначає архітектуру REST API мікросервісу `GoUptime` на основі затвердженої специфікації вимог `spec/srs.md` (включно з Розділом 7 «Адендум: Програмний інтерфейс REST API»). Проєктування виконано за підходом API-First до початку написання серверної реалізації.

---

## 1. Вимоги SRS, що експонуються через REST API

Програмний інтерфейс надає доступ до управління життєвим циклом ресурсів у пам'яті (`In-Memory`) та спостереження за результатами HTTP/HTTPS-перевірок без залучення зовнішніх СУБД чи збереження стану на диск (`REQ-NF-001`, `REQ-NF-005`).

| ID вимоги | Категорія | Опис вимоги у контексті REST API |
| :--- | :--- | :--- |
| `REQ-F-002` | Функціональна | Кожен ресурс має обов'язкові непорожні поля `id` та `url`. |
| `REQ-F-003` | Функціональна | Унікальність ідентифікатора `id` у межах in-memory реєстру ресурсів (відхилення дублікатів із `409 Conflict`). |
| `REQ-F-004` | Функціональна | Валідація синтаксису URL та припустимих схем (`http` або `https`). |
| `REQ-F-005` | Функціональна | Автоматичне призначення дефолтних значень для пропущених числових полів: `interval_seconds=10`, `timeout_ms=2000`, `warning_threshold_ms=800`. |
| `REQ-F-006` | Функціональна | Відхилення `null`, значень `<= 0` (`400 Bad Request`) та порушення інваріанта `warning_threshold_ms >= timeout_ms` (`422 Unprocessable Entity`). |
| `REQ-F-007` | Функціональна | Виконання одного HTTP/HTTPS GET-запиту на перевірку із заданим `timeout_ms` без аналізу тіла відповіді. |
| `REQ-F-008` | Функціональна | Автоматичне слідування максимум за 3 редиректами; перевищення ліміту або цикл дає `DOWN` та `http_code=0`. |
| `REQ-F-010` | Функціональна | Класифікація стану `UP` при кінцевому коді `[200, 400)` та `latency_ms <= warning_threshold_ms`. |
| `REQ-F-011` | Функціональна | Класифікація стану `DEGRADED` при кінцевому коді `[200, 400)` та `warning_threshold_ms < latency_ms <= timeout_ms`. |
| `REQ-F-012` | Функціональна | Класифікація стану `DOWN` при мережевій помилці, тайм-ауті, `latency_ms > timeout_ms` або коді `>= 400`. |
| `REQ-F-013` | Функціональна | Структура результату перевірки `CheckResult` містить рівно 7 обов'язкових полів: `timestamp`, `id`, `url`, `status`, `http_code`, `latency_ms`, `error`. |
| `REQ-F-014` | Функціональна | Формат часової мітки `timestamp` — стандарт RFC3339 у часовій зоні UTC (`Z`). |
| `REQ-NF-006` | Нефункціональна | Fail-Fast валідація вхідного запиту до реєстрації ресурсу або ініціації перевірки. |
| `REQ-NF-007` | Нефункціональна | Стабільна множина машинно-читаних статусів: `{"UP", "DEGRADED", "DOWN"}`. |
| `REQ-API-001` | Адендум API | Читання списку зареєстрованих ресурсів (`GET /api/v1/resources`) та окремого ресурсу за `id` (`GET /api/v1/resources/{id}`). |
| `REQ-API-002` | Адендум API | Реєстрація нового ресурсу в пам'яті (`POST /api/v1/resources`) із дефолтами та Fail-Fast валідацією. |
| `REQ-API-003` | Адендум API | Видалення ресурсу з оперативної пам'яті (`DELETE /api/v1/resources/{id}`). |
| `REQ-API-004` | Адендум API | Примусовий запуск синхронної перевірки (`POST /api/v1/resources/{id}/checks`) та отримання останнього результату (`GET /api/v1/resources/{id}/results`). |
| `REQ-API-005` | Адендум API | Єдина модель помилок `ErrorResponse` без розкриття внутрішніх деталей реалізації Go. |

---

## 2. Декомпозиція предметної області на REST-ресурси

Модель ресурсів побудована у суворій відповідності до Розділу 3 «Глосарій предметної області» специфікації `spec/srs.md`. API абстрагує внутрішні механізми Go (горутини, канали, `sync.Mutex`, `context.Context`) і оперує виключно доменними сутностями.

### 2.1. Ресурс моніторингу (`Resource`)

Представляє цільовий HTTP/HTTPS endpoint, який перебуває під наглядом сервісу в оперативній пам'яті.

| Поле JSON | Тип даних | Обов'язковість (Create / Response) | Обмеження та дефолти | Термін у Глосарії SRS |
| :--- | :--- | :--- | :--- | :--- |
| `id` | `string` | Required / Required | Непорожній рядок (`minLength: 1`), унікальний, `nullable: false` | `ResourceID` |
| `url` | `string (uri)` | Required / Required | Схема `http://` або `https://`, непорожній хост, `nullable: false` | `Resource` (URL) |
| `interval_seconds` | `integer` | Optional / Required | Ціле число `>= 1`, дефолт `10`, `nullable: false` | `IntervalSeconds` |
| `timeout_ms` | `integer` | Optional / Required | Ціле число `>= 1`, дефолт `2000`, `nullable: false` | `TimeoutMs` |
| `warning_threshold_ms` | `integer` | Optional / Required | Ціле число `>= 1`, дефолт `800`, `< timeout_ms`, `nullable: false` | `WarningThresholdMs` |

### 2.2. Результат перевірки (`CheckResult`)

Представляє незмінний (immutable) запис результату однієї завершеної HTTP GET-перевірки ресурсу.

| Поле JSON | Тип даних | Обов'язковість | Обмеження та семантика | Термін у Глосарії SRS |
| :--- | :--- | :--- | :--- | :--- |
| `timestamp` | `string (date-time)` | Required | Формат RFC3339 у UTC (суфікс `Z`) | `CheckResult.Timestamp` |
| `id` | `string` | Required | Ідентифікатор перевіреного ресурсу | `ResourceID` |
| `url` | `string (uri)` | Required | Цільовий URL ресурсу | `Resource` (URL) |
| `status` | `string (enum)` | Required | Значення з множини `{"UP", "DEGRADED", "DOWN"}` | `AvailabilityStatus` |
| `http_code` | `integer` | Required | Кінцевий HTTP-код (`FinalResponse`) або `0` при мережевій помилці / тайм-ауті / циклі редиректів | `HTTPCode` |
| `latency_ms` | `integer` | Required | Виміряна затримка перевірки у мілісекундах (`>= 0`) | `LatencyMs` |
| `error` | `string` | Required | Текст помилки (`NetworkError` / `RequestTimeout`) або `""` при успішній перевірці | `CheckResult.Error` |

---

## 3. Специфікація URI, HTTP-методів та операцій

Базовий префікс API: `/api/v1`. Формат обміну даними: `application/json`.

| Operation ID | HTTP-метод та URI | Параметри | Request Body | Response Schema (Успіх) | Коди стану HTTP |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `listResources` | `GET /api/v1/resources` | — | — | `Resource[]` (масив ресурсів) | `200 OK` |
| `createResource` | `POST /api/v1/resources` | — | `ResourceCreateRequest` | `Resource` | `201 Created`<br>`400 Bad Request`<br>`409 Conflict`<br>`422 Unprocessable Entity` |
| `getResourceById` | `GET /api/v1/resources/{id}` | Path: `id` (`string`, required) | — | `Resource` | `200 OK`<br>`404 Not Found` |
| `deleteResource` | `DELETE /api/v1/resources/{id}` | Path: `id` (`string`, required) | — | Порожнє тіло | `204 No Content`<br>`404 Not Found` |
| `runResourceCheck` | `POST /api/v1/resources/{id}/checks` | Path: `id` (`string`, required) | — | `CheckResult` | `200 OK`<br>`404 Not Found` |
| `getResourceLastResult` | `GET /api/v1/resources/{id}/results` | Path: `id` (`string`, required) | — | `CheckResult` | `200 OK`<br>`404 Not Found` |

*Примітка щодо семантики `POST /api/v1/resources/{id}/checks`:* якщо цільовий зовнішній ресурс недоступний (тайм-аут або код 500), сам REST API `GoUptime` успішно виконує операцію перевірки та повертає HTTP-код `200 OK` із об'єктом `CheckResult`, де поле `"status": "DOWN"`. Код `404 Not Found` повертається лише тоді, коли ресурс із заданим `{id}` не зареєстрований у `GoUptime`.

---

## 4. Схеми даних та Єдина модель помилок (`ErrorResponse`)

Усі клієнтські помилки (`400`, `404`, `409`, `422`) використовують уніфіковану схему `ErrorResponse` (`REQ-API-005`).

### 4.1. Структура `ErrorResponse`

```json
{
  "timestamp": "2026-10-31T19:00:00Z",
  "code": "CONFIG_ERROR",
  "message": "Invalid resource configuration parameters",
  "details": [
    {
      "field": "warning_threshold_ms",
      "issue": "must be strictly less than timeout_ms (REQ-F-006)"
    }
  ]
}
```

### 4.2. Мапінг доменних помилок на HTTP-коди

| HTTP-код | Значення `code` в `ErrorResponse` | Умова виникнення (згідно з SRS) | Пов'язана вимога |
| :--- | :--- | :--- | :--- |
| `400 Bad Request` | `VALIDATION_ERROR` | Відсутні або порожні `id`/`url`; передано `null`; числові параметри `<= 0`; некоректна схема URL (наприклад, `ftp://`). | `REQ-F-002`, `REQ-F-004`, `REQ-F-006` |
| `404 Not Found` | `RESOURCE_NOT_FOUND` | Звернення за `id`, який відсутній в оперативній пам'яті сервісу. | `REQ-API-001`, `REQ-API-003`, `REQ-API-004` |
| `409 Conflict` | `DUPLICATE_RESOURCE_ID` | Спроба створити ресурс із `id`, що вже зареєстрований у пам'яті. | `REQ-F-003`, `REQ-API-002` |
| `422 Unprocessable Entity` | `CONFIG_ERROR` | Типи даних валідні, але порушено доменний інваріант: `warning_threshold_ms >= timeout_ms`. | `REQ-F-006`, `REQ-API-002` |

---

## 5. Таблиця простежуваності: API Operations ↔ REQ-ID ↔ BDD

| Operation ID | Метод та URI | Пов'язані вимоги SRS | BDD-сценарій | Сценарії верифікації Mock API |
| :--- | :--- | :--- | :--- | :--- |
| `listResources` | `GET /api/v1/resources` | `REQ-API-001`, `REQ-F-002`, `REQ-NF-005` | — | `POS-002` |
| `createResource` | `POST /api/v1/resources` | `REQ-API-002`, `REQ-API-005`, `REQ-F-002..006`, `REQ-NF-006` | — | `POS-001`, `NEG-001`, `NEG-002`, `NEG-003`, `NEG-004`, `NEG-005` |
| `getResourceById` | `GET /api/v1/resources/{id}` | `REQ-API-001`, `REQ-API-005`, `REQ-F-002` | — | `POS-002`, `NEG-006` |
| `deleteResource` | `DELETE /api/v1/resources/{id}` | `REQ-API-003`, `REQ-API-005`, `REQ-NF-005` | — | `POS-006`, `NEG-006` |
| `runResourceCheck` | `POST /api/v1/resources/{id}/checks` | `REQ-API-004`, `REQ-F-007..014`, `REQ-NF-004`, `REQ-NF-007`, `REQ-NF-008` | `BDD-001`, `BDD-002`, `BDD-003` | `POS-003`, `POS-004`, `POS-005`, `NEG-006` |
| `getResourceLastResult` | `GET /api/v1/resources/{id}/results` | `REQ-API-004`, `REQ-F-010..014`, `REQ-NF-005`, `REQ-NF-007` | `BDD-001`, `BDD-002`, `BDD-003` | `POS-003`, `NEG-006` |

---

## 6. Позитивні сценарії верифікації контракту (Mock API)

Базова адреса запущеного Stoplight Prism Mock Server: `http://127.0.0.1:4010/api/v1`.

### POS-001 — Реєстрація ресурсу із застосуванням дефолтів (`REQ-F-005`, `REQ-API-002`)

**Передумова:** передаються лише обов'язкові поля `id` та `url`.

**Запит:**

```bash
curl -i -X POST http://127.0.0.1:4010/api/v1/resources \
  -H "Content-Type: application/json" \
  -H "Prefer: code=201,example=CreatedWithDefaults" \
  -d '{"id": "kpi-main", "url": "https://kpi.ua"}'
```

**Очікуваний результат:** код `201 Created`, у тілі відповіді застосовано дефолти `interval_seconds: 10`, `timeout_ms: 2000`, `warning_threshold_ms: 800`.

### POS-002 — Отримання переліку зареєстрованих ресурсів (`REQ-API-001`)

**Запит:**

```bash
curl -i -X GET http://127.0.0.1:4010/api/v1/resources \
  -H "Accept: application/json" \
  -H "Prefer: code=200,example=ResourcesList"
```

**Очікуваний результат:** код `200 OK`, JSON-масив об'єктів `Resource`.

### POS-003 — Перевірка ресурсу зі статусом UP (`BDD-001`, `REQ-F-010`, `REQ-F-013`)

**Умова (BDD-001):** `warning_threshold_ms=800`, кінцевий код `200`, `latency_ms=650` (`650 <= 800`).

**Запит:**

```bash
curl -i -X POST http://127.0.0.1:4010/api/v1/resources/kpi-main/checks \
  -H "Accept: application/json" \
  -H "Prefer: code=200,example=BDD-001-UP"
```

**Очікуваний результат:** код `200 OK`, `status: "UP"`, `http_code: 200`, `latency_ms: 650`, `error: ""`.

### POS-004 — Перевірка ресурсу зі статусом DEGRADED після редиректу (`BDD-002`, `REQ-F-008`, `REQ-F-011`)

**Умова (BDD-002):** `warning_threshold_ms=800`, `timeout_ms=2000`, редирект `301 -> 200` (не більше 3 переходів), загальна `latency_ms=1200` (`800 < 1200 <= 2000`).

**Запит:**

```bash
curl -i -X POST http://127.0.0.1:4010/api/v1/resources/kpi-main/checks \
  -H "Accept: application/json" \
  -H "Prefer: code=200,example=BDD-002-DEGRADED"
```

**Очікуваний результат:** код `200 OK`, `status: "DEGRADED"`, `http_code: 200`, `latency_ms: 1200`, `error: ""`.

### POS-005 — Перевірка ресурсу зі статусом DOWN через тайм-аут або перевищення ліміту редиректів (`BDD-003`, `REQ-F-008`, `REQ-F-012`)

**Умова (BDD-003):** тайм-аут `2000 мс` перевищено або виявлено `> 3` редиректів.

**Запит:**

```bash
curl -i -X POST http://127.0.0.1:4010/api/v1/resources/kpi-main/checks \
  -H "Accept: application/json" \
  -H "Prefer: code=200,example=BDD-003-DOWN-Timeout"
```

**Очікуваний результат:** код `200 OK`, `status: "DOWN"`, `http_code: 0`, `latency_ms: 2000`, `error: "request timeout exceeded 2000ms"`.

### POS-006 — Видалення ресурсу з моніторингу в пам'яті (`REQ-API-003`)

**Запит:**

```bash
curl -i -X DELETE http://127.0.0.1:4010/api/v1/resources/kpi-main \
  -H "Prefer: code=204"
```

**Очікуваний результат:** код `204 No Content`, порожнє тіло відповіді.

---

## 7. Негативні сценарії верифікації контракту (Mock API)

### NEG-001 — Відсутній або порожній `id` чи `url` (`REQ-F-002`, `REQ-API-002`)

**Запит:**

```bash
curl -i -X POST http://127.0.0.1:4010/api/v1/resources \
  -H "Content-Type: application/json" \
  -H "Prefer: code=400,example=EmptyRequiredField" \
  -d '{"id": "", "url": "https://kpi.ua"}'
```

**Очікуваний результат:** код `400 Bad Request`, модель `ErrorResponse` з `code: "VALIDATION_ERROR"`.

### NEG-002 — Дублювання ідентифікатора ресурсу (`REQ-F-003`, `REQ-API-002`)

**Запит:**

```bash
curl -i -X POST http://127.0.0.1:4010/api/v1/resources \
  -H "Content-Type: application/json" \
  -H "Prefer: code=409,example=DuplicateResourceID" \
  -d '{"id": "kpi-main", "url": "https://kpi.ua"}'
```

**Очікуваний результат:** код `409 Conflict`, модель `ErrorResponse` з `code: "DUPLICATE_RESOURCE_ID"`.

### NEG-003 — Непідтримувана схема URL (`ftp://`) (`REQ-F-004`, `REQ-API-002`)

**Запит:**

```bash
curl -i -X POST http://127.0.0.1:4010/api/v1/resources \
  -H "Content-Type: application/json" \
  -H "Prefer: code=400,example=InvalidURLScheme" \
  -d '{"id": "ftp-srv", "url": "ftp://files.example.com/pub"}'
```

**Очікуваний результат:** код `400 Bad Request`, модель `ErrorResponse` з `code: "VALIDATION_ERROR"`.

### NEG-004 — Передача `null` або числа `<= 0` (`REQ-F-006`, `REQ-API-002`)

**Запит:**

```bash
curl -i -X POST http://127.0.0.1:4010/api/v1/resources \
  -H "Content-Type: application/json" \
  -H "Prefer: code=400,example=NullOrNonPositiveValue" \
  -d '{"id": "bad-timeout", "url": "https://kpi.ua", "timeout_ms": 0}'
```

**Очікуваний результат:** код `400 Bad Request`, модель `ErrorResponse` з `code: "VALIDATION_ERROR"`.

### NEG-005 — Порушення граничної умови `warning_threshold_ms >= timeout_ms` (`REQ-F-006`, `RISK-001`)

**Запит:**

```bash
curl -i -X POST http://127.0.0.1:4010/api/v1/resources \
  -H "Content-Type: application/json" \
  -H "Prefer: code=422,example=ThresholdExceedsTimeout" \
  -d '{"id": "bad-threshold", "url": "https://kpi.ua", "timeout_ms": 2000, "warning_threshold_ms": 2000}'
```

**Очікуваний результат:** код `422 Unprocessable Entity`, модель `ErrorResponse` з `code: "CONFIG_ERROR"`.

### NEG-006 — Звернення до неіснуючого ресурсу (`REQ-API-001`, `REQ-API-004`, `REQ-API-005`)

**Запит:**

```bash
curl -i -X GET http://127.0.0.1:4010/api/v1/resources/unknown-id \
  -H "Accept: application/json" \
  -H "Prefer: code=404,example=ResourceNotFound"
```

**Очікуваний результат:** код `404 Not Found`, модель `ErrorResponse` з `code: "RESOURCE_NOT_FOUND"`.