# ignore go cache for tests
export GOFLAGS = -count=1

default: clean lint test

go-run-programs:
	go build .

# go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.45.0
lint:
	golangci-lint run
.PHONY: lint

test:
	 go test -v ./...
.PHONY: test

test-%:
	go test -v ./... -run $*
.PHONY: test-%

clean:
	rm -rf go-run-programs ./tests/test-server{,.exe}
.PHONY: clean
