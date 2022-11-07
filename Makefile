

gofmt:
	gofmt -s -w .
	goimports -w -local github.com/vorlif/spreak ./


lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.1

	@echo Running golangci-lint
	golangci-lint run --fix ./...


test-race:
	go test -race -run=. ./... || exit 1;


test:
	go test -short ./...


coverage:
	go test -short -v -coverprofile cover.out ./...
	go tool cover -func cover.out
	go tool cover -html=cover.out -o coverage.html


clean:
	@echo Cleaning

	go clean -i ./...
	rm -f cover.out
	rm -f coverage.html
	rm -rf dist
