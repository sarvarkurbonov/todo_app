# Todo REST API (Go + Gin + PostgreSQL + Redis)

A minimal REST API with **JWT authentication**, **PostgreSQL** storage, and **Redis** caching.

---

## Features
- User **sign-up / sign-in** (JWT)
- CRUD for **todo lists** and **items**
- **Redis cache** for item lookups
- **Plain SQL migrations** (`migrations/`) — no ORM required

---

## Tech Stack
- Go (Gin)
- PostgreSQL
- Redis
- SQL migrations (up/down files)

---

## Requirements
- Go 1.20+ (or version in `go.mod`)
- PostgreSQL 13+
- Redis 6+
- (Optional) Docker + Docker Compose

---

## Configuration

### App config (`configs/config.yml`)
```yaml
port: 8080

db:
  host: "localhost"   # use "db" when running under Docker Compose
  port: "5432"
  username: "postgres"
  dbName: "postgres"
  sslMode: "disable"
```

### Environment variables (`.env` at project root)
```
DB_PASSWORD=12345678
```

### Redis client (`pkg/cache/redis.go`)
```go
Addr:     "localhost:6379", // use "cache:6379" when using Docker Compose
Password: "CHANGE_ME",
```
> Make sure the Redis password matches your Docker Compose settings, if any.

---

## Database Migrations (Plain SQL)

Apply the schema once before running the API.

**psql (local):**
```bash
psql -h localhost -U postgres -d postgres -f migrations/0001_init.up.sql
```

**Docker (Postgres service name `db`):**
```bash
docker compose cp migrations/0001_init.up.sql db:/init.sql
docker compose exec db psql -U postgres -d postgres -f /init.sql
```

Rollback:
```bash
psql -h localhost -U postgres -d postgres -f migrations/0001_init.down.sql
```

---

## Run

### Local (no Docker)
```bash
go mod tidy
go run ./cmd
# API at http://localhost:8080
```

### Docker Compose
Ensure the app uses container hostnames:

- In `configs/config.yml`: `db.host: "db"`
- In `pkg/cache/redis.go`: `Addr: "cache:6379"`

Then:
```bash
docker compose up --build
# or: docker compose up --build app
```

---

## API Reference

### Auth
- **POST** `/auth/sign-up`
  ```json
  { "name": "Alice", "username": "alice", "password": "secret" }
  ```
- **POST** `/auth/sign-in`
  ```json
  { "username": "alice", "password": "secret" }
  ```
  Response:
  ```json
  { "token": "JWT_TOKEN_HERE" }
  ```
Use the token for protected routes:
```
Authorization: Bearer <JWT_TOKEN_HERE>
```

### Lists
- **POST** `/api/lists/` — create
- **GET** `/api/lists/` — list all
- **GET** `/api/lists/:id` — get by id
- **PUT** `/api/lists/:id` — update
  ```json
  { "title": "New title", "description": "New desc" }
  ```
- **DELETE** `/api/lists/:id` — delete

### Items (under a list)
- **POST** `/api/lists/:id/items/` — create item in list
  ```json
  { "title": "Milk", "description": "2L", "done": false }
  ```
- **GET**  `/api/lists/:id/items/` — all items in a list

### Items (by id)
- **GET** `/api/items/:id`
- **PUT** `/api/items/:id`
  ```json
  { "title": "Milk", "description": "1L", "done": true }
  ```
- **DELETE** `/api/items/:id`

---

## cURL Examples

**Sign up**
```bash
curl -X POST http://localhost:8080/auth/sign-up   -H "Content-Type: application/json"   -d '{"name":"Alice","username":"alice","password":"secret"}'
```

**Sign in**
```bash
curl -X POST http://localhost:8080/auth/sign-in   -H "Content-Type: application/json"   -d '{"username":"alice","password":"secret"}'
```

**Create list (authorized)**
```bash
TOKEN=<paste JWT token>
curl -X POST http://localhost:8080/api/lists/   -H "Authorization: Bearer $TOKEN"   -H "Content-Type: application/json"   -d '{"title":"Groceries","description":"Weekend"}'
```

---

## Troubleshooting

- **`dial tcp [::1]:5432: connect: connection refused`**  
  The app is trying to reach Postgres on localhost from inside Docker.  
  **Fix:** set `db.host: "db"` in `configs/config.yml` when using Docker.

- **`dial tcp: lookup redis on 127.0.0.11:53: no such host`**  
  The code is using a different hostname than your Redis service.  
  **Fix:** set `Addr: "cache:6379"` when using Docker Compose.

- **`NOAUTH Authentication required.`**  
  Redis requires a password you didn’t supply.  
  **Fix:** set the same password in Docker Compose and in `pkg/cache/redis.go`.

- **`pq: relation "users" does not exist`**  
  Tables not created yet.  
  **Fix:** run the migration: `migrations/0001_init.up.sql`.

---

## Project Structure
```
.
├─ cmd/                    # main entry
├─ configs/                # app config (viper)
├─ migrations/             # SQL schema (up/down)
├─ pkg/
│  ├─ cache/               # Redis client
│  ├─ handler/             # HTTP handlers (Gin)
│  ├─ repository/          # database/sql repositories
│  └─ service/             # business logic (JWT, hashing)
├─ todo.go                 # domain models
├─ user.go                 # domain models
└─ docker-compose.yml      # optional
```

---

## Notes
- JWT signing key is hardcoded for dev in `pkg/service/auth.go` — change it for production.
- Migrations are plain SQL files; apply them to any fresh database.


