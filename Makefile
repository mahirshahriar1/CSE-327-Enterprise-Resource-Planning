run:
	@go run main.go

migrate:
	@psql -U postgres -d erp -f db/migration.sql
