.PHONY: api-build api-test lambda-build frontend-build docker-up docker-down lint

api-build:
	cd src/api && go build ./...

api-test:
	cd src/api && go test ./...

lambda-build:
	cd src/lambdas/transaction_processor && GOOS=linux GOARCH=amd64 go build -o bin/transaction_processor

frontend-build:
	cd frontend && npm install && npm run build

docker-up:
	docker-compose up --build

docker-down:
	docker-compose down --remove-orphans
