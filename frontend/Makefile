BUILD_DIR=../backend/dist
SVELTE_KIT_DIR=.svelte-kit
NODE_MODULES=node_modules

# Domyślny cel wykonuje się przy samym 'make'
.PHONY: all
all: install build

.PHONY: build
build:
	@echo "==> Building Svelte app..."
	npm run build
	@echo "==> Done! Built files are in $(BUILD_DIR)/"

.PHONY: install
install:
	@echo "==> Installing frontend dependencies..."
	npm install

.PHONY: dev
dev:
	@echo "==> Starting Svelte dev server..."
	npm run dev

.PHONY: test
test: build
	@echo "==> Running Svelte check (TypeScript & A11y)..."
	npm run check

.PHONY: clean
clean:
	@echo "==> Cleaning up build directories..."
	rm -rf $(BUILD_DIR)
	rm -rf $(SVELTE_KIT_DIR)

.PHONY: clean-all
clean-all: clean
	@echo "==> Removing node_modules..."
	rm -rf $(NODE_MODULES)
	rm -f package-lock.json