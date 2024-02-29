.PHONY: default
default: build

all: clean build test

build:
	mkdir -p bin
	go build -o bin/genai main.go

test: 
	go test -short -coverprofile=cov.out `go list ./... | grep -v vendor/`
	go tool cover -func=cov.out

clean:
	rm -rf cov.out

sonar: 
	sonar-scanner
