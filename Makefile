SRC=*.go
BIN=s3-cli

$(BIN): $(SRC)
	CGO_ENABLED=0 go build -ldflags '-w -s' -o $@ .

clean:
	rm -f $(BIN)

test:
	go test
