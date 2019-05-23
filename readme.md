A simple router written in Go for the [Serverless](https://serverless.com/) framework.

Currently due Go limitations needs to have a separate binary for each function, which is quite annoying. This library allows to have a single binary fo multiple functions.

## Installation
```go get github.com/Napas/go-serverless-router```

or if using [dep](https://github.com/golang/dep):
```dep ensure -v -add github.com/Napas/go-serverless-router```

## Usage example
A small example app could be found in the [example](./example) folder.