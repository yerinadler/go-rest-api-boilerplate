.PHONY: clean
clean:
	rm -rf bin

build:
	CGO_ENABLED=0 go build -o bin/api -v -ldflags="-s -w"  cmd/api/main.go
	CGO_ENABLED=0 go build -o bin/consumer -v -ldflags="-s -w" cmd/consumer/main.go