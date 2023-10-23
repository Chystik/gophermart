SHELL = /bin/bash
.PHONY: autotest dep test race lint gen cover statictest dev-up dev-down

dep:
	go mod download
	go mod tidy

test:
	go test ./...

race:
	go test -v -race ./...

lint:
	/home/user/go/bin/golangci-lint run

gophermart-bin = ./cmd/gophermart/gophermart
autotest:
	go build -o $(gophermart-bin) ./cmd/gophermart
	./gophermarttest \
    -test.v -test.run=^TestGophermart$ \
    -gophermart-binary-path=$(gophermart-bin) \
    -gophermart-host=localhost \
    -gophermart-port=8080 \
    -gophermart-database-uri="postgresql://postgres:postgres@localhost:5432/praktikum?sslmode=disable" \
    -accrual-binary-path=./cmd/accrual/accrual_linux_amd64 \
    -accrual-host=localhost \
    -accrual-port=8787 \
    -accrual-database-uri="postgresql://postgres:postgres@localhost:5432/praktikum?sslmode=disable"
	rm $(gophermart-bin)

gen:
	go generate ./...

cover:
	go test -short -count=1 -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o=coverage.html
	rm coverage.out

statictest:
	go vet -vettool=statictest ./...

dev-up:
	docker-compose -f=docker-compose.dev.yml --env-file=.env.dev up -d

dev-down:
	docker-compose -f=docker-compose.dev.yml --env-file=.env.dev down --rmi local