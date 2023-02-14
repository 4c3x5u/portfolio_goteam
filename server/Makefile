.PHONY: default build run utest itest.u itest.d itest itest.v

default:
	@echo "ERROR: target not specified"

build:
	go build .

run:
	go run .

# run unit tests
utest:
	go test -tags=utest ./...

# set [u]p integration test environment
itest.u:
	docker-compose --project-directory ./itest up -d

# tear [d]own integration test environment
itest.d:
	docker-compose --project-directory ./itest down

# run integration tests
itest:
	@make itest.u
	go test -tags=itest ./itest
	@make itest.d

# run integration tests in [v]erbose mode
itest.v:
	@make itest.u
	go test -v -tags=itest ./itest
	@make itest.d
