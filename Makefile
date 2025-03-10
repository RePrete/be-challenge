PHONY: run build test proto up down

run:
	docker run --rm -p 8080:8080 entity-status-api:latest

build:
	docker build -t entity-status-api:latest .

test:
	go test -v ./...

proto:
	docker run --rm -v $(PWD):/app entity-status-api:latest sh /app/generate.sh

up:
	docker-compose up --force-recreate

down:
	docker-compose down

logs:
	docker-compose logs -f
