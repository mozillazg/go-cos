help:
	@echo "lint             run lint"

.PHONY: lint
lint:
	gofmt -s -w .
	golint .
	go vet
