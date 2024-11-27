.PHONY: build-macos
build-macos:
	@rm -rf bin/mac
	@mkdir -p bin/mac
	@GOOS=darwin GOARCH=arm64 go build -o bin/mac/v2raypanel cmd/client/main.go

server:
	@rm -rf bin/linux/
	@mkdir -p bin/linux
	@GOOS=linux GOARCH=amd64 go build -o bin/linux/v2ray cmd/server/main.go