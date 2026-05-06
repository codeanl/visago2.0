# Visago Backend

Go backend API for Visago.

## Run

```powershell
go run ./cmd/server
```

Default URL: `http://localhost:8080`

## Docker Compose

1. Copy `.env.example` to `.env` and modify secrets if needed.
2. Start services:

```powershell
docker compose up -d --build
```

3. Check backend health:

```powershell
Invoke-WebRequest -UseBasicParsing http://127.0.0.1:8080/api/health
```

Notes:

- Compose includes `mysql` and `backend`.
- The backend uses `SKIP_DB_BOOTSTRAP=true` in Compose, so MySQL database/user are created by the official MySQL image environment variables.
- Uploaded files are persisted in `./uploads`.
- Image upload endpoints now target Qiniu Object Storage when `QINIU_ACCESS_KEY / QINIU_SECRET_KEY / QINIU_BUCKET / QINIU_DOMAIN` are configured.

## API

- `GET /api/health`
- `GET /api/users`
- `POST /api/users`
- `PUT /api/users/{id}`
- `DELETE /api/users/{id}`
- `GET /api/visa/countries`
- `POST /api/visa/countries`
- `PUT /api/visa/countries/{id}`
- `DELETE /api/visa/countries/{id}`
- `GET /api/visa/countries/{id}/visas`
- `POST /api/visa/countries/{id}/visas`
- `GET /api/visa/country-visas`
- `POST /api/visa/country-visas`
- `PUT /api/visa/country-visas/{id}`
- `DELETE /api/visa/country-visas/{id}`
