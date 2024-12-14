lint:
	go fmt ./...

build:
	go build -o ./bin/ftree ./cmd/ftree/main.go

run:
	go run ./cmd/ftree/main.go

install: build
	cp ./bin/ftree ~/.local/bin/ftree
