# TypeScript Fastify Courier API

## Stack

- TypeScript
- Fastify
- Prisma `6.19.3`
- MongoDB
- Zod
- Vitest

## Environment

```text
APP_PORT=3001
MONGODB_URI=mongodb://mongodb:mongodb12345@localhost:27019/?authSource=admin
MONGODB_DATABASE=gradin-courier
DATABASE_URL=mongodb://mongodb:mongodb12345@localhost:27019/gradin-courier?authSource=admin&replicaSet=rs0&directConnection=true
```

## Run

```bash
npm install
npx prisma generate
npx prisma db push
npm run dev
```

With Docker Compose:

```bash
docker compose up --build typescript
```

Local URL: `http://localhost:3001`.

## Tests

```bash
npm run test
npm run build
npm run lint
```

## API Examples

```text
GET    /api/couriers
GET    /api/couriers/:id
POST   /api/couriers
PUT    /api/couriers/:id
DELETE /api/couriers/:id
```

Create:

```json
{"name":"Budiono Hadi Agung","email":"budi@example.com","level":2,"vehicle_type":"motorcycle"}
```

Update:

```json
{"name":"Budi Agung","level":3,"status":"active"}
```

Queries:

```text
/api/couriers?page=1&per_page=10
/api/couriers?search=budi+agung
/api/couriers?level=2,3
/api/couriers?sort=registered_at
/api/couriers?sort=-registered_at
```

Delete uses soft delete. Deleted couriers do not appear in normal list/show results.
