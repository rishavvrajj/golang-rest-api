# Go REST API (Students Management System)

A clean, production-grade, and educational REST API written in Go (Golang) using standard library `net/http` (Go 1.22+ enhanced routing), SQLite, structured logging (`slog`), input validation (`go-playground/validator`), and a layered architecture.

Designed as a premier reference and learning guide for building clean, maintainable, and scalable Go REST APIs.

---

## Project Architecture (Clean Architecture Pattern)

```text
golang-rest-api/
├── cmd/
│   └── students-api/
│       └── main.go               # Application entry point, server startup & graceful shutdown
├── config/
│   └── local.yaml                # Configuration file (Env, Storage Path, Server Address)
├── internal/
│   ├── config/
│   │   └── config.go             # Configuration loader using cleanenv
│   ├── http/
│   │   ├── handlers/
│   │   │   └── student/
│   │   │       └── student.go    # HTTP handlers for student CRUD operations
│   │   ├── middleware/
│   │   │   └── middleware.go     # HTTP middleware (Structured logging & Panic recovery)
│   │   ├── response/
│   │   │   └── response.go       # Centralized JSON response & validation error helpers
│   │   └── router/
│   │       └── router.go         # HTTP router (ServeMux) & route registration
│   ├── models/
│   │   └── student.go            # Domain models and validation tags
│   ├── service/
│   │   └── student_service.go    # Business logic layer
│   └── storage/
│       ├── storage.go            # Repository interface definition
│       └── sqlite/
│           └── sqlite.go         # SQLite database connection & implementation
├── go.mod
└── go.sum
```

---

## Prerequisites

- **Go 1.22+** installed on your system.

---

## Configuration (`config/local.yaml`)

```yaml
env: "local"
storage_path: "storage/storage.db"
http_server:
  address: "localhost:8082"
```

---

## Running the Application

1. **Tidy dependencies**:
   ```bash
   go mod tidy
   ```

2. **Build the application**:
   ```bash
   go build -o bin/students-api.exe cmd/students-api/main.go
   ```

3. **Run the server**:
   ```bash
   go run cmd/students-api/main.go
   ```
   Or specify a custom config file:
   ```bash
   go run cmd/students-api/main.go --config=config/local.yaml
   ```

---

## API Endpoints & Testing (`curl` Examples)

### 1. Health Check (`GET /health`)
```bash
curl -i http://localhost:8082/health
```

### 2. Create a Student (`POST /students`)
```bash
curl -X POST http://localhost:8082/students \
  -H "Content-Type: application/json" \
  -d '{"name": "Alice Smith", "email": "alice@example.com", "age": 22}'
```

### 3. Get All Students (`GET /students`)
```bash
curl http://localhost:8082/students
```

### 4. Get Student by ID (`GET /students/{id}`)
```bash
curl http://localhost:8082/students/1
```

### 5. Update Student (`PUT /students/{id}`)
```bash
curl -X PUT http://localhost:8082/students/1 \
  -H "Content-Type: application/json" \
  -d '{"name": "Alice Johnson", "email": "alice.j@example.com", "age": 23}'
```

### 6. Delete Student (`DELETE /students/{id}`)
```bash
curl -X DELETE http://localhost:8082/students/1
```

---

## Key Learning Concepts Covered

1. **Layered Architecture**: Clear separation of concerns (`cmd` -> `router` -> `handlers` -> `service` -> `storage/repository` -> `models`).
2. **Configuration Management**: Using `cleanenv` to load configuration from YAML files and environment variables.
3. **Structured Logging**: Using Go's built-in `log/slog` package for key-value structured logging.
4. **Modern Go Routing**: Using Go 1.22+ standard library `http.ServeMux` with HTTP method matching and path wildcards (`/students/{id}`).
5. **Middleware Pattern**: Custom middleware for HTTP request logging and panic recovery.
6. **Database Persistence**: Using `database/sql` with SQLite (`mattn/go-sqlite3`) for persistent CRUD operations.
7. **Input Validation & Error Handling**: Validating request payloads using `go-playground/validator` with centralized response formatters.
8. **Graceful Shutdown**: Handling OS signals (`SIGINT`, `SIGTERM`) to cleanly shut down the HTTP server without dropping active connections.
