

deps:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

migrate-db: deps
	migrate -path ./database/migrate/ -database "postgres://postgres:root@localhost:5432/postgres?sslmode=disable" up

run-db:
	docker run --rm -it -p 5432:5432 -e POSTGRES_PASSWORD=root postgres
