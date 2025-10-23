BINARY_NAME=goShop

# The directory to place the final binary into
BUILD_DIR=target

# The 'all' target is the default one executed when you just run 'make'
.PHONY: all
all: build

.PHONY: build
build:
	@echo "==> Building..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "==> Done! Binary is at $(BUILD_DIR)/$(BINARY_NAME)"

.PHONY: test
test:
	@echo "==> Testing..."
	go test ./... -v

.PHONY: clean
clean:
	@echo "==> Cleaning..."
	rm -rf $(BUILD_DIR)

.PHONY: install
install:
	@echo "==> Installing..."
	go mod tidy