build:
	go build -o bin/main main.go

run:
	go run main.go

format:
	go fmt main.go

clean:
	go clean
	rm -rf bin
