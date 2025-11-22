# Petclinic

Small Go project demonstrating a PostgreSQL-backed service for a pet clinic.

## Overview
This repository contains a simple Go service that connects to PostgreSQL using environment configuration. The DB initialization lives in `db/db.go`.

## Prerequisites
- Go 1.18+
- PostgreSQL server
- Git

## Environment
Create a `.env` file at the project root with the following variable:

POSTGRESQL should be a valid Postgres connection string, for example:

POSTGRESQL="postgres://username:password@localhost:5432/dbname?sslmode=disable"

Example `.env`:
```
POSTGRESQL="postgres://user:pass@localhost:5432/petclinic?sslmode=disable"
```

## Setup & Run
1. Install dependencies:
   - The project uses `github.com/joho/godotenv` and `github.com/lib/pq`. Run:
   ```
   go mod tidy
   ```
2. Ensure your `.env` file is present and correct.
3. Initialize/verify DB connection (the project calls `db.InitDB()` in startup).
4. Build and run:
   ```
   go build ./...
   ./your-binary-name
   ```

## Notes
- `db/db.go` reads the `POSTGRESQL` env var using `godotenv`. Make sure `.env` is available if running locally.
- If you see `POSTGRESQL environment variable not set`, confirm the `.env` file path and variable name.

## Troubleshooting
- Connection errors: verify host, port, credentials, and `sslmode`.
- Missing packages: run `go mod tidy` to fetch required modules.

## License
Unlicensed — adapt as needed.
