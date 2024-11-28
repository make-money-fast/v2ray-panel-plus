$(shell mkdir -p bin)

.PHONY: build
build:
	@go build -o bin/v2raypanel cmd/client/main.go

build-server:
	@go build -o bin/v2ray cmd/server/main.go