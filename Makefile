all: build

BINARY_NAME=rxmer
VERSION = 0.1

build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o releases/${BINARY_NAME}${VERSION}-linux-amd64
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o releases/${BINARY_NAME}${VERSION}-darwin-amd64
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o releases/${BINARY_NAME}${VERSION}-windows-amd64.exe

clean:
	rm releases/*