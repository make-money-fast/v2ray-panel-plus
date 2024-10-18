@echo off
go generate
go install github.com/tc-hib/go-winres@latest
go-winres simply --icon web/static/favicon.ico
go build -ldflags "-H windowsgui" -tags=release -o bin/win/v2raypanel.exe cmd/client/main.go

