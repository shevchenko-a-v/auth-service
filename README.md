# auth-service

### Migrator cmd
go run ./cmd/migrator --storage-path=./storage/auth.db --migrations-path=./migrations --migrations-table=migrations

### Migrator cmd for tests
go run ./cmd/migrator --storage-path=./storage/auth.db --migrations-path=./tests/migrations --migrations-table=migrations_test
