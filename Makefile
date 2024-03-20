.PHONY: clean
clean:
	rm -rf bin

build:
	go build -o bin/api -v -ldflags="-s -w"  cmd/api/main.go
	go build -o bin/consumer -v -ldflags="-s -w" cmd/consumer/main.go