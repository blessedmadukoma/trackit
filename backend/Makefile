postgres:
		docker run --name postgres14 -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -d postgres:14-alpine

createdb:
		docker exec -it postgres14 createdb --username=postgres --owner=postgres trackit

dropdb:
		docker exec -it postgres14 dropdb --username=postgres trackit

psql: # login to trackit db in psql terminal
		docker exec -it postgres14 psql -U postgres -d trackit

createmigration:
		migrate -help
		migrate create -ext sql -dir db/migration -seq [ADDNAME]

migrateup:
		migrate -path db/migration -database "postgresql://postgres:postgres@localhost:5432/trackit?sslmode=disable" -verbose up

migratedown:
		migrate -path db/migration -database "postgresql://postgres:postgres@localhost:5432/trackit?sslmode=disable" -verbose down

sqlc:
		sqlc generate

test:
		go test -v -cover ./...

.PHONY: postgres createdb dropdb psql createmigration migrateup migratedown sqlc test