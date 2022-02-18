go-run-programs:
	go build .

lint:
	docker run -it --rm -v $(PWD):/app -w /app golangci/golangci-lint:v1.44.2 golangci-lint run -v --enable-all ./...

test:
	go test -v .
