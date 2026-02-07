.PHONY: build test lint clean

build:
	go build -o contribemblem cmd/contribemblem/main.go

test:
	go test -v ./...

lint:
	go vet ./...
	test -z $$(gofmt -l .)

clean:
	rm -f contribemblem
