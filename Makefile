PACKAGES = $(shell go list -f '{{.ImportPath}}/' ./... | grep -v vendor)

BUILDTAGS :=

.PHONY: build clean deploy


init:
	@echo "Initializing this Makefile dependencies..."
	go mod vendor -v

build:
	@echo "Building..."
	env GOOS=linux go build -ldflags="-s -w" -o ./bin/lambda_get ./lambda_get/main.go
	env GOOS=linux go build -ldflags="-s -w" -o ./bin/lambda_fetch_data ./lambda_fetch_data/main.go

	env GOOS=linux go build -ldflags="-s -w" -o ./bin/local ./local/main.go

tests: ## Runs the go tests
	@echo "+ $@"
	@RUNNING_TESTS=1 go test -v -tags "$(BUILDTAGS) cgo" $(PACKAGES)

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose
