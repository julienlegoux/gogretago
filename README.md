# GoGreTaGo - CovoitAPI in Go

A Go implementation of the CovoitAPI REST API using Gin framework, following Clean Architecture principles.

## Project Structure

```
gogretago/
├── cmd/api/main.go              # Entry point
├── internal/
│   ├── domain/                  # Core business logic
│   ├── application/             # Use cases and DTOs
│   ├── infrastructure/          # External implementations
│   └── presentation/            # HTTP layer
├── config/config.go             # Configuration
├── Dockerfile                   # Container image
└── docker-compose.yml           # Development environment
```

## Tech Stack

- **Framework**: Gin
- **ORM**: GORM with PostgreSQL
- **Authentication**: JWT (golang-jwt/jwt/v5)
- **Password Hashing**: Argon2id
- **Email**: Resend
- **Validation**: go-playground/validator

## Getting Started

### Prerequisites

- Go 1.25+ (or Docker)
- PostgreSQL database

### Configuration

Copy `.env.example` to `.env` and configure:

```bash
cp .env.example .env
```

### Running Locally

```bash
go run cmd/api/main.go
```

### Running with Docker

```bash
docker-compose up --build
```

## API Endpoints

| Method | Path            | Description       |
|--------|-----------------|-------------------|
| GET    | `/health`       | Health check      |
| POST   | `/auth/register`| User registration |
| POST   | `/auth/login`   | User login        |

## License

MIT
