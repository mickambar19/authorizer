generate-mock:
	go generate -v ./...

.ONESHELL:
create-test-html-coverage:
	 go test -coverprofile=coverage.out ./...;
	 go tool cover -html=coverage.out
