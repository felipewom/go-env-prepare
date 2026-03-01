.PHONY: clean install build run test test-race lint

BIN_NAME := prepare

clean:
	rm -f $(BIN_NAME)

install:
	go mod tidy

build:
	go build -o $(BIN_NAME)

run:
	./$(BIN_NAME)

test:
	go test -count=1 ./...

test-race:
	go test -race -count=1 ./...

lint:
	golangci-lint run ./...

install-homebrew:
	cp $(BIN_NAME) /usr/local/bin/$(BIN_NAME)