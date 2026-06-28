BINARY_NAME=look-news
MAIN_PKG=./cmd/web
BIN_DIR=bin

.DEFAULT_GOAL := help

.PHONY: help run build clean test test-cover fmt vet tidy lint

help: ## Mostra esta ajuda
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-12s\033[0m %s\n", $$1, $$2}'

run: ## Roda a aplicação
	go run $(MAIN_PKG)

build: ## Compila o binário em bin/
	go build -o $(BIN_DIR)/$(BINARY_NAME) $(MAIN_PKG)

clean: ## Remove os artefatos de build
	go clean
	rm -rf $(BIN_DIR)

test: ## Roda os testes
	go test ./...

test-cover: ## Roda os testes com cobertura
	go test -cover ./...

fmt: ## Formata o código
	go fmt ./...

vet: ## Analisa o código com go vet
	go vet ./...

tidy: ## Organiza as dependências
	go mod tidy

lint: fmt vet ## Roda fmt e vet
