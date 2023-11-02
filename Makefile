GITCOMMIT := $(shell git rev-parse HEAD)
GITDATE := $(shell git show -s --format='%ct')

LDFLAGSSTRING +=-X main.GitCommit=$(GITCOMMIT)
LDFLAGSSTRING +=-X main.GitDate=$(GITDATE)
LDFLAGS := -ldflags "$(LDFLAGSSTRING)"

GOPATH:=$(shell go env GOPATH)

.PHONY: all
all:
	@make tidy && make build && make test && make lint

.PHONY: update
update:
	@go get -u

.PHONY: tidy
tidy:
	@go mod tidy

.PHONY: run
run:
	@go run ./cmd config --config ./example/mocktimism.toml

.PHONY: build
build:
	@env GO111MODULE=on go build -v $(LDFLAGS) -o bin/mocktimism ./cmd

.PHONY: test
test:
	@go test -v ./...

.PHONY: docker
docker:
	@docker build -t mocktimism:latest .

.PHONY: clean
clean:
	@rm bin/mocktimism

.PHONY: lint
lint:
	@golangci-lint run -E goimports,sqlclosecheck,bodyclose,asciicheck,misspell,errorlint --fix --timeout 5m -e "errors.As" -e "errors.Is" ./...

.PHONY: lintcheck
lintcheck:
	@golangci-lint run -E goimports,sqlclosecheck,bodyclose,asciicheck,misspell,errorlint --timeout 5m -e "errors.As" -e "errors.Is" ./...

# To keep mocktimism as DRY as possible we want to use the existing optimism scripts to generate the genesis allocs for l1
.PHONY: generate-allocs
generate-allocs:
	$(MAKE) -C lib/optimism devnet-allocs
	cp -rf lib/optimism/.devnet/* generated/
