VERSION_API=develop
container_name_api=mawinter-api
.PHONY: build run push stop test migration-test migration-clean

build:
	docker build -t azuki774/$(container_name_api):$(VERSION_API) -f build/Dockerfile .

run:
	docker-compose -f deploy/docker/docker-compose.yml up -d

stop:
	docker-compose -f deploy/docker/docker-compose.yml down

push:
	docker tag $(container_name_api) ghcr.io/azuki774/$(container_name_api):develop
	docker push ghcr.io/azuki774/$(container_name_api):$(VERSION_API)

test:
	gofmt -l -w .
	go test ./... -v -cover

migration-test:
	docker-compose -f deploy/docker/migration-test.yml up --build -d
	sleep 25s
	python3 test/check.py
	# test/run.sh
	docker-compose -f deploy/docker/migration-test.yml down

migration-clean:
	docker-compose -f deploy/docker/migration-test.yml down
