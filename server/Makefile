include .env
export

start:
	./build/ivory

.PHONY: build
build:
	go build -o build/ivory

test:
	go test ./...

clean:
	rm -rf build

update:
	go get -u ./...
	go mod tidy
