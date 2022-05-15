VERSION_API=develop
container_name_api=mawinter-api
.PHONY: build
build:
	docker build -t ghcr.io/azuki774/$(container_name_api):$(VERSION_API) -f build/Dockerfile .

.PHONY: push
push:	
	docker push ghcr.io/azuki774/$(container_name_api):$(VERSION_API)

.PHONY: test
test:
	gofmt -l -w .
	go test ./... -v -cover

.PHONY: overall-test
overall-test:
	docker-compose -f build/overall-test.yml up --build -d
	sleep 25s
	test/run.sh
	docker-compose -f build/overall-test.yml down

.PHONY: overall-clean
overall-clean:
	docker-compose -f build/overall-test.yml down
