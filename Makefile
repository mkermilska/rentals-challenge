start:
	DOCKER_BUILDKIT=1 COMPOSE_DOCKER_CLI_BUILD=1 DEBUG=true docker-compose build
	docker-compose up postgres rentals-api

clear:
	docker-compose down

integration-tests: 
	docker-compose up venom

unit-tests:
	go test ./... -v 