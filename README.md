# Courier Technical Test

This repository contains three backend implementations of the same Courier/Kurir CRUD technical test:

- `laravel/`
- `go/`
- `typescript/`

The business flow follows the original Laravel-oriented brief in `BRD.md`, while the public repository excludes `BRD.md`, `AGENTS.md`, and `.ai/` from git.

## Stack Overview

All implementations use MongoDB for Courier persistence and expose the same public API shape under:

```text
/api/couriers
```

| Implementation | Runtime | Local URL |
| --- | --- | --- |
| Laravel | PHP Laravel with MongoDB | `http://localhost:8000` |
| Go | Go Fiber with MongoDB driver | `http://localhost:8080` |
| TypeScript | Fastify, Prisma, MongoDB | `http://localhost:3001` |

## Docker Compose

Start all services:

```bash
docker compose up --build
```

MongoDB is exposed on `localhost:27019` with this local connection:

```text
mongodb://mongodb:mongodb12345@localhost:27019/
```

For the TypeScript Prisma implementation, MongoDB runs as replica set `rs0`; use `replicaSet=rs0&directConnection=true` in local Prisma `DATABASE_URL`.

All backends use the same database and collection:

- database: `gradin-courier`
- collection: `couriers`

The MongoDB container runs `mongodb/init/01-create-couriers.js` on first initialization to create the `gradin-courier` database, the `couriers` collection, and shared indexes for `name`, `level`, `registered_at`, and sparse unique `email`.

## Run Individually

Laravel:

```bash
cd laravel
composer install
php artisan test
php artisan serve --host=0.0.0.0 --port=8000
```

Go Fiber:

```bash
cd go
go mod tidy
go test ./...
go run ./cmd/api
```

TypeScript Fastify:

```bash
cd typescript
npm install
npx prisma generate
npx prisma db push
npm run test
npm run build
npm run dev
```

## API

```text
GET    /api/couriers
GET    /api/couriers/:id
POST   /api/couriers
PUT    /api/couriers/:id
DELETE /api/couriers/:id
```

Index supports:

- `page`
- `per_page`
- `sort=name|registered_at|created_at`
- `sort=-registered_at` for descending sort
- `search=budi+agung`
- `level=2,3`

Example create payload:

```json
{
  "name": "Budiono Hadi Agung",
  "email": "budi@example.com",
  "phone": "08123456789",
  "level": 2,
  "vehicle_type": "motorcycle",
  "license_plate": "L 1234 AB",
  "status": "active"
}
```
