
# The 'all' target is the default one executed when you just run 'make'
.PHONY: all
all: build

.PHONY: build
build:
	@echo "==> Building..."
	 npm run build

.PHONY: test
test:
	@echo "==> Testing..."
	npm run test:once

.PHONY: clean
clean:
	@echo "==> Cleaning..."
	rm -rf dist node_modules

.PHONY: install
install:
	@echo "==> Installing..."
	npm install