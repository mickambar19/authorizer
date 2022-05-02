all: build
test: lint unit-test
generate-mock:
	go generate -v ./...
