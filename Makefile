APP_NAME=gitglance
BINARY_PATH=bin/${APP_NAME}

all: build run test
.PHONY: all

build: 
	go build -o ${BINARY_PATH}

run: 
	./${BINARY_PATH}

test:
	 go test -v ./...

clean:
	go clean
	rm -rf bin
