.PHONY: default build run utest itest
default:
	@echo "please specify a target\n - build\n - run\
		\n - utest (unit tests)\n - itest (integration tests)"
build:
	go build .
run:
	go run .
utest:
	go test -tags=utest ./...
itest:
	go test -tags=itest ./itest
itest.v:
	go test -v -tags=itest ./itest
