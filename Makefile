#SHELL = /bin/bash

.PHONY: dep
dep:
	go mod download
	go mod tidy

.PHONY: test
test:
	go test ./...

.PHONY: race
race:
	go test -v -race ./...

.PHONY: init
lint:
	/home/user/go/bin/golangci-lint run

accrual-port = $(shell ./.tools/random unused-port)
gophermart-bin = ./cmd/gophermart/gophermart
.PHONY: autotest
autotest:
	go build -o $(gophermart-bin) ./cmd/gophermart
	./.tools/gophermarttest \
    -test.v -test.run=^TestGophermart$ \
    -gophermart-binary-path=$(gophermart-bin) \
    -gophermart-host=localhost \
    -gophermart-port=8080 \
    -gophermart-database-uri="postgresql://postgres:postgres@localhost/praktikum?sslmode=disable" \
    -accrual-binary-path=./cmd/accrual/accrual_linux_amd64 \
    -accrual-host=localhost \
    -accrual-port=$(accrual-port) \
    -accrual-database-uri="postgresql://postgres:postgres@localhost/praktikum?sslmode=disable"
	rm $(gophermart-bin)

.PHONY: gen
gen:
	go generate ./...

.PHONY: cover
cover:
	go test -short -count=1 -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o=coverage.html
	rm coverage.out

.PHONY: statictest
statictest:
	go vet -vettool=./.tools/statictest ./...

.PHONY: dev-up
dev-up:
	docker-compose -f=docker-compose.dev.yml --env-file=.env.dev up -d

.PHONY: dev-down
dev-down:
	docker-compose -f=docker-compose.dev.yml --env-file=.env.dev down --rmi local

.PHONY: dev-autotest
dev-autotest:
	docker-compose -f=docker-compose.dev-test.yml --env-file=.env.dev up -d
	docker-compose -f=docker-compose.dev-test.yml --env-file=.env.dev logs -f tests
	docker-compose -f=docker-compose.dev-test.yml --env-file=.env.dev down --rmi local

.PHONY: accrual
accrual:
	./cmd/accrual/accrual_linux_amd64 -a=localhost:8787 -d="postgresql://postgres:postgres@localhost:5432/praktikum?sslmode=disable"