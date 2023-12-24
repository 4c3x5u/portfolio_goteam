.PHONY: default run test test-v test-u test-uv test-i test-iv

default: 
	@echo "ERROR: target not specified"

build-user:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./build/package/usersvc/ ./cmd/usersvc/main.go
	cp .env ./build/package/usersvc/

build-team: 
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./build/package/teamsvc/ ./cmd/teamsvc/main.go
	cp .env ./build/package/teamsvc/

build-task:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./build/package/tasksvc/ ./cmd/tasksvc/main.go
	cp .env ./build/package/tasksvc/

build-all:
	make build-user
	make build-team
	make build-task
	docker compose build -f ./build/package/docker-compose.yml

run-all:
	make build-all
	docker compose -f ./build/package/docker-compose.yml up

test:
	make test-u
	make test-i

test-v:
	make test-uv
	make test-iv

test-u:
	go test -tags=utest ./...

test-uv:
	go test -v -tags=utest ./...

test-i:
	go test -tags=itest ./test/...

test-iv:
	go test -v -tags=itest ./test/...
