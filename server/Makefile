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

# all tests [v]erbose
test.v:
	make utest.v
	make itest.v

# unit tests
utest:
	go test -tags=utest ./...

# unit tests [v]erbose
utest.v:
	go test -v -tags=utest ./...

# integration tests
itest:
	go test -tags=itest ./itest

# integration tests [v]erbose
itest.v:
	go test -v -tags=itest ./itest

