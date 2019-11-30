DEFAULT_TARGET := run

.PHONY: build run test-panels test

DEFAULT_APP_NAME ?= "node" 
DEFAULT_APP_IP ?= "127.0.0.1"
DEFAULT_APP_PORT ?= "8080"

build: ## Build the node
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o $(DEFAULT_APP_NAME) .

run: build ## Run the node
	@./$(DEFAULT_APP_NAME) $(DEFAULT_APP_IP) $(DEFAULT_APP_PORT)

test-panels: ## Open multi panels... We assume in the port 127.0.0.1:8080 we have a node
	sleep 3
	@./test/panels.sh

test: ## Test all app
	go test ./... -v
