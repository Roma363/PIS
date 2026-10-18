# Матриця простежуваності вимог (Traceability Matrix)

| ID вимоги | Опис вимоги | Критерій прийняття | BDD-сценарій | ID задачі | ID тесту | Код/Модуль |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| **REQ-F-001** | Завантаження з config.json | Тільки локальний config.json; вихід при збої | — | `TASK-004` | `TASK-003` | `internal/config` |
| **REQ-F-002** | Обов'язкові поля id та url | ConfigError при порожніх або відсутніх | — | `TASK-006` | `TASK-005` | `internal/config` |
| **REQ-F-003** | Унікальність ідентифікаторів | ConfigError при дублікатах id | — | `TASK-006` | `TASK-005` | `internal/config` |
| **REQ-F-004** | Валідність схеми HTTP/HTTPS | ConfigError при некоректному синтаксисі/схемі | — | `TASK-006` | `TASK-005` | `internal/config` |
| **REQ-F-005** | Дефолтні числові значення | Застосування дефолтів (10с, 2000мс, 800мс) | — | `TASK-006` | `TASK-005` | `internal/config` |
| **REQ-F-006** | Fail-Fast для null та <= 0 | ConfigError до першого HTTP-запиту | — | `TASK-006` | `TASK-005` | `internal/config` |
| **REQ-F-007** | Виконання HTTP GET без тіла | Рівно 1 GET-запит на перевірку, ігнор body | BDD-001 | `TASK-011` | `TASK-010` | `internal/probe` |
| **REQ-F-008** | Слідування до 3 редиректів | FinalResponse; перевищення/цикл -> DOWN (код 0) | BDD-002, BDD-003 | `TASK-013` | `TASK-012` | `internal/probe` |
| **REQ-F-009** | Незалежні інтервали | Stateless-перевірка без впливу історії | — | `TASK-017` | `TASK-016` | `internal/scheduler` |
| **REQ-F-010** | Критерій статусу UP | Код [200,400) та Latency <= WarningThreshold | BDD-001 | `TASK-008` | `TASK-007` | `internal/domain` |
| **REQ-F-011** | Критерій статусу DEGRADED | Код [200,400) та WarningThreshold < Latency <= Timeout | BDD-002 | `TASK-008` | `TASK-007` | `internal/domain` |
| **REQ-F-012** | Критерій статусу DOWN | Помилка, таймаут, Latency > Timeout, код >= 400 | BDD-003 | `TASK-008`, `TASK-011` | `TASK-007`, `TASK-010` | `internal/domain`, `internal/probe` |
| **REQ-F-013** | JSON-вивід у stdout | Один рядок JSON з 7 обов'язковими полями | BDD-001 | `TASK-015` | `TASK-014` | `internal/reporter` |
| **REQ-F-014** | RFC3339 UTC timestamp | Час у форматі ISO/RFC3339 з UTC-офсетом | — | `TASK-015` | `TASK-014` | `internal/reporter` |
| **REQ-F-015** | Graceful Shutdown | Скасування через контекст без збереження стану | — | `TASK-019` | `TASK-018` | `cmd/gouptime` |
| **REQ-NF-001**| Відсутність БД та UI | Робота виключно з config.json та stdout | — | `TASK-020`, `TASK-021` | `TASK-001` | `cmd/gouptime` |
| **REQ-NF-002**| Паралельність горутин | Ізоляція перевірок між ресурсами | — | `TASK-017` | `TASK-016` | `internal/scheduler` |
| **REQ-NF-003**| Атомарність рядків JSON | Валідний незалежний JSON на кожен рядок | BDD-001 | `TASK-015` | `TASK-014` | `internal/reporter` |
| **REQ-NF-004**| Дотримання таймаутів | Завершення не пізніше TimeoutMs | BDD-003 | `TASK-011` | `TASK-010` | `internal/probe` |
| **REQ-NF-005**| In-memory архітектура | Відсутність залишкового стану після рестарту | — | `TASK-019` | `TASK-018` | `cmd/gouptime` |
| **REQ-NF-006**| Валідація до запуску | Жодного запиту при збої валідації конфігу | — | `TASK-006` | `TASK-003`, `TASK-005` | `internal/config` |
| **REQ-NF-007**| Стабільна множина статусів | Строго { "UP", "DEGRADED", "DOWN" } | BDD-001..003 | `TASK-008` | `TASK-007` | `internal/domain` |
| **REQ-NF-008**| Незалежність від тіла | Однаковий вердикт незалежно від payload | — | `TASK-011` | `TASK-010` | `internal/probe` |