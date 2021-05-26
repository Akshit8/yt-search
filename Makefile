fmt:
	@echo "formatting code"
	go fmt ./...

lint:
	@echo "Linting source code"
	golint ./...

vet:
	@echo "Checking for code issues"
	go vet ./...

test:
	@echo "running tests"
	go test ./...

install:
	@echo "installing external dependencies"
	go mod download

run:
	go run main.go --env .env --address :8000

live:
	reflex -r '\.go' -s -- sh -c "make run"

dev:
	docker-compose -f dev.yml up -d

generate:
	go generate ./...
