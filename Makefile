SRC=*.go
BIN=s3-cli

$(BIN): $(SRC)
	CGO_ENABLED=0 go build -ldflags '-w -s' -o $@ .

lint:
	golangci-lint run

test:
	go test

clean:
	rm -f $(BIN)
