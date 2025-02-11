SHELL=/bin/bash
VERSION_API=latest
container_name_api=mawinter-api
container_name_doc=mawinter-doc

.PHONY: build bin bin-linux-amd64 start stop migration test doc generate

build:
	docker build -t $(container_name_api):$(VERSION_API) -f build/Dockerfile .

bin:
	go build -a -tags "netgo" -installsuffix netgo  -ldflags="-s -w -extldflags \"-static\" \
	-X main.version=$(git describe --tag --abbrev=0) \
	-X main.revision=$(git rev-list -1 HEAD) \
	-X main.build=$(git describe --tags)" \
	-o bin/ ./...

bin-linux-amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags "netgo" -installsuffix netgo  -ldflags="-s -w -extldflags \"-static\" \
	-X main.version=$(git describe --tag --abbrev=0) \
	-X main.revision=$(git rev-list -1 HEAD) \
	-X main.build=$(git describe --tags)" \
	-o bin/ ./...

start: 
	docker compose -f deployment/compose-local.yml up -d

stop:
	docker compose -f deployment/compose-local.yml down

migration:
	cd migration; \
	sql-migrate up -env=local; \
	cd ../

test: 
	gofmt -l .
	go vet ./...
	staticcheck ./...
	go test -v ./...

coverage:
	go test -coverprofile=docs/coverage.out ./...
	go tool cover -html=docs/coverage.out -o docs/coverage.html

generate:
	oapi-codegen -package "openapi" -generate "chi-server" internal/openapi/mawinter-api.yaml > internal/openapi/server.gen.go
	oapi-codegen -package "openapi" -generate "spec" internal/openapi/mawinter-api.yaml > internal/openapi/spec.gen.go
	oapi-codegen -package "openapi" -generate "types" internal/openapi/mawinter-api.yaml > internal/openapi/types.gen.go
