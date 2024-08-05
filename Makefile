# Define directories
BACKEND_DIR=backend
FRONTEND_DIR=frontend

# Backend settings
BACKEND_BINARY=backend

# Frontend settings
FRONTEND_BUILD_DIR=$(FRONTEND_DIR)/build

# Commands
.PHONY: all setup backend frontend clean run-backend run-frontend

all: setup build

setup: setup-backend setup-frontend

setup-backend:
	@echo "Setting up backend..."
	mkdir -p $(BACKEND_DIR)
	cd $(BACKEND_DIR) && \
	python3 -m venv venv \
	&& source venv/bin/activate \
	&& pip install django-cors-headers django djangorestframework pillow \
	&& django-admin startproject backend . \
	&& python manage.py migrate

setup-frontend:
	@echo "Setting up frontend..."
	mkdir -p $(FRONTEND_DIR)
	cd $(FRONTEND_DIR) && npm install

build: build-backend build-frontend

build-backend:
	@echo "Building backend..."
	cd $(BACKEND_DIR) && go build -o $(BACKEND_BINARY) $(BACKEND_MAIN)

build-frontend:
	@echo "Building frontend..."
	cd $(FRONTEND_DIR) && npm run build

run-backend:
	@echo "Running backend..."
	cd $(BACKEND_DIR) && go run $(BACKEND_MAIN)

run-frontend:
	@echo "Running frontend..."
	cd $(FRONTEND_DIR) && npm start

clean:
	@echo "Cleaning up..."
	rm -f $(BACKEND_DIR)/$(BACKEND_BINARY)
	rm -rf $(FRONTEND_BUILD_DIR)
