.PHONY: clean install build run

BIN_NAME := prepare

clean:
	rm -f $(BIN_NAME)

install:
	go mod tidy

build:
	go build -o $(BIN_NAME)

run:
	./$(BIN_NAME)

install-homebrew:
	cp $(BIN_NAME) /usr/local/bin/$(BIN_NAME)