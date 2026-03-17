# TaskFlow API

> A production-grade backend platform for task management, built with Go.
> Exposes both **gRPC** (internal) and **REST/HTTP** (external via gRPC-Gateway) endpoints.

[![CI](https://github.com/glucozo192/glu-project/actions/workflows/ci.yaml/badge.svg)](https://github.com/glucozo192/glu-project/actions/workflows/ci.yaml)
![Go Version](https://img.shields.io/badge/Go-1.24-00ADD8?logo=go)
![License](https://img.shields.io/badge/license-MIT-green)

---

## Key Features

- 🔐 **PASETO Authentication** — Stateless token auth with replay-attack protection (token stored & verified against DB on each request)
- 🛡️ **Role-Based Access Control (RBAC)** — Per-route permission checks using role → permission mappings, with an LRU cache to minimize DB hits
- 🔍 **Full-text Search with Elasticsearch** — Tasks are indexed on creation and searchable via `multi_match` across `title` and `description`
- ⚡ **Async Task Queue** — Background jobs powered by [Asynq](https://github.com/hibiken/asynq) (Redis-backed)
- 🏗️ **Clean Architecture** — Strict separation of `models`, `repositories`, `services`, and `transport` layers
- 🐳 **Docker-first** — Full Docker Compose setup: Postgres, migrations, proto codegen, and SQLC codegen

---

## Architecture

The project follows **Clean Architecture** with clear separation across layers:

```
Transport Layer   cmd/gateway.go (REST) · cmd/user.go (gRPC)
        ↓
Service Layer     internal/user/services · internal/taks/services
        ↓
Repository Layer  internal/*/repositories/postgres · /elasticsearch
        ↓
Data Layer        PostgreSQL (pgx/v5 + SQLC) · Elasticsearch
```

![Clean Architecture Diagram](docs/resource/CleanArchitecture.jpg)

### Project Structure

```
.
├── cmd/                    # CLI entry points (cobra): gateway, user
├── internal/
│   ├── user/               # User & Auth domain
│   │   ├── models/         # Domain types
│   │   ├── repositories/   # Postgres & login repos
│   │   └── services/       # auth.go (Login, Register)
│   └── taks/               # Task domain
│       ├── repositories/
│       │   ├── postgres/   # SQLC-generated queries
│       │   └── elasticsearch/  # Index & full-text search
│       └── services/       # Task CRUD + ES indexing
├── idl/pb/                 # Generated gRPC + gRPC-Gateway stubs
├── proto/                  # Protobuf source definitions
├── pkg/
│   ├── grpc_server/        # gRPC server wrapper
│   ├── http_server/        # HTTP server + auth middleware
│   └── elasticsearch_client/
├── database/postgres/
│   ├── migrations/         # SQL migrations (users, tasks)
│   └── queries/            # SQLC query definitions
├── utils/authenticate/     # PASETO token generation & verification
├── transform/              # Protobuf ↔ domain model converters
├── worker/                 # Asynq background task workers
├── developments/           # Docker Compose, codegen scripts
├── env.example.toml        # Configuration template (copy → env.toml)
└── Makefile
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
| [golang-migrate/migrate](https://github.com/golang-migrate/migrate) | Database schema migrations |
| [hashicorp/golang-lru/v2](https://github.com/hashicorp/golang-lru) | Expirable LRU cache for RBAC permissions |
| [hibiken/asynq](https://github.com/hibiken/asynq) | Async task queue (Redis-backed) |
| [o1egl/paseto](https://github.com/o1egl/paseto) | PASETO v2 token auth |
| [elastic/go-elasticsearch/v8](https://github.com/elastic/go-elasticsearch) | Elasticsearch full-text search |
| [rs/zerolog](https://github.com/rs/zerolog) | Structured JSON logging |
| [spf13/viper](https://github.com/spf13/viper) | TOML configuration |

### Infrastructure
| Component | Role |
|---|---|
| **PostgreSQL 13** | Primary relational database |
| **Elasticsearch 8** | Full-text search engine for tasks |
| **Redis** | Asynq background job queue |
| **Docker / Docker Compose** | Containerized dev environment |
| **Adminer** | Database web UI (`localhost:3037`) |

---

## Services

| Service | Port | Protocol | Description |
|---|---|---|---|
| `user` | `3030` | gRPC | User management & PASETO authentication |
| `gateway` | `3031` | HTTP/REST | gRPC-Gateway reverse proxy with RBAC middleware |

---

## Getting Started

### Prerequisites

- [Go 1.24+](https://golang.org/)
- [Docker & Docker Compose](https://docs.docker.com/compose/)

### 1. Configure Environment

```bash
cp env.example.toml env.toml
# Edit env.toml: set postgres_url, symmetric_key (32 chars), admin credentials
```

### 2. Start Infrastructure

```bash
# PostgreSQL
make start-postgres

# (Optional) Adminer DB UI → http://localhost:3037
make adminer
```

### 3. Run Migrations

```bash
make migrate        # All migrations
make user-migrate   # User-specific migrations only
```

### 4. Start Services

```bash
# Terminal 1 — gRPC user service (port 3030)
make start-user

# Terminal 2 — REST gateway (port 3031)
make start-gateway
```

---

## Authentication Flow

```
1. POST /register → creates user → returns PASETO token
2. POST /login    → verifies credentials + role → returns PASETO token
3. All protected routes → middleware validates token against DB (replay protection)
                       → checks role permissions via RBAC (LRU cached)
```

---

## Security

| Mechanism | Implementation |
|---|---|
| Token type | PASETO v2 (local) — more secure than JWT |
| Replay protection | Token stored in `users.token` column; verified on every request |
| RBAC | Role → permissions mapping; per-route enforcement in HTTP middleware |
| Permission cache | `expirable.LRU` invalidated on role changes |

---

## Development

### Code Generation

```bash
make gen-proto      # Regenerate protobuf stubs (Docker)
make gen-sql        # Regenerate SQLC queries (Docker)
make gen-mock-user  # Regenerate mocks for unit tests
```

### Tests

```bash
go test ./utils/... -v -race
```

### gRPC REPL

```bash
make evans   # Evans → localhost:9091
```

---

## Database Schema

Migrations are in `database/postgres/migrations/`.

**`tasks`**
```sql
CREATE TABLE IF NOT EXISTS tasks (
    id          text PRIMARY KEY,
    title       text,
    description text,
    status      text        NOT NULL DEFAULT 'TaskStatus_TODO',
    user_id     text        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at  timestamptz NOT NULL DEFAULT NOW(),
    updated_at  timestamptz NOT NULL DEFAULT NOW(),
    deleted_at  timestamptz
);
```

Tasks are also indexed into Elasticsearch on write for full-text search.

---

## Makefile Reference

| Command | Description |
|---|---|
| `make start-postgres` | Start PostgreSQL via Docker Compose |
| `make adminer` | Start Adminer DB UI at port 3037 |
| `make migrate` | Run all DB migrations |
| `make user-migrate` | Run user-specific migrations |
| `make start-user` | Start the gRPC user service |
| `make start-gateway` | Start the REST gateway |
| `make gen-proto` | Regenerate protobuf stubs |
| `make gen-sql` | Regenerate SQLC type-safe queries |
| `make gen-mock-user` | Regenerate mocks for user repository |
| `make evans` | Open Evans gRPC REPL |