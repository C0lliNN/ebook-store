.PHONY: dependency unit-test integration-test docker-up docker-down clear

dependency:
	@go get -v ./...

integration-test: export AWS_ACCESS_KEY_ID=test
integration-test: export AWS_SECRET_ACCESS_KEY=test
integration-test: docker-up dependency
	@go test -tags=integration ./...

unit-test: dependency
	@go test -tags=unit ./...

docker-up:
	@docker-compose --file=docker-compose.test.yml --project-name ebook-store-test up -d

docker-down:
	@docker-compose --file=docker-compose.test.yml --project-name ebook-store-test down

migrate-up:
	@migrate -source file:./migration -database postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable up

migrate-down:
	@migrate -source file:./migration -database postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable down 1

clear: docker-down

wire:
	@cd cmd && wire .

generate-mocks:
	@mockery --all --output=./mocks --dir=internal --case=underscore --keeptree

api-docs:
	@swag init -g internal/server/server.go --parseInternal  --generatedTime

start_server:
	@cd cmd && go run .