# Перевірка Mock API

**Дата звіту:** 2026-11-02
**Дата виконання mock-запитів:** 2026-11-02
**Середовище:** Windows PowerShell, Node.js `v24.19.0`, npm `11.17.0`, Stoplight Prism CLI `5.16.0`, curl `8.21.0`
**Специфікація:** `spec/openapi.yaml`
**Базова адреса:** `http://127.0.0.1:4010/api/v1`

Цей звіт перевіряє лише mock-контракт OpenAPI. Prism повертає статичні приклади; він не є реалізацією Go API.

## Відтворюване налаштування PowerShell

Запустіть Prism в одному вікні PowerShell:

```powershell
npx --yes @stoplight/prism-cli mock .\spec\openapi.yaml --host 127.0.0.1 --port 4010
```

Перевірте доступність в іншому вікні PowerShell:

```powershell
Test-NetConnection 127.0.0.1 -Port 4010
```

Для кожного POST-запиту створюйте UTF-8-файл без додаткового екранування та передавайте його через `--data-binary "@file"`:

```powershell
$bodyPath = Join-Path $env:TEMP 'gouptime-resource.json'
[IO.File]::WriteAllText($bodyPath, '{"id":"kpi-main","url":"https://kpi.ua"}', (New-Object Text.UTF8Encoding($false)))
curl.exe --max-time 5 -i -sS -X POST 'http://127.0.0.1:4010/api/v1/resources' `
  -H 'Content-Type: application/json' `
  -H 'Prefer: code=201,example=CreatedWithDefaults' `
  --data-binary "@$bodyPath"
Remove-Item $bodyPath -Force
```

Bash-приклади в `docs/api_design.md` явно позначені як `bash` і використовують `--max-time 5`. Не копіюйте їхній inline-синтаксис `-d` у PowerShell.

## Категорії доказів

| Доказ | Що підтверджує | Чого не підтверджує |
|---|---|---|
| Redocly lint | Структурна валідність документа OpenAPI | Runtime-поведінку |
| Запуск Prism і TCP-перевірка | Mock-сервер прослуховує порт | Бізнес-вимоги |
| HTTP-відповідь/приклад | Доступність маршруту, статусу та форми відповіді | Go-логіку, збереження або зміну стану |
| `sl-violations` | Prism виявив порушення схеми запиту | Відображення статусу в реальному API |

## Шість операцій

| Метод і шлях | Очікуваний результат контракту | Фактичний результат | Статус |
|---|---|---|:---:|
| `GET /api/v1/resources` | `200`, масив ресурсів | `200`, статичний JSON-масив | **PASS — маршрут/приклад** |
| `POST /api/v1/resources` | `201`, ресурс із дефолтами | `201`, приклад `CreatedWithDefaults` через файл | **PASS — маршрут/приклад** |
| `GET /api/v1/resources/kpi-main` | `200`, ресурс | `200`, приклад `Resource` | **PASS — маршрут/приклад** |
| `DELETE /api/v1/resources/kpi-main` | `204`, порожнє тіло | `204`, порожнє тіло | **PASS — маршрут/приклад** |
| `POST /api/v1/resources/kpi-main/checks` | `200`, `CheckResult` | `200`, приклад `BDD-001-UP` | **PASS — маршрут/приклад** |
| `GET /api/v1/resources/kpi-main/results` | `200`, `CheckResult` | `200`, приклад `LastResultUP` | **PASS — маршрут/приклад** |

Ці запити не доводять, що ресурс було створено, видалено, збережено між запитами або фактично перевірено. Prism не підтримує стан ресурсів GoUptime.

## POS- та BDD-сценарії

| Сценарій | Метод/шлях і вхідні дані | Очікуваний результат | Фактичний результат | Статус |
|---|---|---|---|:---:|
| POS-001 | `POST /resources`; `{"id":"kpi-main","url":"https://kpi.ua"}` | `201`; дефолти `10/2000/800` | `201`; відповідний статичний приклад | **PASS — лише приклад** |
| POS-002 | `GET /resources` | `200`; масив ресурсів | `200`; відповідний статичний приклад | **PASS — лише приклад** |
| POS-003 / BDD-001-UP | `POST /resources/kpi-main/checks`; `Prefer: ...BDD-001-UP` | `200`; `UP`, код `200`, latency `650` | Відповідний статичний приклад | **PASS — лише приклад** |
| POS-004 / BDD-002-DEGRADED | `POST /resources/kpi-main/checks`; `Prefer: ...BDD-002-DEGRADED` | `200`; `DEGRADED`, latency `1200` | Відповідний статичний приклад | **PASS — лише приклад** |
| POS-005 / BDD-003-DOWN-Timeout | `POST /resources/kpi-main/checks`; `Prefer: ...BDD-003-DOWN-Timeout` | `200`; `DOWN`, код `0`, помилка timeout | Відповідний статичний приклад | **PASS — лише приклад** |
| POS-006 | `DELETE /resources/kpi-main` | `204`; порожнє тіло | Відповідний статичний приклад | **PASS — лише приклад** |

Приклади не доводять реальну класифікацію latency, обробку редиректів, скасування за timeout, мережеві перевірки або збереження результатів.

## NEG-сценарії

Усі POST-тіла передавалися з UTF-8-файлів через `curl.exe --max-time 5 --data-binary "@file"`.

| Сценарій | Payload | Очікуваний результат контракту | Фактичний результат Prism | Статус |
|---|---|---|---|:---:|
| NEG-001 | `{"id":"","url":"https://kpi.ua"}` | `400`, `VALIDATION_ERROR` | Prism вибрав `422` і приклад `ThresholdExceedsTimeout`; примусовий `code=400,example=EmptyRequiredField` повернув `404` Prism | **BLOCKED** |
| NEG-002 | `{"id":"kpi-main","url":"https://kpi.ua"}` | `409`, `DUPLICATE_RESOURCE_ID` | `409`, статичний приклад дубліката | **PASS — лише приклад** |
| NEG-003 | `{"id":"ftp-srv","url":"ftp://files.example.com/pub"}` | `400`, `VALIDATION_ERROR` | Prism не відтворив очікуваний `400`, аналогічно до NEG-001 | **BLOCKED** |
| NEG-004 | `{"id":"bad-timeout","url":"https://kpi.ua","timeout_ms":0}` | `400`, `VALIDATION_ERROR` | Prism не відтворив очікуваний `400`, аналогічно до NEG-001 | **BLOCKED** |
| NEG-005 | `{"id":"bad-threshold","url":"https://kpi.ua","timeout_ms":2000,"warning_threshold_ms":2000}` | `422`, `CONFIG_ERROR` | `422`, статичний приклад порушення threshold | **PASS — лише приклад** |
| NEG-006 | Без body; `GET /resources/unknown-id` | `404`, `RESOURCE_NOT_FOUND` | `404`, статичний приклад відсутнього ресурсу | **PASS — лише приклад** |

Очікувані коди `400` у контракті не змінювалися. Prism `5.16.0` не забезпечує достатньої поведінки узгодження, щоб їх підтвердити. Міжполе́вий інваріант `warning_threshold_ms < timeout_ms` описаний у контракті, але не може бути виражений двома незалежними ключовими словами `minimum` і потребує окремого contract/runtime-тесту.

## Валідація OpenAPI

Команду успішно виконано 2026-11-02:

```powershell
npx --yes @redocly/cli lint .\spec\openapi.yaml
```

Результат: **PASS**, Redocly повідомив `spec\openapi.yaml: validated in 71ms`.

`logs/openapi-validation.log` відсутній у репозиторії. Щоб зберегти результат:

```powershell
npx --yes @redocly/cli lint .\spec\openapi.yaml 2>&1 | Tee-Object .\logs\openapi-validation.log
```

## Обмеження та наступні дії

- Prism є статичним mock-сервером і не реалізує реєстрацію, унікальність, видалення, scheduler, перевірки ресурсів, редиректи, контроль timeout або JSONL у stdout.
- `NEG-001`, `NEG-003` та `NEG-004` залишаються **BLOCKED**.
- У лабораторній №4 блок можна усунути лише окремим contract-тестом/валідатором очікуваних відображень статусів; це не потребує реалізації Go-сервера.
- В іншому разі відображення статусів потрібно перевірити на реальному Go API у лабораторній №5.
- Не видаляйте приклади, не послаблюйте валідацію, не змінюйте SRS і не збільшуйте тайм-аути, щоб приховати це обмеження.
