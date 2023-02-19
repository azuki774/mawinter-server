SHELL=/bin/bash
VERSION_API=latest
container_name_api=mawinter-api
container_name_register=mawinter-register

.PHONY: build run push stop test migration doc

build:
	docker build -t $(container_name_api):$(VERSION_API) -f build/Dockerfile .
	docker build -t $(container_name_register):$(VERSION_API) -f build/Dockerfile-register .

bin:
	go build -a -tags "netgo" -installsuffix netgo  -ldflags="-s -w -extldflags \"-static\"" -o bin/ ./...

bin-linux-amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags "netgo" -installsuffix netgo  -ldflags="-s -w -extldflags \"-static\"" -o bin/ ./...

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

doc:
	# req: create doc by tbls
	./docs/build_md.sh 2> /dev/null
	cp -a docs/schema/*.svg docs/build/
