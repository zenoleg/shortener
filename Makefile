.PHONY: dc run test lint

dc:
	docker-compose up  --remove-orphans --build

run:
	go build -o shortener cmd/shortener/main.go && ./shortener

test:
	go test -race ./...

lint:
	docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:v1.59.0-alpine golangci-lint run --timeout=10m
