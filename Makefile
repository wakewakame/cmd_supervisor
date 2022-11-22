build:
	docker compose build \
		--build-arg UID=$$(id -u) \
		--build-arg GID=$$(id -g)

up: down build
	docker compose up --detach

down:
	docker compose down

bash:
	docker compose exec api /bin/bash
