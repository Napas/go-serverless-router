.PHONY: build clean deploy

build:
	GOOS=linux go build -v -o bin/service main.go

deploy: build
	sls deploy --verbose
