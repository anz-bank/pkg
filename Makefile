all: deps test

deps:
	go get -v -t -d ./...

test:
	go test ./...
