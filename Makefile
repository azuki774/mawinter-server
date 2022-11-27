VERSION_API=latest
container_name_db=mawinter-db
.PHONY: build run push stop test migration-test migration-clean

build:
	go build -a -tags "netgo" -installsuffix netgo  -ldflags="-s -w -extldflags \"-static\"" -o bin/ ./...
	docker build -t $(container_name_db):$(VERSION_API) -f build/Dockerfile-db .

bin:
	go build -a -tags "netgo" -installsuffix netgo  -ldflags="-s -w -extldflags \"-static\"" -o bin/ ./...

bin-linux-amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags "netgo" -installsuffix netgo  -ldflags="-s -w -extldflags \"-static\"" -o bin/ ./...

start: 
	docker compose -f deployment/compose-local.yml up -d

stop:
	docker compose -f deployment/compose-local.yml down
