build:
	go build -tags netgo -o microauth .
	docker-compose build --no-cache
run: build
	docker-compose up
