postgres: 
	docker run --name postg -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -p 5432\:5432 -v pgdata\:/var/lib/postgresql/data -d postgres

createdb:
	docker exec -it postg createdb --username=postgres --owner=postgres projecttest

dropdb:
	docker exec -it postg dropdb --username=postgres --owner=postgres projecttest

migrateup:
	migrate -path db/migration -database postgresql://postgres:postgres@localhost:5432/projecttest?sslmode=disable -verbose up

migratedown:
	migrate -path db/migration -database postgresql://postgres:postgres@localhost:5432/projecttest?sslmode=disable -verbose down

servergrpc:
	go run grpc_server/main.go

clientgrpc:
	go run grpc_client/main.go

.PHONY: postgres createdb dropdb migrateup