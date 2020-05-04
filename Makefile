.PHONY: examples
examples:
	@go run ./examples

.PHONY: test
test:
	@go test ./tests

.PHONY: fmt
fmt:
	@gofmt -w $(shell find $(shell pwd) -name "*.go")

.PHONY: init
init:
	@git init
	@ln -s $(shell pwd)/hooks/pre-commit $(shell pwd)/.git/hooks/pre-commit || true
	@chmod +x $(shell pwd)/.git/hooks/pre-commit

.PHONY: publish
publish:
	@echo publish
