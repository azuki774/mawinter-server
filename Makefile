VERSION_API=latest
container_name_db=mawinter-db
.PHONY: build run push stop test migration-test migration-clean

build:
	docker build -t $(container_name_db):$(VERSION_API) -f build/Dockerfile-db .

start:
	docker compose -f deployment/compose-local.yml up -d

stop:
	docker compose -f deployment/compose-local.yml down
