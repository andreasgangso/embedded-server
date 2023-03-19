install:
	go install github.com/GeertJohan/go.rice/rice@latest

build:
	rice embed-go
	GOOS=darwin GOARCH=amd64 go build -o bin/darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -o bin/darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build -o bin/win-amd64.exe .
	GOOS=linux GOARCH=amd64 go build -o bin/linux-amd64 .
	GOOS=linux GOARCH=arm64 go build -o bin/linux-arm64 .
	GOOS=linux GOARCH=386 go build -o bin/linux-386 .