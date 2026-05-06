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

## Modules

- 用户管理：用户列表、搜索、新增、编辑、删除
- 签证管理：国家管理、从国家进入该国家签证管理、国家签证新增/编辑/删除
