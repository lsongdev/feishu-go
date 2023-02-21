BINARY_NAME=bin/feishu-bot

build:
	GOARCH=amd64 GOOS=darwin go build -ldflags="-s -w" -o ${BINARY_NAME}-darwin-amd64 main.go
	GOARCH=arm64 GOOS=darwin go build -ldflags="-s -w" -o ${BINARY_NAME}-darwin-arm64 main.go
	GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o ${BINARY_NAME}-linux-amd64 main.go
	GOARCH=amd64 GOOS=windows go build -ldflags="-s -w" -o ${BINARY_NAME}-amd64.exe main.go

clean:
	go clean
	rm ${BINARY_NAME}-darwin
	rm ${BINARY_NAME}-linux
	rm ${BINARY_NAME}-x64.exe
