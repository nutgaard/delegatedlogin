.PHONY: all start test stop

all: start test stop

start:
	docker-compose up --build -d

stop:
	docker-compose down --remove-orphans

test:
	deno test --allow-env --allow-net
