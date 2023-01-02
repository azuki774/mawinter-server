SHELL=/bin/bash
VERSION_API=latest
container_name_api=mawinter-api
container_name_db=mawinter-db
.PHONY: build run push stop test migration-test migration-clean

build:
	docker build -t $(container_name_api):$(VERSION_API) -f build/Dockerfile .
	docker build -t $(container_name_db):$(VERSION_API) -f build/Dockerfile-db .

bin:
	go build -a -tags "netgo" -installsuffix netgo  -ldflags="-s -w -extldflags \"-static\"" -o bin/ ./...

bin-linux-amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags "netgo" -installsuffix netgo  -ldflags="-s -w -extldflags \"-static\"" -o bin/ ./...

start: 
	docker compose -f deployment/compose-local.yml up -d

stop:
	docker compose -f deployment/compose-local.yml down

test: 
	gofmt -l .
	go vet ./...
	staticcheck ./...
	go test -v ./...
