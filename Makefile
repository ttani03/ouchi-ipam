# Variables
APP_NAME=ouchi-ipam

# Targets
.PHONY: all build clean fmt sqlc-gen swag

all: build

build: sqlc-gen swag
	go build -o bin/$(APP_NAME)

clean:
	rm -f bin/$(APP_NAME)

fmt:
	go fmt
	swag fmt

sqlc-gen:
	sqlc generate

swag:
	swag init
