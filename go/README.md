# Go Fiber Courier API

## Stack

- Go
- Fiber
- MongoDB official Go driver
- go-playground/validator

## Environment

Copy `.env.example` values into your environment:

```text
APP_PORT=8080
MONGODB_URI=mongodb://mongodb:mongodb12345@localhost:27019/?authSource=admin
MONGODB_DATABASE=gradin-courier
```

## Run

```bash
go mod tidy
go run ./cmd/api
```

With Docker Compose:

```bash
docker compose up --build go
```

Local URL: `http://localhost:8080`.

## Tests

```bash
go test ./...
go vet ./...
gofmt -w .
```

Tests need MongoDB available through `MONGODB_URI`. Courier documents are stored in the shared `couriers` collection.

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
