VERSION_API=develop
container_name_api=mawinter-api
.PHONY: build run push stop test overall-test overall-clean

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

overall-test:
	docker-compose -f deploy/docker/overall-test.yml up --build -d
	sleep 25s
	test/run.sh
	docker-compose -f deploy/docker/overall-test.yml down

overall-clean:
	docker-compose -f deploy/docker/overall-test.yml down
