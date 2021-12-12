
build: deps
	go build -v ./...

generate: deps clean
	go generate -x ./...

clean:
	-find -name "mocks" -type d | xargs rm -rf
	-find -name "*.coverprofile" -type f | xargs rm

deps:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/onsi/ginkgo/ginkgo@v1
	go install github.com/golang/mock/mockgen@v1.6.0

migrate-db: deps
	migrate -path ./database/migrate/ -database "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" up

clean-db: deps
	migrate -path ./database/migrate/ -database "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" down -all

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

test-unit: generate
	ginkgo -r -race -randomizeAllSpecs -randomizeSuites -trace -progress -cover -skipPackage ./pkg/infrastructure

test-integration: run-db generate
	sleep 10 # Give some time to DB to be launched
	$(MAKE) clean-db
	$(MAKE) migrate-db
	ginkgo -r -race -randomizeAllSpecs -randomizeSuites -trace -progress -cover -p ./pkg/infrastructure
	$(MAKE) kill-db
