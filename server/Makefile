.PHONY: default build run utest itest.setup itest.teardown itest itest.v

default:
	@echo "ERROR: target not specified"

build:
	go build .

run:
	go run .

utest:
	go test -tags=utest ./...

itest.setup:
	docker-compose --project-directory ./itest up -d

itest.teardown:
	docker-compose --project-directory ./itest down

itest:
	@make itest.setup
	go test -tags=itest ./itest
	@make itest.teardown

itest.v:
	@make itest.setup
	go test -v -tags=itest ./itest
	@make itest.teardown
