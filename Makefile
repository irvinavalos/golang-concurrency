build:
	@go build -o ./bin/chat ./...
	@chmod +x ./bin/chat

chat: build
	@./bin/chat
