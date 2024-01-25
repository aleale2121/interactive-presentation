
DB_URL=postgresql://postgres:root@localhost:5432/interactive_presentations?sslmode=disable
BASE_URL=http://localhost:8080
name = init_schema

PRESENTATION_ID =  '461a92d5-ae42-429c-b11f-997b91e197b6'
POLL_ID='906c0b4b-b64d-49c4-81e5-ea274ba534d9'
KEY='A'

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

slide_load_test:
	hey -m GET -c 100 -n 10000 '$(BASE_URL)/presentations/$(PRESENTATION_ID)/polls/current'

vote_load_test:
	hey -n 10000 -c 100 -H "Content-Type: application/json" -m POST -d '{"poll_id":"'"$(POLL_ID)"'","key":"A","client_id":"'"$(shell uuidgen)"'"}' '$(BASE_URL)/presentations/$(PRESENTATION_ID)/polls/current/votes'

clean_cache:
	go clean -testcache