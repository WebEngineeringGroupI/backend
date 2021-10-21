

deps:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

migrate-db: deps
	migrate -path ./database/migrate/ -database "postgres://postgres:root@localhost:5432/postgres?sslmode=disable" up

run-db:
	docker run --rm -it -p 5432:5432 -e POSTGRES_PASSWORD=root postgres

fmt:
	find -iname '*.go' | xargs -L1 gofmt -s -w

lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	golangci-lint run --timeout 1h
