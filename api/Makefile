BINARY_NAME=look-news
MAIN_PKG=./cmd/web
BIN_DIR=bin

FILTRO_DIR=../llm-filter
FILTRO_VENV=$(FILTRO_DIR)/.venv
FILTRO_PYTHON=$(FILTRO_VENV)/bin/python3
FILTRO_SCRIPT=$(FILTRO_DIR)/main.py
FILTRO_QUERY = High-quality factual journalism and news updates strictly about technology, world affairs, geopolitics, global economy, science, engineering breakthroughs and local developments in Angola, with strong priority to new open-source/research projects, tech product launches, AI, machine learning, robotics, space, astronomy and other deep-tech/nerd topics; immediately reject every form of advertisement, sponsored content, press release, product promotion, affiliate marketing, coupon/deal roundups, clickbait and low-quality SEO articles.

.DEFAULT_GOAL := help

.PHONY: help run build clean test test-cover fmt vet tidy lint setup-filtro test-filtro

help: ## Mostra esta ajuda
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-12s\033[0m %s\n", $$1, $$2}'

run: ## Roda a aplicação
	FILTRO_PYTHON=$(FILTRO_PYTHON) FILTRO_SCRIPT=$(FILTRO_SCRIPT) FILTRO_QUERY="$(FILTRO_QUERY)" go run $(MAIN_PKG)

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

setup-filtro: ## Cria o venv do filtro Python e instala as dependências
	python3 -m venv $(FILTRO_VENV)
	$(FILTRO_VENV)/bin/pip install --upgrade pip
	$(FILTRO_VENV)/bin/pip install -r $(FILTRO_DIR)/requirements.txt

test-filtro: ## Testa o filtro Python isolado (sem o Go)
	@echo '{"query":"artificial intelligence","articles":[{"link":"https://a.com","title":"GPT-5 released","summary":"new OpenAI model"},{"link":"https://b.com","title":"Carrot cake recipe","summary":"how to bake"}]}' | $(FILTRO_PYTHON) $(FILTRO_SCRIPT)
