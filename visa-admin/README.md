# Visago Admin

Vue 3 + Element Plus admin panel for Visago.

## Run

Start backend first:

```powershell
cd ..\visa-backend
go run ./cmd/server
```

Then start admin:

```powershell
npm install
npm run dev
```

Default URL: `http://localhost:5174`

The Vite dev server proxies `/api` to `http://localhost:8080`.

## Docker

Build image:

```powershell
docker build -t visago-admin:latest .
```

Run on port `10010`:

```powershell
docker run -d --name visago-admin -p 10010:10010 visago-admin:latest
```

Optional build arg for API base:

```powershell
docker build -t visago-admin:latest --build-arg VITE_API_BASE=https://visagoapi.nova2026.top/api .
```

Or use Docker Compose:

```powershell
docker compose up -d --build
```

## Modules

- 用户管理：用户列表、搜索、新增、编辑、删除
- 签证管理：国家管理、从国家进入该国家签证管理、国家签证新增/编辑/删除
