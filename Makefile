default: 
	@echo "ERROR: target not specified"

db-run:
	docker build -t goteam-db ./build/package/db
	docker run \
		--mount source=data,target=/home/dynamodblocal/data \
		-p 8000:8000 \
		-it goteam-db

db-init:
	./build/package/db/init.sh

usersvc-build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
		-o ./build/package/usersvc/ ./cmd/usersvc/main.go
	cp .env ./build/package/usersvc/

usersvc-run:
	make build-usersvc
	docker build -t goteam-usersvc ./build/package/usersvc
	docker run -p 8080:8080 -it goteam-usersvc

teamsvc-build: 
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o \
		./build/package/teamsvc/ ./cmd/teamsvc/main.go
	cp .env ./build/package/teamsvc/

teamsvc-run:
	make build-teamsvc
	docker build -t goteam-teamsvc ./build/package/teamsvc
	docker run -p 8081:8081 -it goteam-teamsvc

tasksvc-build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
		-o ./build/package/tasksvc/ ./cmd/tasksvc/main.go
	cp .env ./build/package/tasksvc/

tasksvc-run:
	make build-tasksvc
	docker build -t goteam-tasksvc ./build/package/tasksvc
	docker run -p 8082:8082 -it goteam-tasksvc

backend-build:
	make build-usersvc
	make build-teamsvc
	make build-tasksvc

backend-run:
	make build-be
	docker compose -f ./build/package/docker-compose.yml up \
		--build --force-recreate --no-deps

backend-stop:
	docker compose -f ./build/package/docker-compose.yml down

frontend-run:
	cd web && NODE_OPTIONS=--openssl-legacy-provider yarn run start

backend-test:
	make test-u
	make test-i

backend-test-v:
	make test-uv
	make test-iv

backend-test-u:
	go test -tags=utest ./...

backend-test-uv:
	go test -v -tags=utest ./...

backend-test-i:
	go test -tags=itest ./test/...

backend-test-iv:
	go test -v -tags=itest ./test/...
