migrate-up:
	goose -dir pkg/db/migrations postgres "host=$(PSQL_HOST) port=$(PSQL_PORT) user=$(PSQL_USER) password=$(PSQL_PASSWORD) dbname=$(PSQL_DBNAME) sslmode=$(PSQL_SSLMODE)" up

migrate-down:
	goose -dir pkg/db/migrations postgres "host=$(PSQL_HOST) port=$(PSQL_PORT) user=$(PSQL_USER) password=$(PSQL_PASSWORD) dbname=$(PSQL_DBNAME) sslmode=$(PSQL_SSLMODE)" down



