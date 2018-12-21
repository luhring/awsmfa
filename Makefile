all: clean build

build:
	go build -o build/awsmfa

release: test
	goreleaser --rm-dist

install:
	go install ./...

test: build validate
	go test -cover -v ./...

validate:
	@! gofmt -s -d -l . 2>&1 | grep -vE '^\.git/'
	go vet ./...

clean:
	rm -rf build
	go clean

.PHONY: build install test clean release validate
