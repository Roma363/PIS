# Матриця простежуваності вимог

Матриця розділяє чотири різні твердження:

- **Контракт** — вимога представлена REST-операцією або схемою OpenAPI.
- **Mock-сценарій** — Prism повернув статичний приклад відповіді; це не доказ реалізації бізнес-логіки.
- **Лабораторна №4** — фактично перевірена форма контракту та доступність mock-операції.
- **Лабораторна №5** — поведінка, яку можна підтвердити лише реальною Go-реалізацією та тестами.

`—` означає, що вимога не належить до REST-контракту або не має обґрунтованого mock-сценарію в лабораторній №4.

| ID вимоги | Опис за SRS | Критерій приймання за SRS | ID задачі | ID тесту | Код/модуль наступного етапу | API operation (контракт) | Mock-сценарій (лише статичний приклад) | Статус лабораторної №4 |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| **REQ-F-001** | Завантаження конфігурації виключно з локального `config.json`. | Відсутній/недоступний файл завершує процес до HTTP-запитів. | `TASK-004` | `TASK-003` | `internal/config` | — | — | **Lab5** — CLI-конфігурація не є REST-контрактом |
| **REQ-F-002** | Конфігурація містить `resources`; кожен ресурс має непорожні `id` і `url`. | Порожні/відсутні поля дають `ConfigError` до перевірок. | `TASK-006` | `TASK-005` | `internal/config`, `ManagementAPI` | `POST /api/v1/resources` (request schema) | `POS-001`, `NEG-001` | **Contract + blocked mock** — Prism не підтверджує реальну ConfigError |
| **REQ-F-003** | Ідентифікатори ресурсів унікальні. | Дубльований `id` дає `ConfigError`. | `TASK-006` | `TASK-005` | `internal/config`, `ManagementAPI` | `POST /api/v1/resources` (API-002) | `NEG-002` | **Contract + static mock** — унікальність потребує Lab5 |
| **REQ-F-004** | URL ресурсу є коректним HTTP- або HTTPS-URL. | Інша схема, порожній хост або помилковий синтаксис дають `ConfigError`. | `TASK-006` | `TASK-005` | `internal/config`, `ManagementAPI` | `POST /api/v1/resources` (request schema) | `NEG-003` | **Contract + blocked mock** — Prism не підтверджує документований `400` |
| **REQ-F-005** | Пропущені числові поля отримують дефолти `10/2000/800`. | Ресурс без числових полів завантажується з дефолтами. | `TASK-006` | `TASK-005` | `internal/config`, `ManagementAPI` | `POST /api/v1/resources` | `POS-001` | **Contract + static mock** — фактичне застосування дефолтів потребує Lab5 |
| **REQ-F-006** | Відхиляються `null`, значення `<= 0` і `warning_threshold_ms >= timeout_ms`. | Кожна умова дає `ConfigError` до першого HTTP-запиту. | `TASK-006` | `TASK-005` | `internal/config`, `ManagementAPI` | `POST /api/v1/resources` (schema + API-002) | `NEG-004`, `NEG-005` | **Contract + blocked/partial mock** — інваріант і точний `400` потребують Lab5 |
| **REQ-F-007** | Для кожного ресурсу виконується HTTP/HTTPS GET із заданим timeout. | Не більше одного GET за перевірку; response body не аналізується. | `TASK-011` | `TASK-010` | `internal/probe`, `ManagementAPI` | `POST /api/v1/resources/{id}/checks` | `POS-003`, `POS-004`, `POS-005` | **Contract example only** — реальний GET потребує Lab5 |
| **REQ-F-008** | Клієнт слідує максимум за трьома редиректами. | Цикл/перевищення ліміту дає `DOWN`, `HTTPCode=0`. | `TASK-013` | `TASK-012` | `internal/probe`, `ManagementAPI` | `POST /api/v1/resources/{id}/checks` | `POS-004`, `POS-005` | **Contract example only** — Prism не виконує редиректи |
| **REQ-F-009** | Перевірки повторюються з інтервалом кожного ресурсу незалежно. | Попередній результат не змінює класифікацію наступної перевірки. | `TASK-017` | `TASK-016` | `internal/scheduler` | — | — | **Lab5** — DELETE не доводить незалежність scheduler-інтервалів |
| **REQ-F-010** | Статус `UP` для коду `[200,400)` і latency `<= WarningThresholdMs`. | Відповідь у межах порога отримує `UP`. | `TASK-008` | `TASK-007` | `internal/domain`, `ManagementAPI` | `POST /api/v1/resources/{id}/checks`, `GET /api/v1/resources/{id}/results` | `POS-003` | **Contract example only** — класифікація потребує Lab5 |
| **REQ-F-011** | Статус `DEGRADED` для коду `[200,400)` і latency між порогом та timeout. | Відповідь у зазначеному інтервалі отримує `DEGRADED`. | `TASK-008` | `TASK-007` | `internal/domain`, `ManagementAPI` | `POST /api/v1/resources/{id}/checks` | `POS-004` | **Contract example only** — класифікація потребує Lab5 |
| **REQ-F-012** | Статус `DOWN` для мережевої помилки, timeout, latency понад timeout або коду `>=400`. | Кожен випадок дає `DOWN`; без HTTP-відповіді `HTTPCode=0`. | `TASK-008`, `TASK-011` | `TASK-007`, `TASK-010` | `internal/domain`, `internal/probe` | `POST /api/v1/resources/{id}/checks` | `POS-005` | **Contract example only** — Prism не виконує перевірку |
| **REQ-F-013** | Після перевірки записується окремий JSON-рядок у `stdout`. | Рівно один JSON-рядок із сімома визначеними полями. | `TASK-015` | `TASK-014` | `internal/reporter` | — | — | **Lab5** — HTTP JSON-відповідь не є stdout |
| **REQ-F-014** | `timestamp` має бути RFC3339 у UTC. | Timestamp розбирається як RFC3339 і має UTC-офсет. | `TASK-015` | `TASK-014` | `internal/reporter` | `POST /api/v1/resources/{id}/checks`, `GET /api/v1/resources/{id}/results` (схема) | `POS-003`, `POS-004`, `POS-005` | **Contract schema only** — генерацію timestamp потребує Lab5 |
| **REQ-F-015** | Обробка `SIGINT` і `SIGTERM`. | Активні перевірки скасовуються, нові не плануються, стан не записується. | `TASK-019` | `TASK-018` | `cmd/gouptime` | — | — | **Lab5** — не REST-вимога |
| **REQ-NF-001** | Робота без UI, зовнішньої БД і нотифікацій. | Потрібен лише локальний `config.json`, результати йдуть у `stdout`. | `TASK-020`, `TASK-021` | `TASK-001` | `cmd/gouptime` | — | — | **Lab5** — архітектурне обмеження |
| **REQ-NF-002** | Паралельне виконання перевірок у Go-горутині на ресурс. | Повільний ресурс не блокує інший. | `TASK-017` | `TASK-016` | `internal/scheduler` | — | — | **Lab5** — Prism не має scheduler |
| **REQ-NF-003** | Паралельний `stdout` не пошкоджує JSON-рядки. | Кожен stdout-рядок є самостійним валідним JSON. | `TASK-015` | `TASK-014` | `internal/reporter` | — | — | **Lab5** — не перевіряється HTTP mock |
| **REQ-NF-004** | Перевірка завершується не пізніше `timeout_ms` з незначним overhead. | Перевищення timeout скасовується і формує `DOWN`. | `TASK-011` | `TASK-010` | `internal/probe`, `ManagementAPI` | `POST /api/v1/resources/{id}/checks` (response schema) | `POS-005` | **Contract example only** — час виконання Prism не є доказом |
| **REQ-NF-005** | Оперативний стан зберігається лише в пам'яті процесу. | Після restart не відновлюються попередні стани з файлу/БД. | `TASK-019` | `TASK-018` | `cmd/gouptime`, `ManagementAPI` | `GET /api/v1/resources`, `GET /api/v1/resources/{id}`, `DELETE /api/v1/resources/{id}`, `GET /api/v1/resources/{id}/results` | `POS-002`, `POS-006`, `NEG-006` | **Contract + static mock** — Prism не зберігає стан і не доводить Go-архітектуру |
| **REQ-NF-006** | Валідація конфігурації завершується до створення робочих перевірок. | Для невалідної конфігурації не надсилається жоден HTTP-запит. | `TASK-006` | `TASK-003`, `TASK-005` | `internal/config`, `cmd/gouptime`, `ManagementAPI` | `POST /api/v1/resources` (API validation mapping only) | — | **Lab5** — mock не запускає worker checks і не доводить fail-fast |
| **REQ-NF-007** | Значення статусу машинно-читані та стабільні. | `status` належить `{UP, DEGRADED, DOWN}`. | `TASK-008` | `TASK-007` | `internal/domain`, `ManagementAPI` | `POST /api/v1/resources/{id}/checks`, `GET /api/v1/resources/{id}/results` | `POS-003`, `POS-004`, `POS-005` | **Contract example only** — enum перевірено у схемі, логіку статусу потребує Lab5 |
| **REQ-NF-008** | Сервіс не залежить від response body. | Однакові код і latency дають однаковий статус незалежно від body. | `TASK-011` | `TASK-010` | `internal/probe`, `ManagementAPI` | `POST /api/v1/resources/{id}/checks` | — | **Lab5** — Prism не виконує два реальні GET із різними body |
| **REQ-API-001** | Читання списку та ресурсу за `id`. | List дає `200` з масивом; item дає `200` або `404` з `ErrorResponse`. | — | — | `ManagementAPI` | `GET /api/v1/resources`, `GET /api/v1/resources/{id}` | `POS-002`, `NEG-006` | **PASS контрактної форми; static mock** |
| **REQ-API-002** | Реєстрація ресурсу в пам'яті з дефолтами та Fail-Fast validation. | Валідний запит `201`; validation `400`; invariant `422`; duplicate `409`. | — | — | `ManagementAPI` | `POST /api/v1/resources` | `POS-001`, `NEG-001..005` | **PARTIAL** — `NEG-001/003/004` blocked у Prism; state/validation потребують Lab5 |
| **REQ-API-003** | Видалення ресурсу з пам'яті та припинення періодичних перевірок. | Існуючий `id` дає `204`; невідомий — `404`. | — | — | `ManagementAPI`, `internal/scheduler` | `DELETE /api/v1/resources/{id}` | `POS-006`, `NEG-006` | **PASS контрактної форми; static mock** — state/scheduler потребують Lab5 |
| **REQ-API-004** | Запуск синхронної перевірки та читання останнього результату. | Для існуючого ресурсу `200` з `CheckResult`; відсутній — `404`. | — | — | `ManagementAPI`, `internal/probe` | `POST /api/v1/resources/{id}/checks`, `GET /api/v1/resources/{id}/results` | `POS-003..005`, `NEG-006` | **PASS response shape; behavior потребує Lab5** |
| **REQ-API-005** | Єдиний формат помилкових REST-відповідей без внутрішніх деталей. | `400/404/409/422` повертають `ErrorResponse` без стеків і деталей Go. | — | — | `ManagementAPI` | `POST /api/v1/resources`, `GET /api/v1/resources/{id}`, `DELETE /api/v1/resources/{id}`, `POST /api/v1/resources/{id}/checks`, `GET /api/v1/resources/{id}/results` | `NEG-001..006` | **Contract schema + partial mock** — три status mappings blocked у Prism |

## Перевірка узгодженості з OpenAPI

У `spec/openapi.yaml` кожне `x-srs-requirements` посилається на вимогу, яка:

1. існує у `spec/srs.md`;
2. має REST-представлення в адендумі або є безпосередньою схемною/домен­ною передумовою операції;
3. має відповідний запис у цій матриці.

Наявність ID у `x-srs-requirements` означає трасування контракту, а не доказ виконання вимоги Prism. Зокрема, `REQ-F-013`, `REQ-NF-003`, `REQ-NF-006` і scheduler/probe-вимоги мають бути підтверджені реалізаційними тестами лабораторної №5.

## Вимоги, що потребують лабораторної №5

Реальною Go-реалізацією та тестами потрібно перевірити: `REQ-F-001`, `REQ-F-002`, `REQ-F-003`, `REQ-F-004`, `REQ-F-005`, `REQ-F-006`, `REQ-F-007`, `REQ-F-008`, `REQ-F-009`, `REQ-F-010`, `REQ-F-011`, `REQ-F-012`, `REQ-F-013`, `REQ-F-014`, `REQ-F-015`, `REQ-NF-001`, `REQ-NF-002`, `REQ-NF-003`, `REQ-NF-004`, `REQ-NF-005`, `REQ-NF-006`, `REQ-NF-007`, `REQ-NF-008`, а також поведінкові частини `REQ-API-001..005`.

Mock API лабораторної №4 підтверджує лише те, що описані response schemas/examples та маршрути можуть бути віддані Prism; він не замінює ці реалізаційні тести.
