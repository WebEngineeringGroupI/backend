
build: deps
	go build -v ./...

deps:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/onsi/ginkgo/ginkgo@v1

migrate-db: deps
	migrate -path ./database/migrate/ -database "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" up

run-db:
	docker pull postgres
	docker run --name postgres --rm -p 5432:5432 -e POSTGRES_PASSWORD=postgres postgres &

kill-db:
	docker rm -f postgres

bump:
	GOPRIVATE="github.com/WebEngineeringGroupI/*" go get -d -u -v -t ./...

fmt:
	find -iname '*.go' | xargs -L1 gofmt -s -w
	go mod tidy

lint: deps
	golangci-lint run --timeout 1h

test-unit: deps
	ginkgo -r -race -randomizeAllSpecs -randomizeSuites -trace -progress -cover -skipPackage ./pkg/infrastructure

test-integration: run-db deps
	sleep 10 # Give some time to DB to be launched
	$(MAKE) migrate-db
	ginkgo -r -race -randomizeAllSpecs -randomizeSuites -trace -progress -cover -p ./pkg/infrastructure
	$(MAKE) kill-db
