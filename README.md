# ToDo API

A backend microservice-based **ToDo API** built with Go, exposing both **gRPC** and **REST (via gRPC-Gateway)** endpoints. The system supports user authentication with PASETO tokens, task management, and is containerized with Docker.

---

## Architecture

The project follows **Clean Architecture** principles with a clear separation between layers:

![CleanArchitecture](docs/resource/CleanArchitecture.jpg)

```
.
├── cmd/                    # CLI entrypoints (cobra commands: gateway, user)
├── internal/               # Business logic (Clean Architecture layers)
│   ├── user/               # User & Auth domain
│   │   ├── models/         # Domain models
│   │   ├── repositories/   # Data access layer (Postgres)
│   │   ├── services/       # Use-case / service layer
│   │   ├── golibs/         # Shared libs for this domain
│   │   └── constants/      # Domain constants
│   └── taks/               # Task domain
│       ├── models/
│       ├── repositories/
│       └── services/
├── idl/pb/                 # Generated gRPC + gRPC-Gateway protobuf stubs
├── proto/                  # Protobuf source definitions
├── database/
│   └── postgres/
│       ├── migrations/     # SQL migration files (users, tasks)
│       └── queries/        # SQLC query definitions
├── transform/              # Protobuf <-> domain model converters
├── transformhelpers/       # Helper functions for transforms
├── utils/                  # Shared utilities (password hashing, auth, etc.)
├── pkg/                    # Reusable packages
├── worker/                 # Background task workers (Asynq)
├── configs/                # App configuration loading (TOML)
├── developments/           # Docker Compose & codegen scripts
├── docs/                   # Documentation assets
├── env.toml                # Local environment configuration
├── Makefile                # Developer shortcuts
└── main.go                 # Application entrypoint
```

---

## Tech Stack

### Backend
| Library | Purpose |
|---|---|
| [spf13/cobra](https://github.com/spf13/cobra) | CLI command framework |
| [grpc-ecosystem/grpc-gateway/v2](https://github.com/grpc-ecosystem/grpc-gateway) | REST gateway over gRPC |
| [google.golang.org/grpc](https://pkg.go.dev/google.golang.org/grpc) | gRPC server & client |
| [jackc/pgx/v5](https://github.com/jackc/pgx) | PostgreSQL driver |
| [kyleconroy/sqlc](https://github.com/sqlc-dev/sqlc) | Type-safe SQL codegen |
| [golang-migrate/migrate/v4](https://github.com/golang-migrate/migrate) | Database migrations |
| [hashicorp/golang-lru/v2](https://github.com/hashicorp/golang-lru) | In-memory LRU cache |
| [hibiken/asynq](https://github.com/hibiken/asynq) | Async task queue (Redis-backed) |
| [o1egl/paseto](https://github.com/o1egl/paseto) | PASETO token authentication |
| [rs/zerolog](https://github.com/rs/zerolog) | Structured logging |
| [spf13/viper](https://github.com/spf13/viper) | Configuration management (TOML) |
| [google/uuid](https://github.com/google/uuid) | UUID generation |
| [elastic/go-elasticsearch/v8](https://github.com/elastic/go-elasticsearch) | Elasticsearch client |

### Infrastructure
| Component | Role |
|---|---|
| **PostgreSQL** | Primary database |
| **Redis** | Asynq task queue backend |
| **gRPC** | Internal service communication |
| **Docker / Docker Compose** | Containerized development environment |
| **Adminer** | Database UI (port `3037`) |

---

## Services

The application exposes two independently runnable services:

| Service | Default Port | Protocol | Description |
|---|---|---|---|
| `user` | `3030` | gRPC | User management & authentication |
| `gateway` | `3031` | HTTP/REST | gRPC-Gateway reverse proxy |

---

## Getting Started

### Prerequisites
- [Go 1.24+](https://golang.org/)
- [Docker & Docker Compose](https://docs.docker.com/compose/)
- [golang-migrate CLI](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) *(optional, for direct migration)*
- [evans](https://github.com/ktr0731/evans) *(optional, for gRPC REPL)*

### 1. Configure Environment

Copy and edit the config file:
```bash
cp env.toml env.local.toml
# Edit database URL, ports, symmetric key, etc.
```

Default configuration (`env.toml`):
```toml
postgres_url = "postgres://admin:postgres@localhost:5432/postgres?sslmode=disable"
symmetric_key = "FKhIvVYqbJihFRJeMcZrpGwJrcKrgQGz"

[user_service_endpoint]
host = "localhost"
port = 3030

[gateway_service_endpoint]
host = "localhost"
port = 3031
```

### 2. Start Infrastructure

```bash
# Start PostgreSQL
make start-postgres

# Start Adminer (DB UI at http://localhost:3037)
make adminer
```

### 3. Run Database Migrations

```bash
# Run all migrations
make migrate

# Run user-specific migrations only
make user-migrate
```

### 4. Start Services

```bash
# Start the user gRPC service (port 3030)
make start-user

# Start the REST gateway service (port 3031)
make start-gateway
```

---

## Development

### Code Generation

```bash
# Regenerate protobuf stubs (requires Docker)
make gen-proto

# Regenerate SQLC type-safe queries (requires Docker)
make gen-sql

# Regenerate mocks for user repository
make gen-mock-user
```

### gRPC REPL (Evans)

```bash
make evans
# Connects to localhost:9091
```

---

## Database Schema

The database is managed with SQL migration files located in `database/postgres/migrations/`.

### `users` table
Stores registered users with hashed passwords, roles, and timestamps.

### `tasks` table
```sql
CREATE TABLE IF NOT EXISTS tasks (
    id          text PRIMARY KEY,
    title       text,
    description text,
    status      text        NOT NULL DEFAULT 'TaskStatus_TODO',
    user_id     text        NOT NULL,
    created_at  timestamptz NOT NULL DEFAULT NOW(),
    updated_at  timestamptz NOT NULL DEFAULT NOW(),
    deleted_at  timestamptz,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

---

## Authentication

Authentication uses **PASETO v2** tokens. The flow:
1. `Register` → creates user, returns a signed PASETO token.
2. `Login` → validates credentials + role, returns a PASETO token.
3. All protected endpoints require the token in the `Authorization` header.

---

## Makefile Reference

| Command | Description |
|---|---|
| `make start-postgres` | Start PostgreSQL via Docker Compose |
| `make adminer` | Start Adminer DB UI |
| `make migrate` | Run all DB migrations |
| `make user-migrate` | Run user-specific DB migrations |
| `make start-user` | Start the gRPC user service |
| `make start-gateway` | Start the REST gateway service |
| `make gen-proto` | Regenerate protobuf stubs |
| `make gen-sql` | Regenerate SQLC queries |
| `make gen-mock-user` | Regenerate mocks for user repository |
| `make evans` | Open Evans gRPC REPL |
| `make sqlc` | Run sqlc generate (local) |