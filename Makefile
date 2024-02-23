USER = ranmerc
DB_NAME = moviepin
DATABASE_URL = postgres://$(ranmerc)@localhost:5432/$(DB_NAME)?sslmode=disable
MIGRATIONS_PATH = db/migrations

create-db:
	psql -c 'CREATE DATABASE "$(DB_NAME)"' -U $(USER) -d postgres;
	psql -c 'CREATE EXTENSION IF NOT EXISTS "uuid-ossp"' -U $(USER) -d $(DB_NAME);

drop-db:
	psql -c 'DROP DATABASE IF EXISTS "$(DB_NAME)"' -U $(USER) -d postgres

migrate-up:
	migrate -database $(DATABASE_URL) -path $(MIGRATIONS_PATH) up

migrate-down:
	migrate -database $(DATABASE_URL) -path $(MIGRATIONS_PATH) down