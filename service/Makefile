include .env
export

start:
	./build/ivory

.PHONY: build
build:
	go build -o ../build/ivory -C src

clean:
	rm -r build

update:
	go get -u ./...
	go mod tidy
