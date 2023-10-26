SHELL = /bin/bash
.PHONY: autotest dep test race lint gen cover statictest dev-up dev-down accrual

dep:
	go mod download
	go mod tidy

test:
	go test ./...

race:
	go test -v -race ./...

lint:
	/home/user/go/bin/golangci-lint run

accrual-port = $(shell ./.tools/random unused-port)
gophermart-bin = ./cmd/gophermart/gophermart
autotest:
	go build -o $(gophermart-bin) ./cmd/gophermart
	./.tools/gophermarttest \
    -test.v -test.run=^TestGophermart$ \
    -gophermart-binary-path=$(gophermart-bin) \
    -gophermart-host=localhost \
    -gophermart-port=8080 \
    -gophermart-database-uri="postgresql://postgres:postgres@localhost:5432/praktikum?sslmode=disable" \
    -accrual-binary-path=./cmd/accrual/accrual_linux_amd64 \
    -accrual-host=localhost \
    -accrual-port=$(accrual-port) \
    -accrual-database-uri="postgresql://postgres:postgres@localhost:5432/praktikum?sslmode=disable"
	rm $(gophermart-bin)

gen:
	go generate ./...

cover:
	go test -short -count=1 -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o=coverage.html
	rm coverage.out

statictest:
	go vet -vettool=./.tools/statictest ./...

dev-up:
	docker-compose -f=docker-compose.dev.yml --env-file=.env.dev up -d

dev-down:
	docker-compose -f=docker-compose.dev.yml --env-file=.env.dev down --rmi local

accrual:
	./cmd/accrual/accrual_linux_amd64 -a=localhost:8787 -d="postgresql://postgres:postgres@localhost:5432/praktikum?sslmode=disable"