.PHONY: api-build api-test lambda-build frontend-build docker-up docker-down lint

api-build:
	go build ./cmd/api

api-test:
	go test ./...

lambda-build:
	cd cmd/lambdas/transaction_processor && GOOS=linux GOARCH=amd64 go build -o bin/transaction_processor

frontend-build:
	cd frontend && npm install && npm run build

docker-up:
	docker-compose up --build

docker-down:
	docker-compose down --remove-orphans
