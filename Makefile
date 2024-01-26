VERSION := v0.3.3

ldflags := $(LDFLAGS)
ldflags += -X main.version=$(VERSION)
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -ldflags '$(ldflags)'

.PHONY: format lint build
all: format lint build

###############################################################################
###                                  Build                                  ###
###############################################################################

build:
	@echo "🤖 Building supervysor..."
	@go build $(BUILD_FLAGS) -mod=readonly -o "$(PWD)/build/" ./cmd/supervysor
	@echo "✅ Completed build!"

install:
	@echo "🤖 Installing supervysor..."
	@go install -mod=readonly $(BUILD_FLAGS)  ./cmd/supervysor
	@echo "✅ Completed installation!"

###############################################################################
###                          Formatting & Linting                           ###
###############################################################################

gofumpt_cmd=mvdan.cc/gofumpt
golangci_lint_cmd=github.com/golangci/golangci-lint/cmd/golangci-lint

format:
	@echo "🤖 Running formatter..."
	@go run $(gofumpt_cmd) -l -w .
	@echo "✅ Completed formatting!"

lint:
	@echo "🤖 Running linter..."
	@go run $(golangci_lint_cmd) run --timeout=10m
	@echo "✅ Completed linting!"
