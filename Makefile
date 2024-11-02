BINARY_NAME=topic-tracer
all: build

build:
	@ go build -o $(BINARY_NAME) main.go

run: build
	@ ./$(BINARY_NAME) search -q "your search query" -d "your/directory/path"
test:
	@ go test ./...
clean:
	@ rm -f $(BINARY_NAME)
deps:
	@ go mod tidy

.PHONY: all build run test clean deps