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
	make usersvc-build
	docker build -t goteam-usersvc ./build/package/usersvc
	docker run -p 8080:8080 -it goteam-usersvc

teamsvc-build: 
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o \
		./build/package/teamsvc/ ./cmd/teamsvc/main.go
	cp .env ./build/package/teamsvc/

teamsvc-run:
	make teamsvc-build
	docker build -t goteam-teamsvc ./build/package/teamsvc
	docker run -p 8081:8081 -it goteam-teamsvc

tasksvc-build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
		-o ./build/package/tasksvc/ ./cmd/tasksvc/main.go
	cp .env ./build/package/tasksvc/

tasksvc-run:
	make tasksvc-build
	docker build -t goteam-tasksvc ./build/package/tasksvc
	docker run -p 8082:8082 -it goteam-tasksvc

backend-build:
	make usersvc-build
	make teamsvc-build
	make tasksvc-build

backend-run:
	make backend-build
	docker compose -f ./build/package/docker-compose.yml up \
		--build --force-recreate --no-deps

backend-stop:
	docker compose -f ./build/package/docker-compose.yml down

frontend-run:
	cd web && NODE_OPTIONS=--openssl-legacy-provider yarn run start

backend-test-u:
	go test -tags=utest ./...

backend-test-uv:
	go test -v -tags=utest ./...

backend-test-i:
	go test -tags=itest ./test/...

backend-test-iv:
	go test -v -tags=itest ./test/...

backend-test:
	make backend-test-u
	make backend-test-i

backend-test-v:
	make backend-test-uv
	make backend-test-iv

