default: fmt lint install generate

build:
	go build -v ./...

install: build
	go install -v ./...

lint:
	golangci-lint run

generate:
	cd tools; go generate ./...

fmt:
	gofmt -s -w -e .

test:
	go test -v -cover -timeout=120s -parallel=10 ./...

testacc:
	TF_ACC=1 go test -v -cover -timeout 120m ./...

setup-local:
	cd local && ./trino-setup.sh && ./backend-setup.sh && cd -

clean-local:
	./local/clean.sh

.PHONY: fmt lint test testacc build install generate
