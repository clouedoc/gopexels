all: main

main: *.go
	go build -o main main.go
