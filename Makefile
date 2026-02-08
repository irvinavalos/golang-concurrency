build:
	@go build -o ./bin/chat ./...
	@chmod +x ./bin/chat

chat: build
	@./bin/chat

test-race:
	@go clean -testcache
	@go test -race -v ./...
