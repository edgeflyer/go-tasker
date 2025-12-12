# API Contract

- Base URL: `http://localhost:8080`
- Success response wrapper: `{"data": <payload>}`
- Error response wrapper: `{"error": {"code": "<CODE>", "message": "<human readable>"}}`
- Auth: Bearer JWT (`Authorization: Bearer <token>`) required for all `/tasks` routes. Token is obtained via `/auth/login`, expires in 2 hours.
- Task status values: `pending` or `completed`.

## Public endpoints

- `GET /health` — 200 → `{"data":{"status":"ok"}}`
- `GET /ping` — 200 → `{"data":{"message":"pong"}}`
- `GET /demo-error` — always 400 → `{"error":{"code":"DEMO_ERROR","message":"this is a demo error"}}`

## Auth

- `POST /auth/register`

  - Body: `{"username": "string", "password": "string (>=6 chars)"}`
  - 201 → `{"data":{"user": { "id": number, "username": string, "created_at": RFC3339, "updated_at": RFC3339 }}}`
  - Errors: 400 `INVALID_JSON`; 409 `INVALID_USERNAME`/`INVALID_PASSWORD`/`USERNAME_EXISTS`; 500 `INTERNAL_ERROR`.

- `POST /auth/login`
  - Body: `{"username": "string", "password": "string"}`
  - 200 → `{"data":{"token": "jwt", "user": { "id": number, "username": string, "created_at": RFC3339, "updated_at": RFC3339 }}}`
  - Errors: 400 `INVALID_JSON` or other app errors; 401 `INVALID_CREDENTIALS`; 500 `INTERNAL_ERROR` or `TOKEN_ERROR`.

## Tasks (protected, require `Authorization: Bearer <token>`)

- `POST /tasks`

  - Body: `{"title": "string (required)", "description": "string"}`
  - 201 → `{"data": { "id": number, "user_id": number, "title": string, "description": string, "status": "pending", "created_at": RFC3339, "updated_at": RFC3339 }}`
  - Errors: 400 `INVALID_JSON`/`INVALID_TITLE`; 500 `INTERNAL_ERROR`.

- `GET /tasks`

  - 200 → `{"data": [ { "id": number, "user_id": number, "title": string, "description": string, "status": "pending|completed", "created_at": RFC3339, "updated_at": RFC3339 }, ... ]}`
  - Errors: 401 `UNAUTHORIZED`; 500 `INTERNAL_ERROR`.

- `GET /tasks/:id`

  - Params: `id` path param (positive integer)
  - 200 → `{"data": { "id": number, "user_id": number, "title": string, "description": string, "status": "pending|completed", "created_at": RFC3339, "updated_at": RFC3339 }}`
  - Errors: 400 `INVALID_ID`; 401 `UNAUTHORIZED`; 404 `TASK_NOT_FOUND`; 500 `INTERNAL_ERROR`.

- `PUT /tasks/:id`

  - Params: `id` path param (positive integer)
  - Body: `{"title": "string (required)", "description": "string", "status": "pending|completed"}`
  - 200 → `{"data": { "id": number, "user_id": number, "title": string, "description": string, "status": "pending|completed", "created_at": RFC3339, "updated_at": RFC3339 }}`
  - Errors: 400 `INVALID_ID`/`INVALID_JSON`/`INVALID_TITLE`/`INVALID_STATUS`; 401 `UNAUTHORIZED`; 404 `TASK_NOT_FOUND`; 500 `INTERNAL_ERROR`.

- `DELETE /tasks/:id`
  - Params: `id` path param (positive integer)
  - 200 → `{"data":{"message":"task deleted"}}`
  - Errors: 400 `INVALID_ID`; 401 `UNAUTHORIZED`; 404 `TASK_NOT_FOUND`; 500 `INTERNAL_ERROR`.
