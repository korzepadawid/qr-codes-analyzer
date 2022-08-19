run:
	go run main.go

db:
	docker run --name postgres-qr-codes -p5432:5432 -d -e POSTGRES_PASSWORD=postgres -e POSTGRES_USER=postgres -e POSTGRES_DB=qr-codes postgres:14.4-alpine

migrate-up:
	migrate -path db/migrations/ -database postgresql://postgres:postgres@localhost:5432/qr-codes?sslmode=disable -verbose up

migrate-down:
	migrate -path db/migrations/ -database postgresql://postgres:postgres@localhost:5432/qr-codes?sslmode=disable -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

mock-store:
	mockgen --build_flags=--mod=mod -package mockdb -destination db/mock/store.go github.com/korzepadawid/qr-codes-analyzer/db/sqlc Store

mock-token:
	mockgen --build_flags=--mod=mod -package mocktoken -destination token/mock/token.go github.com/korzepadawid/qr-codes-analyzer/token Provider

mock-password:
	mockgen --build_flags=--mod=mod -package mockpassword -destination util/mock/password.go github.com/korzepadawid/qr-codes-analyzer/util PasswordService