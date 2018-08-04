help:
	@echo "test             run test"
	@echo "lint             run lint"
	@echo "example          run examples"

.PHONY: test
test:
	go test -v -cover -coverprofile cover.out
	go tool cover -html=cover.out -o cover.html
	-open cover.html

.PHONY: lint
lint:
	gofmt -s -w . _example
	goimports -w . _example
	golint . _example
	go vet

.PHONY: example
example:
	cd _example && bash test.sh
