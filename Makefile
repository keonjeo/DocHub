.PHONY: fmt run air test dep build
# NAME = $(notdir $(shell pwd))

fmt:
	@echo "formatting..."
	@gofmt -w ./../dochub
air:
	make fmt && air
run:
	make fmt && go run *.go
dep:
	make fmt && go get -v -x && go mod download && go mod tidy
build:
	make fmt && GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build
test:
	go test -count=1
