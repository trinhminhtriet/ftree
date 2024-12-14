lint:
	go fmt ./...

build:
	go build -o ./bin/ftree main.go

run:
	go run main.go

install: build
	cp ./bin/ftree ~/.local/bin/ftree
