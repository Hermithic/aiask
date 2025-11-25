# AIHelp CLI Makefile
# Cross-platform build script

BINARY_NAME=aiask
VERSION?=1.0.0
BUILD_DIR=build
DIST_DIR=dist

# Go build flags
LDFLAGS=-ldflags "-s -w -X github.com/Hermithic/aiask/internal/cli.Version=$(VERSION)"

# Default target
.PHONY: all
all: clean build

# Clean build directory
.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)
	mkdir -p $(BUILD_DIR)

# Build for current platform
.PHONY: build
build:
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/aiask

# Build for all platforms
.PHONY: build-all
build-all: clean build-windows build-linux build-darwin

# Build for Windows (amd64 and arm64)
.PHONY: build-windows
build-windows:
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/aiask
	GOOS=windows GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-arm64.exe ./cmd/aiask

# Build for Linux (amd64 and arm64)
.PHONY: build-linux
build-linux:
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/aiask
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 ./cmd/aiask

# Build for macOS (amd64 and arm64)
.PHONY: build-darwin
build-darwin:
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/aiask
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/aiask

# Install locally
.PHONY: install
install: build
	cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/$(BINARY_NAME)

# Run tests
.PHONY: test
test:
	go test -v ./...

# Run linter
.PHONY: lint
lint:
	golangci-lint run

# Download dependencies
.PHONY: deps
deps:
	go mod download
	go mod tidy

# Create release archives
.PHONY: release
release: build-all
	mkdir -p $(BUILD_DIR)/release
	# Windows
	cd $(BUILD_DIR) && zip release/$(BINARY_NAME)-$(VERSION)-windows-amd64.zip $(BINARY_NAME)-windows-amd64.exe
	cd $(BUILD_DIR) && zip release/$(BINARY_NAME)-$(VERSION)-windows-arm64.zip $(BINARY_NAME)-windows-arm64.exe
	# Linux
	cd $(BUILD_DIR) && tar -czvf release/$(BINARY_NAME)-$(VERSION)-linux-amd64.tar.gz $(BINARY_NAME)-linux-amd64
	cd $(BUILD_DIR) && tar -czvf release/$(BINARY_NAME)-$(VERSION)-linux-arm64.tar.gz $(BINARY_NAME)-linux-arm64
	# macOS
	cd $(BUILD_DIR) && tar -czvf release/$(BINARY_NAME)-$(VERSION)-darwin-amd64.tar.gz $(BINARY_NAME)-darwin-amd64
	cd $(BUILD_DIR) && tar -czvf release/$(BINARY_NAME)-$(VERSION)-darwin-arm64.tar.gz $(BINARY_NAME)-darwin-arm64

# Build .deb package for Linux
.PHONY: deb
deb: build-linux
	mkdir -p $(BUILD_DIR)/deb/$(BINARY_NAME)_$(VERSION)_amd64/DEBIAN
	mkdir -p $(BUILD_DIR)/deb/$(BINARY_NAME)_$(VERSION)_amd64/usr/local/bin
	cp $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(BUILD_DIR)/deb/$(BINARY_NAME)_$(VERSION)_amd64/usr/local/bin/$(BINARY_NAME)
	cp $(DIST_DIR)/deb/DEBIAN/control $(BUILD_DIR)/deb/$(BINARY_NAME)_$(VERSION)_amd64/DEBIAN/
	sed -i 's/VERSION/$(VERSION)/g' $(BUILD_DIR)/deb/$(BINARY_NAME)_$(VERSION)_amd64/DEBIAN/control
	dpkg-deb --build $(BUILD_DIR)/deb/$(BINARY_NAME)_$(VERSION)_amd64
	mv $(BUILD_DIR)/deb/$(BINARY_NAME)_$(VERSION)_amd64.deb $(BUILD_DIR)/release/

# Generate checksums
.PHONY: checksums
checksums:
	cd $(BUILD_DIR)/release && sha256sum * > checksums.txt

# Help
.PHONY: help
help:
	@echo "AIask CLI Build Targets:"
	@echo "  make build       - Build for current platform"
	@echo "  make build-all   - Build for all platforms"
	@echo "  make install     - Install to GOPATH/bin"
	@echo "  make test        - Run tests"
	@echo "  make deps        - Download dependencies"
	@echo "  make release     - Create release archives"
	@echo "  make deb         - Build .deb package"
	@echo "  make clean       - Clean build directory"
	@echo ""
	@echo "Variables:"
	@echo "  VERSION          - Set version (default: $(VERSION))"

