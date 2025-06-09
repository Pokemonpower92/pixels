# Pixels

This is the backend for the pixels application.
It takes any png and creates a pixelated replica of it.

Requires go 1.22^

## QuickStart

1. Install [sql-migrate](https://github.com/rubenv/sql-migrate):

   ```
   go install github.com/rubenv/sql-migrate/...@latest
   ```

2. Install [sqlc](https://github.com/sqlc-dev/sqlc):

   ```
   go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
   ```

3. Generate sqlc files:

   ```
   sqlc generate -f ./internal/sqlc/sqlc.yml
   ```

4. Run the tests:

   ```
   make test
   ```

5. Start the local server(ensure docker is installed and running):

   ```
   docker compose up
   ```
