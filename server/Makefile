.PHONY: default build run test test.v utest utest.v itest itest.v

default:
	@echo "ERROR: target not specified"

build:
	go build .

run:
	go run .

# all tests
test:
	make utest
	make itest

# all tests verbose
test.v:
	make utest.v
	make itest.v

# unit tests
utest:
	go test -tags=utest ./...

# unit tests verbose
utest.v:
	go test -v -tags=utest ./...

# integration tests
itest:
	go test -tags=itest ./itest

# integration tests verbose
itest.v:
	go test -v -tags=itest ./itest

