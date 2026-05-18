# Laravel Courier API

## Stack

- Laravel
- MongoDB via `mongodb/laravel-mongodb`
- PHPUnit feature tests

## Environment

```text
APP_PORT=8000
MONGODB_URI=mongodb://mongodb:mongodb12345@localhost:27019/?authSource=admin
MONGODB_DATABASE=gradin-courier
```

## Run

```bash
composer install
php artisan serve --host=0.0.0.0 --port=8000
```

With Docker Compose:

```bash
docker compose up --build laravel
```

Local URL: `http://localhost:8000`.

## Tests

```bash
php artisan test
./vendor/bin/pint --test
```

Tests need MongoDB available through `MONGODB_URI`. Courier documents are stored in the shared `couriers` collection.

## API Examples

```text
GET    /api/couriers
GET    /api/couriers/{courier}
POST   /api/couriers
PUT    /api/couriers/{courier}
DELETE /api/couriers/{courier}
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
