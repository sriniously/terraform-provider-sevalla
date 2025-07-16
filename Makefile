default: install

build:
	go build -v ./...

install: build
	go install -v ./...

test:
	go test -v -cover -timeout=120s -parallel=10 ./...

testacc:
	TF_ACC=1 go test -v -cover -timeout 120m ./...

fmt:
	gofmt -s -w -e .

lint:
	golangci-lint run

generate:
	go generate ./...

docs:
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate -provider-name sevalla

clean:
	go clean -testcache
	rm -rf dist/

release:
	goreleaser release --clean

snapshot:
	goreleaser release --snapshot --clean

.PHONY: build install test testacc fmt lint generate docs clean release snapshot