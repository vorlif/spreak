
XSPREAK_PATH=cmd/xspreak

gofmt:
	gofmt -s -w .
	cd $(XSPREAK_PATH) && gofmt -s -w .
	goimports -w -local github.com/vorlif/spreak ./
	cd $(XSPREAK_PATH) && goimports -w -local github.com/vorlif/spreak ./


lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.45.2

	@echo Running golangci-lint
	golangci-lint run --fix ./...
	cd $(XSPREAK_PATH) && golangci-lint run --fix ./...


test-race:
	go test -race -run=. ./... || exit 1;


test:
	go test -short ./...
	cd $(XSPREAK_PATH) && go test -short ./...


coverage:
	go test -short -v -coverprofile cover.out ./...
	go tool cover -func cover.out
	go tool cover -html=cover.out -o coverage.html

coverage-cli:
	cd $(XSPREAK_PATH) && go test -short -v -coverprofile cover.out ./...
	cd $(XSPREAK_PATH) && go tool cover -func cover.out
	cd $(XSPREAK_PATH) && go tool cover -html=cover.out -o coverage.html


install-cli:
	cd $(XSPREAK_PATH) && go install


clean:
	@echo Cleaning

	go clean -i ./...
	rm -f cover.out
	rm -f coverage.html
	rm -rf dist
