.PHONY: default run test test-v test-u test-uv test-i test-iv

default: 
	@echo "ERROR: target not specified"

run: 
	go run ./cmd/api

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
