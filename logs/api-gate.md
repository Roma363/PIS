# Звіт проходження контрольної точки API-GATE

## 1. Metadata контрольної точки

| Поле | Значення |
|---|---|
| **Проєкт** | GoUptime |
| **Дата** | 31.10.2026 |
| **Етап SDD** | `API-GATE` |
| **Лабораторна робота** | №4 — API-First Specification |
| **Гілка** | `feature/api-specification` |
| **Арбітр** | `Human — Roman Benza` |
| **Рішення** | **APPROVED / GO** |

## 2. Чеклист обов'язкових критеріїв API-GATE

| № | Критерій перевірки | Статус | Підтвердження |
|---:|---|:---:|---|
| 1 | OpenAPI 3.0.3-контракт створено та він відповідає затвердженому SRS. | **PASS** | [`spec/openapi.yaml`](../spec/openapi.yaml), [`spec/srs.md`](../spec/srs.md) |
| 2 | Усі необхідні API-операції описані з коректними `operationId`, HTTP-методами та шляхами: `listResources`, `createResource`, `getResourceById`, `deleteResource`, `runResourceCheck`, `getResourceLastResult`. | **PASS** | [`spec/openapi.yaml`](../spec/openapi.yaml), [`docs/api_design.md`](../docs/api_design.md) |
| 3 | Схеми `ResourceCreateRequest`, `Resource`, `CheckResult`, `AvailabilityStatus` та `ErrorResponse` містять обов'язкові поля, типи, обмеження і стабільні enum-значення. | **PASS** | [`spec/openapi.yaml`](../spec/openapi.yaml), [`docs/api_design.md`](../docs/api_design.md) |
| 4 | Дефолтні значення `interval_seconds=10`, `timeout_ms=2000`, `warning_threshold_ms=800` та доменні обмеження відображені в API-контракті. | **PASS** | [`spec/srs.md`](../spec/srs.md), [`spec/openapi.yaml`](../spec/openapi.yaml), [`spec/traceability_matrix.md`](../spec/traceability_matrix.md) |
| 5 | Розмежування помилок API визначене однозначно: `400`, `404`, `409`, `422` з єдиною моделлю `ErrorResponse`. | **PASS** | [`spec/openapi.yaml`](../spec/openapi.yaml), [`docs/api_design.md`](../docs/api_design.md), [`spec/traceability_matrix.md`](../spec/traceability_matrix.md) |
| 6 | OpenAPI YAML пройшов синтаксичну та структурну валідацію рекомендованим валідатором. | **PASS** | [`logs/openapi-validation.log`](openapi-validation.log), [`spec/openapi.yaml`](../spec/openapi.yaml) |
| 7 | Позитивні, BDD та негативні сценарії `POS-001..006`, `BDD-001..003`, `NEG-001..006` описані та пов'язані з операціями API. | **PASS** | [`docs/api_design.md`](../docs/api_design.md), [`spec/traceability_matrix.md`](../spec/traceability_matrix.md), [`logs/mock-api-verification.md`](mock-api-verification.md) |
| 8 | JSON Schema валідація для обов'язкових полів, URL, числових меж, `nullable: false` та доменних інваріантів формалізована в контракті. | **PASS** | [`spec/openapi.yaml`](../spec/openapi.yaml), [`spec/srs.md`](../spec/srs.md), [`logs/mock-api-verification.md`](mock-api-verification.md) |
| 9 | Простежуваність вимог до API-операцій і mock-сценаріїв завершена для `REQ-F-002..014`, `REQ-NF-003..008` та `REQ-API-001..005`. | **PASS** | [`spec/traceability_matrix.md`](../spec/traceability_matrix.md), [`spec/openapi.yaml`](../spec/openapi.yaml) |
| 10 | Контракт зберігає межі системи: in-memory API без БД, веб-інтерфейсу, зовнішніх нотифікацій і запису стану на диск. | **PASS** | [`spec/srs.md`](../spec/srs.md), [`docs/api_design.md`](../docs/api_design.md), [`spec/openapi.yaml`](../spec/openapi.yaml), [`spec/traceability_matrix.md`](../spec/traceability_matrix.md) |

> **Примітка щодо живого Mock API:** протокол містить усі команди та payload для `POS-001..006`, `BDD-001..003` і `NEG-001..006`. На момент виконання контрольної точки `127.0.0.1:4010` не приймав TCP-з'єднання (`connection refused`), тому фактичні HTTP-відповіді Prism не були отримані. Це не змінює PASS для контракту, його валідації та повноти тестового протоколу; live execution має бути повторено після запуску сервера.

## 3. Контроль ризиків

| Ризик | Перевірка mitigation | Результат |
|---|---|---|
| `RISK-001` / `RISK-008` | Перевірено точні межі: `UP` при `latency_ms <= warning_threshold_ms`, `DEGRADED` при `warning_threshold_ms < latency_ms <= timeout_ms`, `DOWN` при `latency_ms > timeout_ms`. Перевірено інваріант `warning_threshold_ms < timeout_ms`. | **PASS** — оператори зафіксовані у [`spec/srs.md`](../spec/srs.md), [`spec/openapi.yaml`](../spec/openapi.yaml) та [`docs/api_design.md`](../docs/api_design.md). |
| `RISK-013` | Перевірено відхилення `null` через `nullable: false` та значень `<= 0` через `minimum: 1` у JSON Schema OpenAPI; негативні payload винесені у `NEG-001`, `NEG-003`, `NEG-004`, `NEG-005`. | **PASS** — [`spec/openapi.yaml`](../spec/openapi.yaml), [`docs/api_design.md`](../docs/api_design.md), [`logs/mock-api-verification.md`](mock-api-verification.md). |
| `RISK-015` | Підтверджено відсутність Scope Creep: немає БД, зовнішніх нотифікацій, веб-інтерфейсу чи збереження стану на диск. | **PASS** — межі зафіксовані у [`spec/srs.md`](../spec/srs.md), [`docs/api_design.md`](../docs/api_design.md) та [`spec/openapi.yaml`](../spec/openapi.yaml). |

## 4. Аудит пропозицій AI-агента та корективи людини

У межах Human-in-the-Loop аудиту зафіксовано такі рішення:

- Людина заборонила додавання ендпоінта для збереження конфігурації на диск (`PUT /config/save`) та історії перевірок у базі даних. Межі `REQ-NF-001` і `REQ-NF-005` залишено без змін.
- Людина скоригувала схему `ResourceCreateRequest`, додавши явну заборону `nullable: false` для полів, де `null` не допускається, відповідно до `REQ-F-006`.
- Людина затвердила розмежування кодів помилок: `400 Bad Request` для синтаксису, типів і значень `<=0`; `409 Conflict` для дубліката `id` за `REQ-F-003`; `422 Unprocessable Entity` для порушення `warning_threshold_ms >= timeout_ms` за `REQ-F-006`.
- Вихідний код у `/src` та тести у `/tests` не змінювалися. Реалізацію відкладено до Лабораторної роботи №5 відповідно до TDD-послідовності.

## 5. Резолюція Арбітра

### Рішення

**GO — дозволити перехід до Лабораторної роботи №5 (TDD-реалізація).**

API-контракт, JSON Schema-обмеження, трасування вимог, контроль Scope та Human-in-the-Loop корективи визнані достатніми для переходу до реалізації. Реалізація API, scheduler, probe, reporter і тестів має виконуватися на наступному етапі без розширення затвердженого Scope.

### Підпис арбітра

**Human — Roman Benza**