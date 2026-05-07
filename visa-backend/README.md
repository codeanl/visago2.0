# Visago Backend

Go backend API for Visago.

## Run

```powershell
go run ./cmd/server
```

Default URL: `http://localhost:8080`

## Docker Compose

1. You can start directly with defaults, or copy `.env.example` to `.env` first if you want to customize secrets and ports.
2. Start services:

```powershell
docker compose up -d --build
```

3. Check backend health:

```powershell
Invoke-WebRequest -UseBasicParsing http://127.0.0.1:8081/api/health
```

Notes:

- Compose includes `mysql` and `backend`.
- The Docker image sources are switched to the Daocloud mirror, so `docker compose up -d --build` is less likely to fail on Docker Hub timeouts in mainland China networks.
- The backend is exposed to the host on `127.0.0.1:8081` by default in Compose, so it can run alongside a locally started backend on `8080`.
- MySQL is exposed to the host on `127.0.0.1:3307` by default, while the backend still connects to `mysql:3306` inside Compose.
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
