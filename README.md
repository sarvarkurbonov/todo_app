# **Todo REST API — GitHub Codespaces Ready**

A clean, minimal REST API built with **Go (Gin)**, **PostgreSQL**, **Redis**, and **JWT auth** — configured to run smoothly in **GitHub Codespaces**.

---

## **Features**
- **JWT authentication**: sign-up / sign-in
- **Todo lists & items**: CRUD endpoints
- **Redis caching** for item lookups
- **Plain SQL migrations** (`migrations/`) — no ORM required

---

## **Run in GitHub Codespaces (Recommended)**

### **1) Open in Codespaces**
- Click **Code → Codespaces → Create codespace on main**.

### **2) Set database password**
Create a **`.env`** file in the project root:
```
DB_PASSWORD=12345678
```

### **3) Start Postgres + Redis**
If your compose file is named **`docker-compose.yml`**:
```bash
docker compose up --build -d db cache
```
If it’s **`Docker-compose.yml`** (capital D):
```bash
docker compose -f Docker-compose.yml up --build -d db cache
```

### **4) Apply the database migration (one-time)**
```bash
docker compose cp migrations/0001_init.up.sql db:/init.sql
docker compose exec db psql -U postgres -d postgres -f /init.sql
```
> Re-run this step if you reset or recreate the Postgres volume.

### **5) Configure the app for containers**
Edit **`configs/config.yml`**:
```yaml
port: 8080

db:
  host: "db"      # IMPORTANT inside Docker/Codespaces
  port: "5432"
  username: "postgres"
  dbName: "postgres"
  sslMode: "disable"
```

If Redis runs via Compose with service name **`cache`**, edit **`pkg/cache/redis.go`**:
```go
Addr:     "cache:6379",
Password: "CHANGE_ME", // match your compose password if set
```

### **6) Start the API**
With **`docker-compose.yml`**:
```bash
docker compose up --build app
```
With **`Docker-compose.yml`**:
```bash
docker compose -f Docker-compose.yml up --build app
```
Codespaces will **auto-forward port 8080** → open it from the **PORTS** tab.

---

## **Alternative: Run without Docker (inside Codespaces)**
1) Ensure Postgres + Redis are available in the environment.  
2) Use `db.host: "localhost"` in `configs/config.yml` and `Addr: "localhost:6379"` in `pkg/cache/redis.go`.  
3) Apply migration:
```bash
psql -h localhost -U postgres -d postgres -f migrations/0001_init.up.sql
```
4) Run:
```bash
go mod tidy
go run ./cmd
```

---

## **API Summary**

### **Auth**
- **POST** `/auth/sign-up`  
  ```json
  { "name": "Alice", "username": "alice", "password": "secret" }
  ```
- **POST** `/auth/sign-in` → returns `{ "token": "JWT..." }`

Use on protected routes:
```
Authorization: Bearer <JWT_TOKEN>
```

### **Lists**
- **POST** `/api/lists/` — create  
- **GET** `/api/lists/` — list all  
- **GET** `/api/lists/:id` — get by id  
- **PUT** `/api/lists/:id` — update  
  ```json
  { "title": "New title", "description": "New desc" }
  ```
- **DELETE** `/api/lists/:id` — delete  

### **Items (under a list)**
- **POST** `/api/lists/:id/items/` — create item in list  
  ```json
  { "title": "Milk", "description": "2L", "done": false }
  ```
- **GET**  `/api/lists/:id/items/` — all items in a list  

### **Items (by id)**
- **GET** `/api/items/:id`  
- **PUT** `/api/items/:id`  
  ```json
  { "title": "Milk", "description": "1L", "done": true }
  ```
- **DELETE** `/api/items/:id`

---

## **cURL Examples**

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

## **Troubleshooting (Common in Codespaces)**

- **`dial tcp [::1]:5432: connect: connection refused`**  
  The app is trying to reach Postgres on localhost inside Docker.  
  **Fix:** set `db.host: "db"` in `configs/config.yml`.

- **`dial tcp: lookup redis on 127.0.0.11:53: no such host`**  
  The code tries `redis`, but your service is named **`cache`**.  
  **Fix:** set Redis address to `"cache:6379"` in `pkg/cache/redis.go`.

- **`NOAUTH Authentication required.`**  
  Redis requires a password you didn’t supply.  
  **Fix:** set the same password in Compose **and** in `pkg/cache/redis.go`.

- **`pq: relation "users" does not exist`**  
  Tables not created yet.  
  **Fix:** re-run the migration step in this README.

---

## **Project Structure**
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
└─ docker-compose.yml      # or Docker-compose.yml
```

---

## **Notes**
- **JWT signing key** is hardcoded for dev in `pkg/service/auth.go` → **change for production**.
- **Migrations** are plain SQL — apply them for every fresh database volume.
- In Codespaces, use the **PORTS** tab to open **http://localhost:8080**.

**License:** MIT (or your choice)
