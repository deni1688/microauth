build:
	docker compose build --no-cache
run: build
	docker compose up
run-silent: build
	docker compose up -d
test:
	make run-silent && sleep 1 && \
	go test -v tests/api_integration_test.go
clean:
	docker compose down && \
	docker rmi microauth:latest


