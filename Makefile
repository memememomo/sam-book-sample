DOCKER_YAML=-f docker-compose.yml

DOCKER=COMPOSE_PROJECT_NAME=sam-book-sample docker-compose $(DOCKER_YAML)

build:
	$(DOCKER) build --no-cache

clean:
	find handlers -name main -type f | xargs rm

handler:
	$(DOCKER) run sam-local ./scripts/build-handlers.sh

go-dep:
	$(DOCKER) run go-test dep ensure

go-lint:
	$(DOCKER) run go-test ./scripts/go-lint.sh

go-test:
	$(DOCKER) run go-test ./scripts/go-test.sh '${PACKAGE}' '${ARGS}'

generate-random-name:
	$(DOCKER) run sam-local ./scripts/generate-random-name.sh

create-swagger-bucket:
	$(DOCKER) run sam-local ./scripts/create-swagger-bucket.sh

deploy: clean
	$(DOCKER) run sam-local dep ensure
	$(DOCKER) run sam-local ./scripts/build-handlers.sh
	$(DOCKER) run sam-local ./scripts/deploy.sh

delete-stack:
	$(DOCKER) run sam-local ./scripts/delete-stack.sh
