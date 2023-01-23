build:
	GOOS=linux GOARCH=amd64 go build -tags netgo -o microauth .
	docker-compose build --no-cache
run: build
	docker-compose up
