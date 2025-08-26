Todo REST API (Go + Gin + Postgres + Redis)

A minimal todo API with JWT auth, PostgreSQL storage, and Redis caching.

Tech Stack

Go (Gin)

PostgreSQL

Redis

Migrations: plain SQL (in migrations/)

Quick Start
1) Configure database

Edit configs/config.yml (defaults shown):

port: 8080

db:
  host: "localhost"    # set to "db" if using Docker
  port: "5432"
  username: "postgres"
  dbName: "postgres"
  sslMode: "disable"


Set DB password via .env (used at runtime):

DB_PASSWORD=12345678

2) Apply migrations (one time)

Run the SQL file to create tables:

Using psql:

psql -h <db_host> -U postgres -d postgres -f migrations/0001_init.up.sql


If you’re running Docker Compose and the Postgres service is named db:

docker exec -i <db_container_name> psql -U postgres -d postgres < migrations/0001_init.up.sql


If you ever need to drop everything:

psql -h <db_host> -U postgres -d postgres -f migrations/0001_init.down.sql

3) Redis config

Edit pkg/cache/redis.go.

For local Redis: keep Addr: "localhost:6379", set Password.

For Docker Compose (service name cache), change to:

Addr:     "cache:6379",
Password: "CHANGE_ME",

4) Run

Local (no Docker):

go run ./cmd


Docker Compose:

docker compose up --build
# App → http://localhost:8080

API
Auth

POST /auth/sign-up
Body:

{ "name": "Alice", "username": "alice", "password": "secret" }


POST /auth/sign-in
Body:

{ "username": "alice", "password": "secret" }


Response:

{ "token": "JWT_TOKEN_HERE" }


Use the token for protected routes:
Authorization: Bearer <JWT_TOKEN_HERE>

Lists

POST /api/lists/ – create list
Body:

{ "title": "Groceries", "description": "Weekend shopping" }


GET /api/lists/ – list all

GET /api/lists/:id – get by id

PUT /api/lists/:id – update
Body (any field):

{ "title": "New Title", "description": "New Desc" }


DELETE /api/lists/:id

Items (under a list)

POST /api/lists/:id/items/ – add item to list
Body:

{ "title": "Milk", "description": "2L", "done": false }


GET /api/lists/:id/items/ – get all items for a list

Items (by id)

GET /api/items/:id

PUT /api/items/:id
Body (any field):

{ "title": "Milk", "description": "1L", "done": true }


DELETE /api/items/:id

Environment & Secrets

DB password: from .env → DB_PASSWORD.

JWT key (dev only): hardcoded in pkg/service/auth.go (signingKey = "asd234asd").
Change it for production.

Troubleshooting

dial tcp [::1]:5432: connect: connection refused
The app is trying to reach Postgres on localhost from inside Docker.
Use the service name (db) in configs/config.yml when using Compose.

dial tcp [::1]:6379: connect: connection refused or no such host redis
Update Redis address:

Local: localhost:6379

Docker Compose: cache:6379

NOAUTH Authentication required.
Set the correct Redis password in pkg/cache/redis.go.

pq: relation "users" does not exist
Run the migration: migrations/0001_init.up.sql.

Project Structure (short)
.
├─ cmd/                 # main entry
├─ configs/             # app config (viper)
├─ migrations/          # SQL schema (up/down)
├─ pkg/
│  ├─ cache/            # Redis client
│  ├─ handler/          # HTTP handlers (Gin)
│  ├─ repository/       # database/sql repositories
│  └─ service/          # business logic (JWT, hashing, etc.)
├─ todo.go              # Todo domain models
├─ user.go              # User domain model
└─ docker-compose.yml   # optional: run db/cache/app together

Notes

Migrations are plain SQL and must be applied before using the API.

Redis is used to cache item lookups; adjust address/password to your setup.

JWT TTL is 1 hour by default (see pkg/service/auth.go).
