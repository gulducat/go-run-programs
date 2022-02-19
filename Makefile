go-run-programs:
	go build ./cmd/go-run-programs

lint:
	docker run -it --rm -v $(PWD):/app -w /app golangci/golangci-lint:v1.44.2 golangci-lint run -v --enable-all ./...

test:
	go test -v .

clean:
	rm -f go-run-programs

.PHONY: lint test clean
