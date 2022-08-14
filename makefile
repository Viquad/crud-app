POSTGRES_USER=postgres
HOST=db
POST=5432
POSTGRES_DB=postgres
OPTIONS=sslmode=disable

POSTGRES_URI=postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(HOST):$(PORT)/$(POSTGRES_DB)?$(OPTIONS)

run:
	docker compose up

stop:
	docker compose down
	
clean:
	docker image rm crud-app_backend

docker: wait-postgres migrate
	./main

wait-postgres:
	./wait-for-postgres.sh db

migrate:
	goose -dir ./schema -table schema_migrations postgres $(POSTGRES_URI) up-to 20220814121331