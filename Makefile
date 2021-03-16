build:
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/main main.go

run:
	go run main.go

format:
	go fmt main.go

clean:
	go clean
	rm -rf bin

deploy: clean build
	yarn deploy
