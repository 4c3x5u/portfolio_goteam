default: 
	@echo "ERROR: target not specified"

run-db:
	docker build -t goteam-db ./build/package/db
	docker run \
		--mount source=data,target=/home/dynamodblocal/data \
		-p 8000:8000 \
		-it goteam-db

init-db:
	./build/package/db/init.sh

build-usersvc:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
		-o ./build/package/usersvc/ ./cmd/usersvc/main.go
	cp .env ./build/package/usersvc/

run-usersvc:
	make build-usersvc
	docker build -t goteam-usersvc ./build/package/usersvc
	docker run -p 8080:8080 -it goteam-usersvc

build-teamsvc: 
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o \
		./build/package/teamsvc/ ./cmd/teamsvc/main.go
	cp .env ./build/package/teamsvc/

run-teamsvc:
	make build-teamsvc
	docker build -t goteam-teamsvc ./build/package/teamsvc
	docker run -p 8081:8081 -it goteam-teamsvc

build-tasksvc:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
		-o ./build/package/tasksvc/ ./cmd/tasksvc/main.go
	cp .env ./build/package/tasksvc/

run-tasksvc:
	make build-tasksvc
	docker build -t goteam-tasksvc ./build/package/tasksvc
	docker run -p 8082:8082 -it goteam-tasksvc

build-b:
	make build-usersvc
	make build-teamsvc
	make build-tasksvc

run-b:
	make build-be
	docker compose -f ./build/package/docker-compose.yml up \
		--build --force-recreate --no-deps

stop-b:
	docker compose -f ./build/package/docker-compose.yml down

run-f:
	cd web && NODE_OPTIONS=--openssl-legacy-provider yarn run start

test-b:
	make test-u
	make test-i

test-bv:
	make test-uv
	make test-iv

test-bu:
	go test -tags=utest ./...

test-buv:
	go test -v -tags=utest ./...

test-bi:
	go test -tags=itest ./test/...

test-biv:
	go test -v -tags=itest ./test/...
