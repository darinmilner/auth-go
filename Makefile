.PHONY: migrate

migrate:
	migrate -path db/migrations -database "postgresql://root:password@localhost:5432/auth?sslmode=disable" -verbose up