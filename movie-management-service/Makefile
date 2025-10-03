DB_USER = $(MOVIEPIN_DB_USER)
DB_NAME = $(MOVIEPIN_DB_NAME)

DATABASE_URL = postgres://$(DB_USER)@$(MOVIEPIN_DB_HOST):$(MOVIEPIN_DB_PORT)/$(DB_NAME)?sslmode=disable

MIGRATIONS_PATH = db/migrations

create-db:
	psql -c 'CREATE DATABASE "$(DB_NAME)"' -U $(DB_USER) -d postgres;
	psql -c 'CREATE EXTENSION IF NOT EXISTS "uuid-ossp"' -U $(DB_USER) -d $(DB_NAME);

drop-db:
	psql -c 'DROP DATABASE IF EXISTS "$(DB_NAME)"' -U $(DB_USER) -d postgres

migrate-up:
	migrate -database $(DATABASE_URL) -path $(MIGRATIONS_PATH) up

migrate-down:
	migrate -database $(DATABASE_URL) -path $(MIGRATIONS_PATH) down