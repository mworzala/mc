
build:
	go build -o mc -ldflags "-X main.source=yes" .

test:
	go test -v ./...
