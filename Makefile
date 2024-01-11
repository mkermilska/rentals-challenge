start:
	DOCKER_BUILDKIT=1 COMPOSE_DOCKER_CLI_BUILD=1 COMPOSE_HTTP_TIMEOUT=180 DEBUG=true docker-compose build
	docker-compose up

clear:
	docker-compose down

