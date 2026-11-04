VERSION ?=
SERVICE ?=

.PHONY: check-version check-service setup deploy-user-service kafka-init migrate-up migrate-down migrate-force migrate-version seed up build down

check-version:
	@test -n "$(VERSION)" || (echo "VERSION is required. Example: make migrate-force VERSION=1" && exit 1)
check-service:
	@test -n "$(SERVICE)" || (echo "SERVICE is required. Example: make up SERVICE=user-service" && exit 1)

setup: build
	docker compose up -d --wait kafka greencart-redis
	docker compose up -d --wait user-service-db
	$(MAKE) kafka-init
	$(MAKE) migrate-up SERVICE=user-service
	$(MAKE) seed SERVICE=user-service
	docker compose up -d user-service

deploy-user-service:
	docker compose build user-service
	docker compose up -d --wait kafka greencart-redis user-service-db
	$(MAKE) kafka-init
	$(MAKE) migrate-up SERVICE=user-service
	docker compose up -d user-service

kafka-init:
	docker compose exec -T kafka /opt/kafka/bin/kafka-topics.sh --bootstrap-server kafka:19092 --create --if-not-exists --topic auth.email-verification.v1 --partitions 3 --replication-factor 1 --config retention.ms=259200000
	docker compose exec -T kafka /opt/kafka/bin/kafka-topics.sh --bootstrap-server kafka:19092 --create --if-not-exists --topic auth.forgot-password.v1 --partitions 3 --replication-factor 1 --config retention.ms=259200000

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