default:
	@echo "please specify a target\n - build\n - run\
		\n - utest (unit tests)\n - itest (integration tests)"
build:
	go build .
run:
	go run .
utest:
	go test -v -tags=utest ./...
itest:
	go test -v -tags=itest ./...

