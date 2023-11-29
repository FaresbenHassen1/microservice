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

grpc:
	go run main.go grpc

rest:
	go run main.go rest

protoc:
	protoc --proto_path=proto --go_out=proto --go_opt=paths=source_relative     --go-grpc_out=proto
 --go-grpc_opt=paths=source_relative     proto/*.proto

.PHONY: postgres createdb dropdb migrateup