POSTGRES_USER=postgres
HOST=db
POST=5432
POSTGRES_DB=postgres
OPTIONS=sslmode=disable

POSTGRES_URI=postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(HOST):$(PORT)/$(POSTGRES_DB)?$(OPTIONS)

# docker commands
run: up

stop: down

up:
	docker compose up --force-recreate --build

down:
	docker compose down --rmi local

clean:
	docker image rm crud-app_backend

# migrate commands
migrate-create:
	docker compose --profile migrate run --rm create-migration $(SEQ)

migrate-drop:
	docker compose --profile migrate run --rm migrate-drop

migrate-up:
	docker compose --profile migrate run --rm migrate

swag-init:
	docker compose --profile swag run --rm swag-init   