all: main

main: *.go
	go build -o gopexels main.go

windows: *.go
	GOOS=windows GOARCH=amd64 go build -o gopexels.exe main.go
