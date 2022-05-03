ALL_GO_FILES = $(shell find . -iname '*.go' -not -path './vendor/*')
ENTRYPOINT = $(shell find ./cmd/kak-raw-send/ -iname '*.go')
GO_COVER_FLAG ?= -cover -covermode atomic
CGO_ENABLED ?= 1
GO_TEST_FLAGS ?= -race -shuffle on
CLEAN_FILES = cover.out kak-raw-send
VERSION ?= $(shell git describe --tags HEAD --always)

kak-raw-send: $(ALL_GO_FILES)
	CGO_ENABLED=$(CGO_ENABLED) go build $(GO_BUILD_FLAGS) -o $@ $(ENTRYPOINT)

test: $(ALL_GO_FILES) $(GO_MOD_FILES)
	go test $(GO_TEST_FLAGS) ./...

pre-commit: format test lint

show-coverage: cover.out
	go tool cover -html=$<

cover.out: $(ALL_GO_FILES) $(GO_MOD_FILES)
	$(MAKE) test GO_TEST_FLAGS="$(GO_TEST_FLAGS) $(GO_COVER_FLAG) -coverprofile $@"

format:
	goimports -w -local $$(go list -m) $(ALL_GO_FILES)
	golangci-lint run --no-config --disable-all -E gofumpt --fix

lint:
	golangci-lint run $(LINT_FLAGS)

todo: $(ALL_GO_FILES)
	@grep -niE '//\s?(TODO|FIXME|XXX|NOTE):?\s' $^

clean:
	rm -rf $(CLEAN_FILES)
	go clean -testcache

.PHONY: test show-coverage pre-commit format lint todo clean
