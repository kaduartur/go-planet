.PHONY: run test

run:
	go run -race ./cmd/api/main.go

test:
	go test ./... -covermode=atomic -race -count=1