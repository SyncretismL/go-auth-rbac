lint:
	golangci-lint run ./... --no-config
tests:
	go test ./...
run:
	docker-compose up
stop:
	docker-compose stop