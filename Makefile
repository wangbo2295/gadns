# GADNS Makefile

# 项目信息
NAME := gadns
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Go 相关
GOCMD := go
GOBUILD := $(GOCMD) build
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod
GOFMT := $(GOCMD) fmt
GOVET := $(GOCMD) vet

# 构建相关
LDFLAGS := -ldflags \
	"-X 'github.com/wangbo2295/gadns/cmd.Version=$(VERSION)' \
	-X 'github.com/wangbo2295/gadns/cmd.BuildTime=$(BUILD_TIME)' \
	-X 'github.com/wangbo2295/gadns/cmd.GitCommit=$(GIT_COMMIT)'"

# 目录
CONFIG_DIR := ~/.gadns

# 颜色输出
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[1;33m
NC := \033[0m # No Color

.PHONY: all build test clean install deps fmt vet lint help run

## all: 默认目标，构建项目
all: build

## build: 构建 CLI 工具
build:
	@echo "$(GREEN)Building $(NAME)...$(NC)"
	$(GOBUILD) $(LDFLAGS) -o bin/$(NAME) .
	@echo "$(GREEN)Build complete: bin/$(NAME)$(NC)"

## build-dev: 开发模式构建（带调试信息）
build-dev:
	@echo "$(GREEN)Building $(NAME) (dev mode)...$(NC)"
	$(GOBUILD) -gcflags="all=-N -l" $(LDFLAGS) -o bin/$(NAME) .
	@echo "$(GREEN)Build complete: bin/$(NAME)$(NC)"

## test: 运行所有测试
test:
	@echo "$(GREEN)Running tests...$(NC)"
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	@echo "$(GREEN)Tests complete!$(NC)"
	@$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report: coverage.html$(NC)"

## test-quick: 快速测试（不带竞态检测）
test-quick:
	@echo "$(GREEN)Running tests (quick)...$(NC)"
	$(GOTEST) ./...

## clean: 清理构建文件
clean:
	@echo "$(YELLOW)Cleaning...$(NC)"
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "$(GREEN)Clean complete!$(NC)"

## install: 安装 CLI 工具到 $GOPATH/bin 或 $HOME/go/bin
install: build
	@echo "$(GREEN)Installing $(NAME)...$(NC)"
	@install -d $$HOME/go/bin 2>/dev/null || install -d $(GOPATH)/bin
	@install bin/$(NAME) $$HOME/go/bin/$(NAME) 2>/dev/null || install bin/$(NAME) $(GOPATH)/bin/$(NAME)
	@echo "$(GREEN)Installed to: $$HOME/go/bin/$(NAME)$(NC)"

## uninstall: 卸载 CLI 工具
uninstall:
	@echo "$(YELLOW)Uninstalling $(NAME)...$(NC)"
	@rm -f $$HOME/go/bin/$(NAME) 2>/dev/null || rm -f $(GOPATH)/bin/$(NAME)
	@echo "$(GREEN)Uninstall complete!$(NC)"

## deps: 下载依赖
deps:
	@echo "$(GREEN)Downloading dependencies...$(NC)"
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "$(GREEN)Dependencies updated!$(NC)"

## fmt: 格式化代码
fmt:
	@echo "$(GREEN)Formatting code...$(NC)"
	$(GOFMT) -s -w .
	@echo "$(GREEN)Formatted!$(NC)"

## vet: 静态检查
vet:
	@echo "$(GREEN)Running go vet...$(NC)"
	$(GOVET) ./...
	@echo "$(GREEN)Vet complete!$(NC)"

## lint: 代码检查（需要 golangci-lint）
lint:
	@echo "$(GREEN)Running linter...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "$(YELLOW)golangci-lint not found. Install with:$(NC)"; \
		echo "$(YELLOW)curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$HOME/go/bin latest$(NC)"; \
	fi

## init-config: 初始化配置文件
init-config:
	@echo "$(GREEN)Initializing config...$(NC)"
	@mkdir -p $(CONFIG_DIR)
	@if [ ! -f $(CONFIG_DIR)/tencent.yaml ]; then \
		echo "secret_id: \"your_secret_id\"" > $(CONFIG_DIR)/tencent.yaml; \
		echo "secret_key: \"your_secret_key\"" >> $(CONFIG_DIR)/tencent.yaml; \
		echo "region: \"ap-guangzhou\"" >> $(CONFIG_DIR)/tencent.yaml; \
		echo "domain: \"example.com\"" >> $(CONFIG_DIR)/tencent.yaml; \
		echo "$(GREEN)Created $(CONFIG_DIR)/tencent.yaml$(NC)"; \
	else \
		echo "$(YELLOW)$(CONFIG_DIR)/tencent.yaml already exists$(NC)"; \
	fi

## run: 构建并运行 CLI（用于开发）
run: build
	@echo "$(GREEN)Running $(NAME)...$(NC)"
	./bin/$(NAME) $(ARGS)

## help: 显示帮助信息
help:
	@echo "$(GREEN)GADNS Makefile$(NC)"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  $(GREEN)all$(NC)           - 默认目标，构建项目"
	@echo "  $(GREEN)build$(NC)         - 构建 CLI 工具"
	@echo "  $(GREEN)build-dev$(NC)     - 开发模式构建"
	@echo "  $(GREEN)test$(NC)          - 运行所有测试"
	@echo "  $(GREEN)test-quick$(NC)    - 快速测试"
	@echo "  $(GREEN)clean$(NC)         - 清理构建文件"
	@echo "  $(GREEN)install$(NC)       - 安装到 $$HOME/go/bin"
	@echo "  $(GREEN)uninstall$(NC)     - 卸载 CLI 工具"
	@echo "  $(GREEN)deps$(NC)          - 下载依赖"
	@echo "  $(GREEN)fmt$(NC)           - 格式化代码"
	@echo "  $(GREEN)vet$(NC)           - 静态检查"
	@echo "  $(GREEN)lint$(NC)          - 代码检查"
	@echo "  $(GREEN)init-config$(NC)   - 初始化配置文件"
	@echo "  $(GREEN)run$(NC)           - 构建并运行（开发用）"
	@echo "  $(GREEN)help$(NC)          - 显示此帮助"
	@echo ""
	@echo "Examples:"
	@echo "  make                    # 构建"
	@echo "  make test               # 运行测试"
	@echo "  make install             # 安装到系统"
	@echo "  make run ARGS=\"help\"    # 运行命令"
	@echo ""
	@echo "Config: $(CONFIG_DIR)/tencent.yaml"
