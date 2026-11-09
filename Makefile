VERSION ?=
SERVICE ?=

.PHONY: check-version check-service setup deploy-user-service kafka-init migrate-up migrate-down migrate-force migrate-version seed up build down

check-version:
	@test -n "$(VERSION)" || (echo "VERSION is required. Example: make migrate-force VERSION=1" && exit 1)
check-service:
	@test -n "$(SERVICE)" || (echo "SERVICE is required. Example: make up SERVICE=user-service" && exit 1)

setup: build
	docker compose up -d --wait kafka greencart-redis user-service-db s3
	$(MAKE) kafka-init
	$(MAKE) s3-init
	$(MAKE) migrate-up SERVICE=user-service
	$(MAKE) seed SERVICE=user-service
	docker compose up -d user-service
	docker compose up -d nginx

deploy-user-service:
	docker compose build user-service
	docker compose up -d --wait kafka greencart-redis user-service-db s3
	$(MAKE) migrate-up SERVICE=user-service
	docker compose up -d user-service
	$(MAKE) nginx-reload

nginx-reload:
	docker compose exec nginx nginx -t
	docker compose exec nginx nginx -s reload

nginx-restart:
	docker compose up -d --no-deps --force-recreate nginx

kafka-init:
	docker compose run --rm --no-deps kafka_initializer

s3-init:
	docker compose run --rm --no-deps s3_initializer

migrate-up: check-service
	docker compose run --rm $(SERVICE) migrate up

migrate-down: check-service
	docker compose run --rm $(SERVICE) migrate down

migrate-force: check-service check-version
	docker compose run --rm $(SERVICE) migrate force -- $(VERSION)

migrate-version: check-service
	docker compose run --rm $(SERVICE) migrate version

seed: check-service
	docker compose run --rm $(SERVICE) seed

up: check-service
	docker compose up -d $(SERVICE)

build:
	docker compose build

down:
	docker compose down