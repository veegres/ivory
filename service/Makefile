include .env
export

start:
	./build/ivory

.PHONY: build
build:
	go build -C src -o ../build/ivory

test:
	go test -C src ./...

clean:
	rm -r build

update:
	go get -u ./...
	go mod tidy
