.PHONY: all build clean install

all: build

build: build-frontend build-backend

clean: clean-frontend clean-backend

install: install-frontend install-backend

test: test-frontend test-backend

build-frontend:
	@echo "Building frontend..."
	@$(MAKE) -C frontend build

# Target to build the backend
build-backend:
	@echo "Building backend..."
	@$(MAKE) -C backend build

clean-frontend:
	@echo "Cleaning frontend..."
	@$(MAKE) -C frontend clean

clean-backend:
	@echo "Cleaning backend..."
	@$(MAKE) -C backend clean

install-frontend:
	@echo "Installing frontend dependencies..."
	@$(MAKE) -C frontend install

install-backend:
	@echo "Installing backend dependencies..."
	@$(MAKE) -C backend install

test-frontend:
	@echo "Testing frontend..."
	@$(MAKE) -C frontend test

test-backend:
	@echo "Testing backend..."
	@$(MAKE) -C backend test