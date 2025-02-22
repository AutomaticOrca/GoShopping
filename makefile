ENV_FILE = app.env

include $(ENV_FILE)


postgres:
	@docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=$(DB_USERNAME) -e POSTGRES_PASSWORD=$(DB_PASSWORD) -d postgres:12-alpine

createdb:
	@docker exec -it postgres12 createdb --username=$(DB_USERNAME) --owner=root $(DB_NAME)

dropdb:
	@docker exec -it postgres12 dropdb $(DB_NAME)

migrateup:
	@migrate -path $(CURDIR)/internal/database/migration -database "postgresql://$(DB_USERNAME):$(DB_PASSWORD)@localhost:5432/$(DB_NAME)?sslmode=disable" -verbose up

migratedown:
	@migrate -path $(CURDIR)/internal/database/migration -database "postgresql://$(DB_USERNAME):$(DB_PASSWORD)@localhost:5432/$(DB_NAME)?sslmode=disable" -verbose down

.PHONY: postgres createdb dropdb migrateup migratedown