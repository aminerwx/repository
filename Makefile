run:
	@go run main.go

build:
	@go build -o bin/repository main.go

clean:
	@rm bin/repository

test:
	go test ./... -coverprofile=./cover.out

cover: test
	@go tool cover -html=./cover.out

all: clean build
	@./bin/repository
