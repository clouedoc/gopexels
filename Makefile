all: gopexels

gopexels: *.go
	go build -o gopexels *.go
