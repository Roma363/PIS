# Mock API Verification Protocol

**Base URL:** `http://127.0.0.1:4010/api/v1`
**Specification:** `spec/openapi.yaml`
**Executed:** 2026-10-31T23:13:11.4572442+02:00
**Overall result:** BLOCKED — the mock server was not listening on `127.0.0.1:4010`. No scenario reached HTTP or JSON Schema validation.

Connectivity evidence: `Test-NetConnection 127.0.0.1 -Port 4010` returned `TcpTestSucceeded: False`; `curl.exe` returned exit code `7` and `Failed to connect to 127.0.0.1 port 4010 ... Connection refused`. Therefore HTTP status, response headers, and JSON response bodies are unavailable for every scenario.

## POS-001 — Create resource with defaults

**Expected:** 201 Created; defaults 10/2000/800

### Request

```text
curl -i -s -X POST http://127.0.0.1:4010/api/v1/resources
  -H "Accept: application/json"
  -H "Prefer: code=201"
  -H "Prefer: example=CreatedWithDefaults"
  -H "Content-Type: application/json"
  --data-raw '{"id":"kpi-main","url":"https://kpi.ua"}'
```

### Actual response

```text
HTTP status: unavailable (connection refused before an HTTP response)
HTTP headers: unavailable
JSON body: unavailable
curl exit code: 7
curl error: Failed to connect to 127.0.0.1 port 4010: Connection refused
```

**Scenario result:** NOT EXECUTED against Prism because the server was unavailable.

---

## POS-002 — List resources

**Expected:** 200 OK; JSON array of resources

### Request

```text
curl -i -s -X GET http://127.0.0.1:4010/api/v1/resources
  -H "Accept: application/json"
  -H "Prefer: code=200"
  -H "Prefer: example=ResourcesList"
```

### Actual response

```text
HTTP status: unavailable (connection refused before an HTTP response)
HTTP headers: unavailable
JSON body: unavailable
curl exit code: 7
curl error: Failed to connect to 127.0.0.1 port 4010: Connection refused
```

**Scenario result:** NOT EXECUTED against Prism because the server was unavailable.

---

## POS-003 / BDD-001 — Check resource - UP

**Expected:** 200 OK; UP, http_code 200, latency 650, empty error

### Request

```text
curl -i -s -X POST http://127.0.0.1:4010/api/v1/resources/kpi-main/checks
  -H "Accept: application/json"
  -H "Prefer: code=200"
  -H "Prefer: example=BDD-001-UP"
```

### Actual response

```text
HTTP status: unavailable (connection refused before an HTTP response)
HTTP headers: unavailable
JSON body: unavailable
curl exit code: 7
curl error: Failed to connect to 127.0.0.1 port 4010: Connection refused
```

**Scenario result:** NOT EXECUTED against Prism because the server was unavailable.

---

## POS-004 / BDD-002 — Check resource after redirect - DEGRADED

**Expected:** 200 OK; DEGRADED, http_code 200, latency 1200, empty error

### Request

```text
curl -i -s -X POST http://127.0.0.1:4010/api/v1/resources/kpi-main/checks
  -H "Accept: application/json"
  -H "Prefer: code=200"
  -H "Prefer: example=BDD-002-DEGRADED"
```

### Actual response

```text
HTTP status: unavailable (connection refused before an HTTP response)
HTTP headers: unavailable
JSON body: unavailable
curl exit code: 7
curl error: Failed to connect to 127.0.0.1 port 4010: Connection refused
```

**Scenario result:** NOT EXECUTED against Prism because the server was unavailable.

---

## POS-005 / BDD-003 — Check resource timeout - DOWN

**Expected:** 200 OK; DOWN, http_code 0, timeout error

### Request

```text
curl -i -s -X POST http://127.0.0.1:4010/api/v1/resources/kpi-main/checks
  -H "Accept: application/json"
  -H "Prefer: code=200"
  -H "Prefer: example=BDD-003-DOWN-Timeout"
```

### Actual response

```text
HTTP status: unavailable (connection refused before an HTTP response)
HTTP headers: unavailable
JSON body: unavailable
curl exit code: 7
curl error: Failed to connect to 127.0.0.1 port 4010: Connection refused
```

**Scenario result:** NOT EXECUTED against Prism because the server was unavailable.

---

## POS-006 — Delete resource

**Expected:** 204 No Content; empty body

### Request

```text
curl -i -s -X DELETE http://127.0.0.1:4010/api/v1/resources/kpi-main
  -H "Accept: application/json"
  -H "Prefer: code=204"
```

### Actual response

```text
HTTP status: unavailable (connection refused before an HTTP response)
HTTP headers: unavailable
JSON body: unavailable
curl exit code: 7
curl error: Failed to connect to 127.0.0.1 port 4010: Connection refused
```

**Scenario result:** NOT EXECUTED against Prism because the server was unavailable.

---

## NEG-001 — Empty required id

**Expected:** 400 Bad Request; VALIDATION_ERROR

### Request

```text
curl -i -s -X POST http://127.0.0.1:4010/api/v1/resources
  -H "Accept: application/json"
  -H "Prefer: code=400"
  -H "Prefer: example=EmptyRequiredField"
  -H "Content-Type: application/json"
  --data-raw '{"id":"","url":"https://kpi.ua"}'
```

### Actual response

```text
HTTP status: unavailable (connection refused before an HTTP response)
HTTP headers: unavailable
JSON body: unavailable
curl exit code: 7
curl error: Failed to connect to 127.0.0.1 port 4010: Connection refused
```

**Scenario result:** NOT EXECUTED against Prism because the server was unavailable.

---

## NEG-002 — Duplicate resource id

**Expected:** 409 Conflict; DUPLICATE_RESOURCE_ID

### Request

```text
curl -i -s -X POST http://127.0.0.1:4010/api/v1/resources
  -H "Accept: application/json"
  -H "Prefer: code=409"
  -H "Prefer: example=DuplicateResourceID"
  -H "Content-Type: application/json"
  --data-raw '{"id":"kpi-main","url":"https://kpi.ua"}'
```

### Actual response

```text
HTTP status: unavailable (connection refused before an HTTP response)
HTTP headers: unavailable
JSON body: unavailable
curl exit code: 7
curl error: Failed to connect to 127.0.0.1 port 4010: Connection refused
```

**Scenario result:** NOT EXECUTED against Prism because the server was unavailable.

---

## NEG-003 — Unsupported URL scheme

**Expected:** 400 Bad Request; VALIDATION_ERROR

### Request

```text
curl -i -s -X POST http://127.0.0.1:4010/api/v1/resources
  -H "Accept: application/json"
  -H "Prefer: code=400"
  -H "Prefer: example=InvalidURLScheme"
  -H "Content-Type: application/json"
  --data-raw '{"id":"ftp-srv","url":"ftp://files.example.com/pub"}'
```

### Actual response

```text
HTTP status: unavailable (connection refused before an HTTP response)
HTTP headers: unavailable
JSON body: unavailable
curl exit code: 7
curl error: Failed to connect to 127.0.0.1 port 4010: Connection refused
```

**Scenario result:** NOT EXECUTED against Prism because the server was unavailable.

---

## NEG-004 — Non-positive timeout

**Expected:** 400 Bad Request; VALIDATION_ERROR

### Request

```text
curl -i -s -X POST http://127.0.0.1:4010/api/v1/resources
  -H "Accept: application/json"
  -H "Prefer: code=400"
  -H "Prefer: example=NullOrNonPositiveValue"
  -H "Content-Type: application/json"
  --data-raw '{"id":"bad-timeout","url":"https://kpi.ua","timeout_ms":0}'
```

### Actual response

```text
HTTP status: unavailable (connection refused before an HTTP response)
HTTP headers: unavailable
JSON body: unavailable
curl exit code: 7
curl error: Failed to connect to 127.0.0.1 port 4010: Connection refused
```

**Scenario result:** NOT EXECUTED against Prism because the server was unavailable.

---

## NEG-005 — Warning threshold equals timeout

**Expected:** 422 Unprocessable Entity; CONFIG_ERROR

### Request

```text
curl -i -s -X POST http://127.0.0.1:4010/api/v1/resources
  -H "Accept: application/json"
  -H "Prefer: code=422"
  -H "Prefer: example=ThresholdExceedsTimeout"
  -H "Content-Type: application/json"
  --data-raw '{"id":"bad-threshold","url":"https://kpi.ua","timeout_ms":2000,"warning_threshold_ms":2000}'
```

### Actual response

```text
HTTP status: unavailable (connection refused before an HTTP response)
HTTP headers: unavailable
JSON body: unavailable
curl exit code: 7
curl error: Failed to connect to 127.0.0.1 port 4010: Connection refused
```

**Scenario result:** NOT EXECUTED against Prism because the server was unavailable.

---

## NEG-006 — Unknown resource

**Expected:** 404 Not Found; RESOURCE_NOT_FOUND

### Request

```text
curl -i -s -X GET http://127.0.0.1:4010/api/v1/resources/unknown-id
  -H "Accept: application/json"
  -H "Prefer: code=404"
  -H "Prefer: example=ResourceNotFound"
```

### Actual response

```text
HTTP status: unavailable (connection refused before an HTTP response)
HTTP headers: unavailable
JSON body: unavailable
curl exit code: 7
curl error: Failed to connect to 127.0.0.1 port 4010: Connection refused
```

**Scenario result:** NOT EXECUTED against Prism because the server was unavailable.

---

