.PHONY: default build run utest itest.u itest.d itest itest.v

default:
	@echo "ERROR: target not specified"

build:
	go build .

run:
	go run .

# unit tests
utest:
	go test -tags=utest ./...

# integration tests
itest:
	go test -tags=itest ./itest

# integration tests [v]erbose
itest.v:
	go test -v -tags=itest ./itest
