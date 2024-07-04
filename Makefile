startPG:
	docker stop simple_bank_pg || true
	docker run --name simple_bank_pg -e POSTGRES_USER=root -e POSTGRES_PASSWORD=admin -p 5432:5432 -d postgres:16-alpine

createdb:
	docker exec -it simple_bank_pg createdb --username=root --owner=root simple_bank

removedb:
	docker exec -it simple_bank_pg dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:admin@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:admin@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

generateStoreMock:
	mockgen -destination db/sqlc/mock/store.go -package mockdb -build_flags=--mod=mod github.com/xianlinbox/simple_bank/db/sqlc Store

test:
	go test -v --cover ./...

start-server:
	go run main.go


.PHONY: startPG createdb removedb