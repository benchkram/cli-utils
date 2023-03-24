
VERSION=$(shell git describe --abbrev=0 --tags)
COMMIT=$(shell git rev-parse --short HEAD)

build:
	go build -o cli_example "-ldflags=-X main.Version=$(VERSION) -X main.CommitHash=$(COMMIT)" main.go
