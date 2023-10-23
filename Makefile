# Define your project name and binary name
PROJECT_NAME := mocktimism
BINARY_NAME := $(PROJECT_NAME)

# Set the project directory
PROJ_DIR := cmd

# Set the output directory
BIN_DIR := bin

# Specify the Go compiler
GO := go

# Define build flags
BUILD_FLAGS := -v

# Define targets
.PHONY: all build clean

all: build

build:
	@echo "Building $(BINARY_NAME)..."
	$(GO) build $(BUILD_FLAGS) -o ./$(PROJ_DIR)/$(BIN_DIR)/$(BINARY_NAME) ./$(PROJ_DIR)

clean:
	@echo "Cleaning up..."
	rm -rf ./$(PROJ_DIR)/$(BIN_DIR)