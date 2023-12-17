
DB_URL=postgresql://postgres:root@localhost:5432/interactive_presentations?sslmode=disable

name = init_schema

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

build:
	go build -o server server.go

run:
	./server
	
sqlc: 
	sqlc generate

test:
	go test -timeout 3m -v -cover -short ./...

bench_test:
	go test -bench=. -count 10 -run=^# ./...

